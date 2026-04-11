package rest

import (
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/models"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

func MapAuthDtoSDKToRest(authAudit *dto.AuthAuditDTO) *models.AuthAuditDTO {
	if authAudit == nil {
		return nil
	}

	return &models.AuthAuditDTO{
		Source:       authAudit.Source,
		EventDate:    authAudit.EventDate.Format(time.RFC3339),
		Event:        authAudit.Event,
		Status:       authAudit.Status,
		RequestID:    authAudit.RequestID,
		TraceID:      authAudit.TraceID,
		Username:     authAudit.Username,
		AccessToken:  authAudit.AccessToken,
		RefreshToken: authAudit.RefreshToken,
	}
}

func MapDataDtoSDKToRest(dataAudit *dto.DataAuditDTO) *models.DataAuditDTO {
	if dataAudit == nil {
		return nil
	}

	return &models.DataAuditDTO{
		Source:           dataAudit.Source,
		EventDate:        dataAudit.EventDate.Format(time.RFC3339),
		Event:            dataAudit.Event,
		Status:           dataAudit.Status,
		RequestID:        dataAudit.RequestID,
		TraceID:          dataAudit.TraceID,
		Username:         dataAudit.Username,
		InternalTypeName: dataAudit.InternalTypeName,
		TypeName:         dataAudit.TypeName,
		TypeDescription:  dataAudit.TypeDescription,
		InstanceID:       dataAudit.InstanceID,
		InstanceName:     dataAudit.InstanceName,
		Values:           MapDataValueDTOsSDKToRest(dataAudit.Values),
	}
}

func MapDataValueDtoSDKToRest(dataAuditValue *dto.DataAuditValueDTO) *models.DataAuditValueDTO {
	if dataAuditValue == nil {
		return nil
	}

	return &models.DataAuditValueDTO{
		Name:        dataAuditValue.Name,
		Description: dataAuditValue.Description,
		Before:      dataAuditValue.Before,
		After:       dataAuditValue.After,
	}
}

func MapDataValueDTOsSDKToRest(dataAuditValues []*dto.DataAuditValueDTO) []*models.DataAuditValueDTO {
	if len(dataAuditValues) == 0 {
		return nil
	}

	res := make([]*models.DataAuditValueDTO, 0, len(dataAuditValues))
	for _, dataValue := range dataAuditValues {
		res = append(res, MapDataValueDtoSDKToRest(dataValue))
	}

	return res
}
