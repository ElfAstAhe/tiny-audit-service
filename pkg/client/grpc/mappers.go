package grpc

import (
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MapAuthDtoSDKToGRPC transforms a client-level AuthAuditDTO data contract into a strongly-typed wire-ready gRPC pb.AuthAudit message using the Opaque API builder.
func MapAuthDtoSDKToGRPC(authData *dto.AuthAuditDTO) *pb.AuthAudit {
	if utils.IsNil(authData) {
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

// MapDataDtoSDKToGRPC converts an inbound client-level DataAuditDTO data payload schema into a decoupled network delivery pb.DataAudit protobuf aggregate.
func MapDataDtoSDKToGRPC(data *dto.DataAuditDTO) *pb.DataAudit {
	if utils.IsNil(data) {
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

// MapDataValueDtoSDKToGRPC maps an internal data modification record field variant schema onto a structured network pb.DataAuditValue element payload.
func MapDataValueDtoSDKToGRPC(dataValue *dto.DataAuditValueDTO) *pb.DataAuditValue {
	if utils.IsNil(dataValue) {
		return nil
	}

	return pb.DataAuditValue_builder{
		Name:        &dataValue.Name,
		Description: &dataValue.Description,
		Before:      &dataValue.Before,
		After:       &dataValue.After,
	}.Build()
}

// MapDataValueDTOsSDKToGRPC processes a slice array of client change records, converting and allocating them into a sequential gRPC list collection array.
func MapDataValueDTOsSDKToGRPC(dataValues []*dto.DataAuditValueDTO) []*pb.DataAuditValue {
	if len(dataValues) == 0 {
		return nil
	}

	res := make([]*pb.DataAuditValue, 0, len(dataValues))
	for _, dataValue := range dataValues {
		res = append(res, MapDataValueDtoSDKToGRPC(dataValue))
	}

	return res
}
