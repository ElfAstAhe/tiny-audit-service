package rest

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/client/audit"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

const (
	dataClientNameTemplate string = "rest-data-audit-client-%s"
)

type DataAuditClient struct {
	*client.BaseAuditClient[*dto.DataAuditDTO]
	name   string
	client audit.ClientService
	opts   *AuditClientOptions[*dto.DataAuditDTO]
	log    logger.Logger
}

var _ client.AuditClient[*dto.DataAuditDTO] = (*DataAuditClient)(nil)
var _ client.DataAuditClient = (*DataAuditClient)(nil)
var _ container.Runner = (*DataAuditClient)(nil)

func NewDataAuditClient(options ...AuditClientOption[*dto.DataAuditDTO]) (*DataAuditClient, error) {
	opts := NewAuditClientOptions[*dto.DataAuditDTO]()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, errs.NewTlCommonError("NewAuthAuditClient", "client options validation failed", err)
	}
	u, err := url.Parse(opts.BaseURL)
	if err != nil {
		return nil, errs.NewTlCommonError("NewAuthAuditClient", "failed to parse base url", err)
	}
	// instance
	res := &DataAuditClient{
		name:   fmt.Sprintf(authClientNameTemplate, opts.Pool.Name),
		opts:   opts,
		log:    opts.Pool.Logger.GetLogger(fmt.Sprintf(authClientNameTemplate, opts.Pool.Name)),
		client: audit.NewClientWithBearerToken(u.Host, u.Path, u.Scheme, ""),
	}
	// base
	base, err := client.NewBaseAuditClient[*dto.DataAuditDTO](
		client.WithName[*dto.DataAuditDTO](opts.Pool.Name),
		client.WithAuditAction[*dto.DataAuditDTO](res.auditAction),
		client.WithLogger[*dto.DataAuditDTO](opts.Pool.Logger),
		client.WithPoolWorkerCount[*dto.DataAuditDTO](opts.Pool.WorkerCount),
		client.WithPoolDataCapacity[*dto.DataAuditDTO](opts.Pool.DataCapacity),
		client.WithPoolCompleteProcess[*dto.DataAuditDTO](opts.Pool.CompleteProcess),
		client.WithPoolStopTimeout[*dto.DataAuditDTO](opts.Pool.StopTimeout),
	)
	if err != nil {
		return nil, errs.NewTlCommonError("NewAuthAuditClient", "failed to create base audit client", err)
	}
	// setup
	res.BaseAuditClient = base

	return res, nil
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
	_, err := dac.client.PostAPIV1AuditDataContext(
		ctx,
		audit.NewPostAPIV1AuditDataParams().
			WithTimeout(dac.opts.ReadTimeout).
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

func (dac *DataAuditClient) GetName() string {
	return dac.name
}

func (dac *DataAuditClient) GetConfig() *AuditClientOptions[*dto.DataAuditDTO] {
	return dac.opts
}

func (dac *DataAuditClient) GetLogger() logger.Logger {
	return dac.log
}
