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

type DataAuditClient struct {
	*client.BaseAuditClient[*dto.DataAuditDTO]
	client audit.ClientService
	conf   *AuditClientConfig
	log    logger.Logger
}

var _ client.AuditClient[*dto.DataAuditDTO] = (*DataAuditClient)(nil)
var _ client.DataAuditClient = (*DataAuditClient)(nil)

func NewDataAuditClient(
	conf *AuditClientConfig,
	tokenProvider auth.TokenProvider,
	log logger.Logger,
) *DataAuditClient {
	res := &DataAuditClient{
		log:    log.GetLogger("DataAuditClient"),
		conf:   conf,
		client: audit.NewClientWithBearerToken(conf.Host, conf.BasePath, conf.Scheme, ""),
	}
	res.BaseAuditClient = client.NewBaseAuditClient[*dto.DataAuditDTO](
		"data",
		conf.poolConf,
		res.auditAction,
		tokenProvider,
		log,
	)

	return res
}

func (dac *DataAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *dto.DataAuditDTO,
	token string,
) error {
	dac.GetLogger().Debugf("data audit action worker %d start", workerIndex)
	defer dac.GetLogger().Debugf("data audit action worker %d finish", workerIndex)

	clientDTO := MapDataDtoSDKToRest(data)
	// request
	_, err := dac.client.PostAPIV1AuditData(
		audit.NewPostAPIV1AuditDataParams().
			WithContext(ctx).
			WithTimeout(dac.conf.ReadTimeout).
			WithInput(clientDTO),
		func(op *runtime.ClientOperation) {
			op.AuthInfo = httptransport.BearerToken(token)
		},
	)
	if err != nil {
		dac.IncLostCounter()
		return errs.NewCommonError("data audit action failed", err)
	}

	return nil
}

func (dac *DataAuditClient) GetConfig() *AuditClientConfig {
	return dac.conf
}

func (dac *DataAuditClient) GetLogger() logger.Logger {
	return dac.log
}
