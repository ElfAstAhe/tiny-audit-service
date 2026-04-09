package grpc

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapDataAuditValueDTOToGRPC(dataAuditValue *dto.DataAuditValueDTO) *pb.DataAuditValue {
	if dataAuditValue == nil {
		return nil
	}

	return pb.DataAuditValue_builder{
		Name:        &dataAuditValue.Name,
		Description: &dataAuditValue.Description,
		Before:      &dataAuditValue.Before,
		After:       &dataAuditValue.After,
	}.Build()
}

func MapDataAuditValueGRPCToDTO(dataAuditValue *pb.DataAuditValue) *dto.DataAuditValueDTO {
	if dataAuditValue == nil {
		return nil
	}

	return &dto.DataAuditValueDTO{
		Name:        dataAuditValue.GetName(),
		Description: dataAuditValue.GetDescription(),
		Before:      dataAuditValue.GetBefore(),
		After:       dataAuditValue.GetAfter(),
	}
}

func MapDataAuditValuesGRPCToDTO(dataAuditValues []*pb.DataAuditValue) []*dto.DataAuditValueDTO {
	if dataAuditValues == nil {
		return make([]*dto.DataAuditValueDTO, 0)
	}

	res := make([]*dto.DataAuditValueDTO, 0, len(dataAuditValues))
	for _, data := range dataAuditValues {
		res = append(res, MapDataAuditValueGRPCToDTO(data))
	}

	return res
}

func MapDataAuditValuesDTOToGRPC(dataAuditValues []*dto.DataAuditValueDTO) []*pb.DataAuditValue {
	if len(dataAuditValues) == 0 {
		return make([]*pb.DataAuditValue, 0)
	}
	res := make([]*pb.DataAuditValue, 0, len(dataAuditValues))
	for _, data := range dataAuditValues {
		res = append(res, MapDataAuditValueDTOToGRPC(data))
	}

	return res
}

func MapDataAuditDTOToGRPC(dataAudit *dto.DataAuditDTO) *pb.DataAudit {
	if dataAudit == nil {
		return nil
	}

	return pb.DataAudit_builder{
		Id:               &dataAudit.ID,
		Source:           &dataAudit.Source,
		EventDate:        timestamppb.New(dataAudit.EventDate),
		Event:            &dataAudit.Event,
		Status:           &dataAudit.Status,
		RequestId:        &dataAudit.RequestID,
		TraceId:          &dataAudit.TraceID,
		Username:         &dataAudit.Username,
		InternalTypeName: &dataAudit.InternalTypeName,
		TypeName:         &dataAudit.TypeName,
		TypeDescription:  &dataAudit.TypeDescription,
		InstanceId:       &dataAudit.InstanceID,
		InstanceName:     &dataAudit.InstanceName,
		Values:           MapDataAuditValuesDTOToGRPC(dataAudit.Values),
		CreateAt:         timestamppb.New(dataAudit.CreatedAt),
		UpdateAt:         timestamppb.New(dataAudit.UpdatedAt),
	}.Build()
}

func MapDataAuditsDTOToGRPC(dataAudits []*dto.DataAuditDTO) []*pb.DataAudit {
	if len(dataAudits) == 0 {
		return make([]*pb.DataAudit, 0)
	}

	res := make([]*pb.DataAudit, 0, len(dataAudits))
	for _, data := range dataAudits {
		res = append(res, MapDataAuditDTOToGRPC(data))
	}

	return res
}

func MapDataAuditGRPCToDTO(dataAudit *pb.DataAudit) *dto.DataAuditDTO {
	if dataAudit == nil {
		return nil
	}

	return &dto.DataAuditDTO{
		ID:               dataAudit.GetId(),
		Source:           dataAudit.GetSource(),
		EventDate:        dataAudit.GetEventDate().AsTime(),
		Event:            dataAudit.GetEvent(),
		Status:           dataAudit.GetStatus(),
		RequestID:        dataAudit.GetRequestId(),
		TraceID:          dataAudit.GetTraceId(),
		Username:         dataAudit.GetUsername(),
		InternalTypeName: dataAudit.GetInternalTypeName(),
		TypeName:         dataAudit.GetTypeName(),
		TypeDescription:  dataAudit.GetTypeDescription(),
		InstanceID:       dataAudit.GetInstanceId(),
		InstanceName:     dataAudit.GetInstanceName(),
		Values:           MapDataAuditValuesGRPCToDTO(dataAudit.GetValues()),
		CreatedAt:        dataAudit.GetCreateAt().AsTime(),
		UpdatedAt:        dataAudit.GetUpdateAt().AsTime(),
	}
}

func MapAuditInstanceGRPCToDTO(req *pb.ListByInstanceRequest) *dto.AuditInstanceDTO {
	if req == nil {
		return nil
	}

	return &dto.AuditInstanceDTO{
		TypeName:   req.GetTypeName(),
		InstanceID: req.GetInstanceId(),
		Limit:      int(req.GetLimit()),
		Offset:     int(req.GetOffset()),
	}
}
