package mapper

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

func MapAuthAuditModelToDTO(model *domain.AuthAudit) *dto.AuthAuditDTO {
	if model == nil {
		return nil
	}

	return &dto.AuthAuditDTO{
		ID:           model.ID,
		Source:       model.Source,
		EventDate:    model.EventDate,
		Event:        model.Event,
		Status:       model.Status,
		RequestID:    model.RequestID,
		Username:     model.Username,
		AccessToken:  model.AccessToken,
		RefreshToken: model.RefreshToken,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func MapAuthAuditModelsToDTOs(models []*domain.AuthAudit) []*dto.AuthAuditDTO {
	res := make([]*dto.AuthAuditDTO, 0, len(models))

	for _, model := range models {
		data := MapAuthAuditModelToDTO(model)
		if data != nil {
			res = append(res, data)
		}
	}

	return res
}

func MapAuthAuditDTOToModel(data *dto.AuthAuditDTO) *domain.AuthAudit {
	if data == nil {
		return nil
	}

	return &domain.AuthAudit{
		ID:           data.ID,
		Source:       data.Source,
		EventDate:    data.EventDate,
		Event:        data.Event,
		Status:       data.Status,
		RequestID:    data.RequestID,
		Username:     data.Username,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
}
