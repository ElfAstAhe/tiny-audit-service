package rest

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/client/audit"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/models"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
	"github.com/go-openapi/runtime"
	oapirtcli "github.com/go-openapi/runtime/client"
)

type AuthAuditClient struct {
	*BaseAuditClient[*models.AuthAuditDTO]
}

var _ AuditClient[*models.AuthAuditDTO] = (*AuthAuditClient)(nil)

func NewAuthAuditClient(
	conf *AuditClientConfig,
	tokenProvider auth.TokenProvider,
	log logger.Logger,
) *AuthAuditClient {
	res := &AuthAuditClient{}
	res.BaseAuditClient = NewBaseAuditClient[*models.AuthAuditDTO](
		"auth",
		conf,
		res.auditAction,
		tokenProvider,
		log,
	)

	return res
}

func (aac *AuthAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *models.AuthAuditDTO,
	token string,
) error {
	// request
	_, err := aac.client.PostAPIV1AuditAuth(
		audit.NewPostAPIV1AuditAuthParams().
			WithContext(ctx).
			WithTimeout(aac.conf.timeout).
			WithInput(data),
		func(op *runtime.ClientOperation) {
			op.AuthInfo = oapirtcli.BearerToken(token)
		},
	)
	if err != nil {
		aac.totalLost.Add(1)
		return errs.NewCommonError("auth audit action failed", err)
	}

	return nil
}
