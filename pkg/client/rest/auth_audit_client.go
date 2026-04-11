package rest

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/client/audit"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

type AuthAuditClient struct {
	*client.BaseAuditClient[*dto.AuthAuditDTO]
	client audit.ClientService
	conf   *AuditClientConfig
	log    logger.Logger
}

var _ client.AuditClient[*dto.AuthAuditDTO] = (*AuthAuditClient)(nil)
var _ client.AuthAuditClient = (*AuthAuditClient)(nil)

func NewAuthAuditClient(
	conf *AuditClientConfig,
	tokenProvider auth.TokenProvider,
	log logger.Logger,
) *AuthAuditClient {
	res := &AuthAuditClient{
		log:    log.GetLogger("AuthAuditClient"),
		conf:   conf,
		client: audit.NewClientWithBearerToken(conf.Host, conf.BasePath, conf.Scheme, ""),
	}
	res.BaseAuditClient = client.NewBaseAuditClient[*dto.AuthAuditDTO](
		"auth",
		conf.poolConf,
		res.auditAction,
		tokenProvider,
		log,
	)

	return res
}

func (aac *AuthAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *dto.AuthAuditDTO,
	token string,
) error {
	aac.GetLogger().Debugf("auth audit action worker %d start", workerIndex)
	defer aac.GetLogger().Debugf("auth audit action worker %d finish", workerIndex)

	clientDTO := MapAuthDtoSDKToRest(data)
	// request
	_, err := aac.client.PostAPIV1AuditAuth(
		audit.NewPostAPIV1AuditAuthParams().
			WithContext(ctx).
			WithTimeout(aac.conf.ReadTimeout).
			WithInput(clientDTO),
		func(op *runtime.ClientOperation) {
			op.AuthInfo = httptransport.BearerToken(token)
		},
	)
	if err != nil {
		return errs.NewCommonError("auth audit action failed", err)
	}

	return nil
}

func (aac *AuthAuditClient) GetConfig() *AuditClientConfig {
	return aac.conf
}

func (aac *AuthAuditClient) GetLogger() logger.Logger {
	return aac.log
}
