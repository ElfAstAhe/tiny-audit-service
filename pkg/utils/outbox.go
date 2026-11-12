package utils

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/http"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
)

// DefaultBuildDataAuditOutbox формирует сущность транзакционного аудит-лога (domain.DataAuditOutbox)
// со статусом ожидания отправки (domain.OutboxStatusPending) и JSON-полезной нагрузкой из переданного билдера.
func DefaultBuildDataAuditOutbox(_ context.Context, builder *DataAuditBuilder) (*domain.DataAuditOutbox, error) {
	if utils.IsNil(builder) {
		return nil, errs.NewInvalidArgumentError("builder", "is required")
	}
	audit := builder.Build()
	payload, err := BuildJSONPayload(audit)
	if err != nil {
		return nil, errs.NewCommonError("unable to build payload", err)
	}

	res := domain.NewEmptyDataAuditOutbox()
	res.InstanceID = audit.InstanceID
	res.ContentType = http.MediaTypeApplicationJSON
	res.Status = domain.OutboxStatusPending
	res.Payload = payload
	res.CreatedAt = audit.EventDate

	return res, nil
}
