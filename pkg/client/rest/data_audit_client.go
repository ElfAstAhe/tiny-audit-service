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

type DataAuditClient struct {
	*BaseAuditClient[*models.DataAuditDTO]
}

var _ AuditClient[*models.DataAuditDTO] = (*DataAuditClient)(nil)

func NewDataAuditClient(
	conf *AuditClientConfig,
	tokenProvider auth.TokenProvider,
	log logger.Logger,
) *DataAuditClient {
	res := &DataAuditClient{}
	res.BaseAuditClient = NewBaseAuditClient[*models.DataAuditDTO](
		"data",
		conf,
		res.auditAction,
		tokenProvider,
		log,
	)

	return res
}

func (dac *DataAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *models.DataAuditDTO,
	token string,
) error {
	// request
	_, err := dac.client.PostAPIV1AuditData(
		audit.NewPostAPIV1AuditDataParams().
			WithContext(ctx).
			WithTimeout(dac.conf.timeout).
			WithInput(data),
		func(op *runtime.ClientOperation) {
			op.AuthInfo = oapirtcli.BearerToken(token)
		},
	)
	if err != nil {
		dac.totalLost.Add(1)
		return errs.NewCommonError("data audit action failed", err)
	}

	return nil
}
