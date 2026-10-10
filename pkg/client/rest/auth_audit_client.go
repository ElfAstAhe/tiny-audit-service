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
	authClientNameTemplate string = "rest-auth-audit-client-%s"
)

type AuthAuditClient struct {
	*client.BaseAuditClient[*dto.AuthAuditDTO]
	name   string
	client audit.ClientService
	opts   *AuditClientOptions[*dto.AuthAuditDTO]
	log    logger.Logger
}

var _ client.AuditClient[*dto.AuthAuditDTO] = (*AuthAuditClient)(nil)
var _ client.AuthAuditClient = (*AuthAuditClient)(nil)
var _ container.Runner = (*AuthAuditClient)(nil)

func NewAuthAuditClient(options ...AuditClientOption[*dto.AuthAuditDTO]) (*AuthAuditClient, error) {
	opts := NewAuditClientOptions[*dto.AuthAuditDTO]()
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
	res := &AuthAuditClient{
		name:   fmt.Sprintf(authClientNameTemplate, opts.Pool.Name),
		opts:   opts,
		log:    opts.Pool.Logger.GetLogger(fmt.Sprintf(authClientNameTemplate, opts.Pool.Name)),
		client: audit.NewClientWithBearerToken(u.Host, u.Path, u.Scheme, ""),
	}
	// base
	base, err := client.NewBaseAuditClient[*dto.AuthAuditDTO](
		client.WithName[*dto.AuthAuditDTO](opts.Pool.Name),
		client.WithAuditAction[*dto.AuthAuditDTO](res.auditAction),
		client.WithLogger[*dto.AuthAuditDTO](opts.Pool.Logger),
		client.WithPoolWorkerCount[*dto.AuthAuditDTO](opts.Pool.WorkerCount),
		client.WithPoolDataCapacity[*dto.AuthAuditDTO](opts.Pool.DataCapacity),
		client.WithPoolCompleteProcess[*dto.AuthAuditDTO](opts.Pool.CompleteProcess),
		client.WithPoolStopTimeout[*dto.AuthAuditDTO](opts.Pool.StopTimeout),
		client.WithTokenProvider[*dto.AuthAuditDTO](opts.TokenProvider),
	)
	if err != nil {
		return nil, errs.NewTlCommonError("NewAuthAuditClient", "failed to create base audit client", err)
	}
	// setup
	res.BaseAuditClient = base

	return res, nil
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
	_, err := aac.client.PostAPIV1AuditAuthContext(
		ctx,
		audit.NewPostAPIV1AuditAuthParams().
			WithTimeout(aac.opts.ReadTimeout).
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

func (aac *AuthAuditClient) GetName() string {
	return aac.name
}

func (aac *AuthAuditClient) GetOpts() *AuditClientOptions[*dto.AuthAuditDTO] {
	return aac.opts
}

func (aac *AuthAuditClient) GetLogger() logger.Logger {
	return aac.log
}
