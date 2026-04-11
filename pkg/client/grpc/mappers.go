package grpc

import (
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAuthDtoSDKToGRPC(authData *dto.AuthAuditDTO) *pb.AuthAudit {
	if authData == nil {
		return nil
	}

	return pb.AuthAudit_builder{
		Source:       &authData.Source,
		EventDate:    timestamppb.New(authData.EventDate),
		Event:        &authData.Event,
		Status:       &authData.Status,
		RequestId:    &authData.RequestID,
		TraceId:      &authData.TraceID,
		Username:     &authData.Username,
		AccessToken:  &authData.AccessToken,
		RefreshToken: &authData.RefreshToken,
	}.Build()
}

func MapDataDtoSDKToGRPC(data *dto.DataAuditDTO) *pb.DataAudit {
	if data == nil {
		return nil
	}

	return pb.DataAudit_builder{
		Source:           &data.Source,
		EventDate:        timestamppb.New(data.EventDate),
		Event:            &data.Event,
		Status:           &data.Status,
		RequestId:        &data.RequestID,
		TraceId:          &data.TraceID,
		Username:         &data.Username,
		InternalTypeName: &data.InternalTypeName,
		TypeName:         &data.TypeName,
		TypeDescription:  &data.TypeDescription,
		InstanceId:       &data.InstanceID,
		InstanceName:     &data.InstanceName,
		Values:           MapDataValueDTOsSDKToGRPC(data.Values),
	}.Build()
}

func MapDataValueDtoSDKToGRPC(dataValue *dto.DataAuditValueDTO) *pb.DataAuditValue {
	if dataValue == nil {
		return nil
	}

	return pb.DataAuditValue_builder{
		Name:        &dataValue.Name,
		Description: &dataValue.Description,
		Before:      &dataValue.Before,
		After:       &dataValue.After,
	}.Build()
}

func MapDataValueDTOsSDKToGRPC(dataValues []*dto.DataAuditValueDTO) []*pb.DataAuditValue {
	if dataValues == nil {
		return nil
	}

	res := make([]*pb.DataAuditValue, 0, len(dataValues))
	for _, dataValue := range dataValues {
		res = append(res, MapDataValueDtoSDKToGRPC(dataValue))
	}

	return res
}
