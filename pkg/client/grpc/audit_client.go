package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	libauth "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	auditClientNameTemplate string = "gRPC-audit-client-%s"
)

// AuditClient реализует интерфейс client.AuditClient, предоставляя высокопроизводительный,
// потокобезопасный дженерик-клиент для асинхронной отправки логов аудита по протоколу gRPC.
//
// Инкапсулирует пул воркеров BaseAuditClient, управляет низкоуровневым жизненным циклом сетевого
// соединения grpc.ClientConn, механизмами KeepAlive и реализует интерфейс credentials.PerRPCCredentials
// для автоматического внедрения авторизационных токенов в метаданные каждого RPC-вызова.
type AuditClient[D any] struct {
	*client.BaseAuditClient[D]
	name          string
	conn          *grpc.ClientConn
	opts          *Options[D]
	tokenProvider libauth.TokenProvider
	log           logger.Logger
}

// Проверяем соответствие типа строгим доменным контрактам на этапе компиляции для DTO авторизации и данных
var _ client.AuditClient[*dto.AuthAuditDTO] = (*AuditClient[*dto.AuthAuditDTO])(nil)
var _ client.AuditClient[*dto.DataAuditDTO] = (*AuditClient[*dto.DataAuditDTO])(nil)

// NewAuditClient — фабричный конструктор gRPC-клиента аудита.
// Выполняет сквозную валидацию опций, инициализирует вложенный пул воркеров и настраивает изолированное логирование.
func NewAuditClient[D any](options ...Option[D]) (*AuditClient[D], error) {
	opts := NewOptions[D]()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, errs.NewTlCommonError("NewAuditClient", "client options validation failed", err)
	}
	// instance
	res := &AuditClient[D]{
		name:          fmt.Sprintf(auditClientNameTemplate, opts.Pool.Name),
		opts:          opts,
		tokenProvider: opts.TokenProvider,
		log:           opts.Pool.Logger.GetLogger(fmt.Sprintf(auditClientNameTemplate, opts.Pool.Name)),
	}
	// base
	base, err := client.NewBaseAuditClient[D](
		client.WithName[D](res.name),
		client.WithAuditAction[D](opts.AuditAction),
		client.WithTokenProvider[D](opts.TokenProvider),
		client.WithLogger[D](opts.Pool.Logger),
		client.WithPoolWorkerCount[D](opts.Pool.WorkerCount),
		client.WithPoolDataCapacity[D](opts.Pool.DataCapacity),
		client.WithPoolCompleteProcess[D](opts.Pool.CompleteProcess),
		client.WithPoolStopTimeout[D](opts.Pool.StopTimeout),
	)
	if err != nil {
		return nil, errs.NewTlCommonError("NewAuditClient", "base client create failed", err)
	}
	// setup
	res.BaseAuditClient = base

	return res, nil
}

// Start осуществляет инициализацию сетевого транспортного соединения gRPC и запускает внутренний пул воркеров.
func (ac *AuditClient[D]) Start(ctx context.Context) error {
	ac.GetLogger().Debugf("Start start")
	defer ac.GetLogger().Debugf("Start finish")
	// connection
	conn, err := ac.createConnection(ctx)
	if err != nil {
		return errs.NewCommonError("gRPC base audit client create connection failed", err)
	}
	ac.conn = conn
	// pool
	err = ac.BaseAuditClient.Start(ctx)
	if err != nil {
		_ = conn.Close()

		return errs.NewCommonError("gRPC base audit client start pool failed", err)
	}

	return nil
}

// Stop выполняет координируемую остановку пула воркеров (досылку буфера) и закрывает gRPC-соединение ClientConn.
func (ac *AuditClient[D]) Stop(ctx context.Context) error {
	ac.GetLogger().Debugf("Stop start")
	defer ac.GetLogger().Debugf("Stop finish")

	var stopErrs []error
	stopErrs = append(stopErrs, ac.BaseAuditClient.Stop(ctx))
	if !utils.IsNil(ac.conn) {
		stopErrs = append(stopErrs, ac.conn.Close())
	}
	err := errors.Join(stopErrs...)
	if err != nil {
		return errs.NewCommonError("gRPC base audit client stop failed", err)
	}

	return nil
}

