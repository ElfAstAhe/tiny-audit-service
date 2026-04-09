package mapper

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

func MapDataAuditValueModelToDTO(model *domain.DataAuditValue) *dto.DataAuditValueDTO {
	if model == nil {
		return nil
	}

	return &dto.DataAuditValueDTO{
		Name:        model.Name,
		Description: model.Description,
		Before:      model.Before,
		After:       model.After,
	}
}

func MapDataAuditValueModelsToDTOs(models []*domain.DataAuditValue) []*dto.DataAuditValueDTO {
	res := make([]*dto.DataAuditValueDTO, 0, len(models))

	for _, model := range models {
		data := MapDataAuditValueModelToDTO(model)
		if data != nil {
			res = append(res, data)
		}
	}

	return res
}

func MapDataAuditValueDTOToModel(data *dto.DataAuditValueDTO) *domain.DataAuditValue {
	if data == nil {
		return nil
	}

	return &domain.DataAuditValue{
		Name:        data.Name,
		Description: data.Description,
		Before:      data.Before,
		After:       data.After,
	}
}

func MapDataAuditValueDTOsToModels(data []*dto.DataAuditValueDTO) []*domain.DataAuditValue {
	res := make([]*domain.DataAuditValue, 0, len(data))

	for _, item := range data {
		model := MapDataAuditValueDTOToModel(item)
		if model != nil {
			res = append(res, model)
		}
	}

	return res
}

func MapDataAuditModelToDTO(model *domain.DataAudit) *dto.DataAuditDTO {
	if model == nil {
		return nil
	}

	return &dto.DataAuditDTO{
		ID:               model.ID,
		Source:           model.Source,
		EventDate:        model.EventDate,
		Event:            model.Event,
		Status:           model.Status,
		RequestID:        model.RequestID,
		TraceID:          model.TraceID,
		Username:         model.Username,
		InstanceID:       model.InstanceID,
		InstanceName:     model.InstanceName,
		InternalTypeName: model.InternalTypeName,
		TypeName:         model.TypeName,
		TypeDescription:  model.TypeDescription,
		Values:           MapDataAuditValueModelsToDTOs(model.Values),
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}

func MapDataAuditDTOToModel(data *dto.DataAuditDTO) *domain.DataAudit {
	if data == nil {
		return nil
	}

	return &domain.DataAudit{
		ID:               data.ID,
		Source:           data.Source,
		EventDate:        data.EventDate,
		Event:            data.Event,
		Status:           data.Status,
		RequestID:        data.RequestID,
		TraceID:          data.TraceID,
		Username:         data.Username,
		InstanceID:       data.InstanceID,
		InstanceName:     data.InstanceName,
		InternalTypeName: data.InternalTypeName,
		TypeName:         data.TypeName,
		TypeDescription:  data.TypeDescription,
		Values:           MapDataAuditValueDTOsToModels(data.Values),
		CreatedAt:        data.CreatedAt,
		UpdatedAt:        data.UpdatedAt,
	}
}

func MapDataAuditModelsToDTOs(models []*domain.DataAudit) []*dto.DataAuditDTO {
	res := make([]*dto.DataAuditDTO, 0, len(models))

	for _, model := range models {
		data := MapDataAuditModelToDTO(model)
		if data != nil {
			res = append(res, data)
		}

	}

	return res
}
