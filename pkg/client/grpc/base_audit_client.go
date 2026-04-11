package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	transportauth "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type BaseAuditClient[D any] struct {
	*client.BaseAuditClient[D]
	conn          *grpc.ClientConn
	conf          *AuditClientConfig
	tokenProvider transportauth.TokenProvider
	log           logger.Logger
}

var _ client.AuditClient[*dto.AuthAuditDTO] = (*BaseAuditClient[*dto.AuthAuditDTO])(nil)
var _ client.AuditClient[*dto.DataAuditDTO] = (*BaseAuditClient[*dto.DataAuditDTO])(nil)

func NewBaseAuditClient[D any](
	name string,
	conf *AuditClientConfig,
	tokenProvider transportauth.TokenProvider,
	auditAction client.AuditAction[D],
	log logger.Logger,
) *BaseAuditClient[D] {
	res := &BaseAuditClient[D]{
		conf:          conf,
		tokenProvider: tokenProvider,
		log:           log.GetLogger("BaseAuditClient"),
	}

	res.BaseAuditClient = client.NewBaseAuditClient[D](
		name,
		conf.poolConf,
		auditAction,
		tokenProvider,
		log,
	)

	return res
}

func (bac *BaseAuditClient[D]) Start(ctx context.Context) error {
	// connection
	conn, err := bac.createGRPCConnection(ctx)
	if err != nil {
		return errs.NewCommonError("gRPC base audit client create connection failed", err)
	}
	// pool
	err = bac.BaseAuditClient.Start(ctx)
	if err != nil {
		_ = conn.Close()

		return errs.NewCommonError("gRPC base audit client start pool failed", err)
	}

	bac.conn = conn

	return nil
}

func (bac *BaseAuditClient[D]) Stop(stopTimeout time.Duration) error {
	var stopErrs []error
	stopErrs = append(stopErrs, bac.BaseAuditClient.Stop(stopTimeout))
	stopErrs = append(stopErrs, bac.conn.Close())
	err := errors.Join(stopErrs...)
	if err != nil {
		return errs.NewCommonError("gRPC base audit client stop failed", err)
	}

	return nil
}

func (bac *BaseAuditClient[D]) createGRPCConnection(ctx context.Context) (*grpc.ClientConn, error) {
	// transport credentials
	var transportCreds credentials.TransportCredentials
	if bac.conf.GRPCConf.Secure {
		transportCreds = credentials.NewClientTLSFromCert(nil, "")
	} else {
		transportCreds = insecure.NewCredentials()
	}
	// kacp params
	kacp := keepalive.ClientParameters{
		Time:                bac.conf.GRPCConf.KATime,
		Timeout:             bac.conf.GRPCConf.KATimeOut,
		PermitWithoutStream: bac.conf.GRPCConf.KAPermitWOStream,
	}

	// dial options
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithPerRPCCredentials(bac),
		grpc.WithKeepaliveParams(kacp),
	}

	// connection
	conn, err := grpc.NewClient(bac.conf.Target, dialOptions...)
	if err != nil {
		return nil, errs.NewCommonError("gRPC base audit client create connection failed", err)
	}

	// conn timeout
	if bac.conf.GRPCConf.ConnTimeout > 0 {
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, bac.conf.GRPCConf.ConnTimeout)
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

func (bac *BaseAuditClient[D]) GetConfig() *AuditClientConfig {
	return bac.conf
}

func (bac *BaseAuditClient[D]) GetLogger() logger.Logger {
	return bac.log
}

func (bac *BaseAuditClient[D]) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := bac.tokenProvider.GetAccessToken()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		auth.DefaultMetadataName: token,
	}, nil
}

func (bac *BaseAuditClient[D]) RequireTransportSecurity() bool {
	return bac.conf.GRPCConf.Secure
}