// buildTransportCredentials конфигурирует параметры безопасности TLS или возвращает небезопасные учетные данные.
func (ac *AuditClient[D]) buildTransportCredentials(opts *Options[D]) credentials.TransportCredentials {
	if utils.IsNil(opts) {
		return nil
	}
	if utils.IsNil(opts.CertPool) {
		return insecure.NewCredentials()
	}

	return credentials.NewClientTLSFromCert(opts.CertPool, opts.CertServerNameOverride)
}

// buildKeepAliveParams формирует низкоуровневую структуру параметров KeepAlive для контроля живости gRPC-канала.
func (ac *AuditClient[D]) buildKeepAliveParams(opts *Options[D]) keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:                opts.KeepAliveTime,
		Timeout:             opts.KeepAliveTimeout,
		PermitWithoutStream: opts.KeepAlivePermitWithoutStream,
	}
}

// buildDialOptions агрегирует все конфигурационные параметры, параметры безопасности, KeepAlive и интерцепторы в срез опций DialOption.
func (ac *AuditClient[D]) buildDialOptions(opts *Options[D]) []grpc.DialOption {
	if utils.IsNil(opts) {
		return []grpc.DialOption{}
	}

	return []grpc.DialOption{
		grpc.WithTransportCredentials(ac.buildTransportCredentials(opts)),
		grpc.WithKeepaliveParams(ac.buildKeepAliveParams(opts)),
		grpc.WithPerRPCCredentials(ac),
		grpc.WithIdleTimeout(opts.IdleTimeout),
	}
}

// createConnection инициирует создание gRPC-клиента и блокирует поток до перехода канала в состояние Ready или по таймауту.
func (ac *AuditClient[D]) createConnection(ctx context.Context) (*grpc.ClientConn, error) {
	ac.GetLogger().Debugf("createConnection start")
	defer ac.GetLogger().Debugf("createConnection finish")

	// dial options
	dialOptions := ac.buildDialOptions(ac.opts)

	// connection
	conn, err := grpc.NewClient(ac.opts.Target, dialOptions...)
	if err != nil {
		return nil, errs.NewCommonError("gRPC base audit client create connection failed", err)
	}

	// conn timeout
	if ac.opts.ConnTimeout > 0 {
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, ac.opts.ConnTimeout)
		defer timeoutCancel()
		for {
			state := conn.GetState()
			if state == connectivity.Ready {
				break
			}
			if !conn.WaitForStateChange(timeoutCtx, state) {
				_ = conn.Close()
				return nil, errs.NewCommonError("gRPC base audit client connection timeout", nil)
			}
		}
	}

	return conn, nil
}

// GetOpts возвращает ссылку на конфигурационный паспорт Options текущего gRPC-клиента.
func (ac *AuditClient[D]) GetOpts() *Options[D] {
	return ac.opts
}

// GetLogger возвращает изолированный экземпляр структурированного логгера компонента.
func (ac *AuditClient[D]) GetLogger() logger.Logger {
	return ac.log
}

// GetRequestMetadata извлекает актуальный токен авторизации и внедряет его в заголовки метаданных каждого RPC-вызова.
func (ac *AuditClient[D]) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	token, err := ac.tokenProvider.GetAccessToken()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		auth.DefaultMetadataName: token,
	}, nil
}

// RequireTransportSecurity сообщает gRPC-рантайму, требуется ли обязательное TLS-шифрование для передачи токенов.
func (ac *AuditClient[D]) RequireTransportSecurity() bool {
	return !utils.IsNil(ac.opts.CertPool)
}
