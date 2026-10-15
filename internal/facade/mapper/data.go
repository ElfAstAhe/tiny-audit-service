package mapper

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

// MapDataAuditValueModelToDTO трансформирует внутреннюю доменную модель атрибута в соответствующий объект передачи данных (DTO).
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

// MapDataAuditValueModelsToDTOs веерно конвертирует срез доменных моделей полей в срез DTO структур, исключая пустые указатели.
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

// MapDataAuditValueDTOToModel трансформирует входящий объект передачи данных (DTO) атрибута во внутреннюю доменную модель.
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

// MapDataAuditValueDTOsToModels веерно конвертирует входящий срез DTO структур в срез внутренних доменных моделей полей.
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

// MapDataAuditModelToDTO преобразует корневой доменный агрегат лога аудита в плоскую структуру DTO с каскадным маппингом полей.
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

// MapDataAuditDTOToModel преобразует внешнюю структуру DTO в корневой доменный агрегат лога аудита для последующей бизнес-валидации.
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

// MapDataAuditModelsToDTOs выполняет массовую трансформацию пачки доменных моделей логов в срез DTO для отдачи в API-слой.
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
