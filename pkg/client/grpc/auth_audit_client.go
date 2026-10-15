package grpc

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

const (
	authAuditClientNameTemplate string = "gRPC-auth-audit-client-%s"
)

// AuthAuditClient реализует интерфейсы client.AuditClient и client.AuthAuditClient,
// предоставляя специализированную обертку над базовым gRPC-клиентом для фиксации событий аутентификации.
//
// Инкапсулирует типизированный AuditClient, генерирует строго типизированный gRPC-сервис pb.AuthAuditServiceClient,
// управляет трассировкой жизненного цикла через логгер и связывает воркеры с удаленным вызовом процедур (RPC).
type AuthAuditClient struct {
	*AuditClient[*dto.AuthAuditDTO]
	name string
	svc  AuditorService[*pb.AuthAuditRequest]
	log  logger.Logger
}

// Контролируем реализацию контрактов SDK аудит-клиента на этапе компиляции
var _ client.AuditClient[*dto.AuthAuditDTO] = (*AuthAuditClient)(nil)
var _ client.AuthAuditClient = (*AuthAuditClient)(nil)

// NewAuthAuditClient — фабричный конструктор gRPC-клиента аудита авторизационных событий.
// Настраивает опции, проверяет инварианты и инжектирует кастомный метод auditAction в базовый дженерик-клиент.
func NewAuthAuditClient(options ...Option[*dto.AuthAuditDTO]) (*AuthAuditClient, error) {
	// options
	opts := NewOptions[*dto.AuthAuditDTO]()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, errs.NewTlCommonError("NewAuthAuditClient", "gRPC audit client options validation failed", err)
	}
	// instance
	res := &AuthAuditClient{
		name: fmt.Sprintf(authAuditClientNameTemplate, opts.Pool.Name),
		log:  opts.Pool.Logger.GetLogger(fmt.Sprintf(authAuditClientNameTemplate, opts.Pool.Name)),
	}
	// embed
	embed, err := NewAuditClient[*dto.AuthAuditDTO](append(options, withAuditAction[*dto.AuthAuditDTO](res.auditAction))...)
	if err != nil {
		return nil, errs.NewTlCommonError("NewAuthAuditClient", "audit client creation failed", err)
	}
	// setup
	res.AuditClient = embed

	return res, nil
}

// Start осуществляет запуск сетевых ресурсов базового клиента и инициализирует строго типизированный gRPC-сервис.
func (aac *AuthAuditClient) Start(ctx context.Context) error {
	aac.GetLogger().Debugf("Start start")
	defer aac.GetLogger().Debugf("Start finish")

	if err := aac.AuditClient.Start(ctx); err != nil {
		return err
	}

	aac.svc = pb.NewAuthAuditServiceClient(aac.conn)

	return nil
}

// Stop осуществляет фазу контролируемой мягкой остановки клиента, зануляет ссылку на gRPC-сервис и деаллоцирует ресурсы соединения.
func (aac *AuthAuditClient) Stop(ctx context.Context) error {
	aac.GetLogger().Debugf("Stop start")
	defer aac.GetLogger().Debugf("Stop finish")

	aac.svc = nil

	return aac.AuditClient.Stop(ctx)
}

// auditAction — это внутренняя прикладная функция (экшен), обрабатывающая извлеченный из очереди лог в рамках конкретного воркера.
// Выполняет маппинг в Protobuf-структуру и осуществляет блокирующий сетевой вызов к удаленному gRPC-серверу аудита.
func (aac *AuthAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *dto.AuthAuditDTO,
	token string,
) error {
	aac.GetLogger().Debugf("auditAction start")
	defer aac.GetLogger().Debugf("auditAction finish")

	// transform
	clientDTO := MapAuthDtoSDKToGRPC(data)

	// request
	_, err := aac.svc.Audit(
		ctx,
		pb.AuthAuditRequest_builder{
			Data: clientDTO,
		}.Build(),
	)
	if err != nil {
		return errs.NewCommonError("auth audit action failed", err)
	}

	return nil
}
