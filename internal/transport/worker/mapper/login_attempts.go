package mapper

import (
	"encoding/json"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
)

func MapMessageToLoginAttemptWorkerJob(data amqp.Message) (*dto.LoginAttemptWorkerJob, error) {
	if utils.IsNil(data) {
		return nil, nil
	}

	localDto, err := MapMessageBinDataToLoginAttemptDTO(data.GetPayload())
	if err != nil {
		return nil, err
	}

	return &dto.LoginAttemptWorkerJob{
		Data:    localDto,
		Message: data,
	}, nil
}

func MapMessageBinDataToLoginAttemptDTO(data []byte) (*dto.LoginAttemptEventDTO, error) {
	if len(data) == 0 {
		return nil, errs.NewCommonError("empty data", nil)
	}
	res := &dto.LoginAttemptEventDTO{}

	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, errs.NewCommonError("unmarshal login attempt dto", err)
	}

	return res, nil
}

func MapLoginAttemptEventDTOToAuthAuditModel(data *dto.LoginAttemptEventDTO) *domain.AuthAudit {
	if utils.IsNil(data) {
		return nil
	}

	res := &domain.AuthAudit{
		Source:    data.NodeName,
		EventDate: data.EventDate,
		Event:     domain.AuthEventLogin,
		RequestID: data.RequestID,
		TraceID:   data.TraceID,
		Username:  data.Username,
	}
	res.Status = domain.AuditStatusFail
	if data.Success {
		res.Status = domain.AuditStatusSuccess
	}

	return res
}
