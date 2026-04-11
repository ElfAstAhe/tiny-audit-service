package grpc

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	transportauth "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

type AuthAuditClient struct {
	*BaseAuditClient[*dto.AuthAuditDTO]
	client pb.AuthAuditServiceClient
	log    logger.Logger
}

var _ client.AuditClient[*dto.AuthAuditDTO] = (*AuthAuditClient)(nil)
var _ client.AuthAuditClient = (*AuthAuditClient)(nil)

func NewAuthAuditClient(
	conf *AuditClientConfig,
	tokenProvider transportauth.TokenProvider,
	log logger.Logger,
) *AuthAuditClient {
	res := &AuthAuditClient{
		log: log.GetLogger("AuthAuditClient"),
	}

	res.BaseAuditClient = NewBaseAuditClient[*dto.AuthAuditDTO](
		"auth",
		conf,
		tokenProvider,
		res.auditAction,
		log,
	)

	return res
}

func (acc *AuthAuditClient) Start(ctx context.Context) error {
	if err := acc.BaseAuditClient.Start(ctx); err != nil {
		return err
	}

	acc.client = pb.NewAuthAuditServiceClient(acc.conn)

	return nil
}

func (acc *AuthAuditClient) Stop(stopTimeout time.Duration) error {
	if err := acc.BaseAuditClient.Stop(stopTimeout); err != nil {
		return err
	}

	return nil
}

func (acc *AuthAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *dto.AuthAuditDTO,
	token string,
) error {
	// transform
	clientDTO := MapAuthDtoSDKToGRPC(data)

	// request
	_, err := acc.client.Audit(
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
