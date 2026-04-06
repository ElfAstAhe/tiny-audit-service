package grpc

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAuthAuditGRPCToDTO(authAudit *pb.AuthAudit) *dto.AuthAuditDTO {
	if authAudit == nil {
		return nil
	}

	return &dto.AuthAuditDTO{
		ID:           authAudit.GetId(),
		Source:       authAudit.GetSource(),
		EventDate:    authAudit.GetEventDate().AsTime(),
		Event:        authAudit.GetEvent(),
		Status:       authAudit.GetStatus(),
		RequestID:    authAudit.GetRequestId(),
		TraceID:      authAudit.GetTraceId(),
		Username:     authAudit.GetUsername(),
		AccessToken:  authAudit.GetAccessToken(),
		RefreshToken: authAudit.GetRefreshToken(),
		CreatedAt:    authAudit.GetCreateAt().AsTime(),
		UpdatedAt:    authAudit.GetUpdateAt().AsTime(),
	}
}

func MapAuthAuditDTOToGRPC(authAudit *dto.AuthAuditDTO) *pb.AuthAudit {
	if authAudit == nil {
		return nil
	}

	return pb.AuthAudit_builder{
		Id:           &authAudit.ID,
		Source:       &authAudit.Source,
		EventDate:    timestamppb.New(authAudit.EventDate),
		Event:        &authAudit.Event,
		Status:       &authAudit.Status,
		RequestId:    &authAudit.RequestID,
		TraceId:      &authAudit.TraceID,
		Username:     &authAudit.Username,
		AccessToken:  &authAudit.AccessToken,
		RefreshToken: &authAudit.RefreshToken,
		CreateAt:     timestamppb.New(authAudit.CreatedAt),
		UpdateAt:     timestamppb.New(authAudit.UpdatedAt),
	}.Build()
}

func MapAuthAuditDTOsToGRPCs(authAudits []*dto.AuthAuditDTO) []*pb.AuthAudit {
	res := make([]*pb.AuthAudit, 0, len(authAudits))

	for _, authAudit := range authAudits {
		res = append(res, MapAuthAuditDTOToGRPC(authAudit))
	}

	return res
}

func MapAuthAuditUserGRPCToDTO(auditUser *pb.AuthListByUsernameRequest) *dto.AuditUserDTO {
	if auditUser == nil {
		return nil
	}

	return &dto.AuditUserDTO{
		Username: auditUser.GetUsername(),
		Limit:    int(auditUser.GetLimit()),
		Offset:   int(auditUser.GetOffset()),
	}
}
