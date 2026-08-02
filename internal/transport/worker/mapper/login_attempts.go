package mapper

import (
	"encoding/json"

	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	appdto "github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
)

func ToAuthAuditDTO(data amqp.Message) (*dto.AuthAuditDTO, error) {
	if utils.IsNil(data) {
		return nil, nil
	}
	if len(data.GetPayload()) == 0 {
		return nil, nil
	}
	res := &dto.AuthAuditDTO{
		AuthAuditDTO: &appdto.AuthAuditDTO{},
		Message:      data,
	}
	if err := json.Unmarshal(data.GetPayload(), res); err != nil {
		return nil, err
	}

	return res, nil
}
