package grpc

import (
	"context"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthAuditGRPCService struct {
	pb.UnimplementedAuthAuditServiceServer
	authAuditFacade facade.AuthAuditFacade
}

var _ pb.AuthAuditServiceServer = (*AuthAuditGRPCService)(nil)

func NewAuthAuditGRPCService(authAuditFacade facade.AuthAuditFacade) *AuthAuditGRPCService {
	return &AuthAuditGRPCService{
		authAuditFacade: authAuditFacade,
	}
}

func (aas *AuthAuditGRPCService) Audit(ctx context.Context, req *pb.AuthAuditRequest) (*emptypb.Empty, error) {
	err := aas.authAuditFacade.Audit(ctx, MapAuthAuditGRPCToDTO(req.GetData()))
	if err != nil {
		return nil, MapToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (aas *AuthAuditGRPCService) ListByPeriod(ctx context.Context, req *pb.ListByPeriodRequest) (*pb.AuthAuditInstances, error) {
	dtoRes, err := aas.authAuditFacade.ListByPeriod(ctx, MapAuditPeriodGRPCToDTO(req))
	if err != nil {
		return nil, MapToGrpcError(err)
	}

	return pb.AuthAuditInstances_builder{
		Instances: MapAuthAuditDTOsToGRPCs(dtoRes),
	}.Build(), nil
}

func (aas *AuthAuditGRPCService) ListByUsername(ctx context.Context, req *pb.AuthListByUsernameRequest) (*pb.AuthAuditInstances, error) {
	dtoRes, err := aas.authAuditFacade.ListByUsername(ctx, MapAuthAuditUserGRPCToDTO(req))
	if err != nil {
		return nil, MapToGrpcError(err)
	}

	return pb.AuthAuditInstances_builder{
		Instances: MapAuthAuditDTOsToGRPCs(dtoRes),
	}.Build(), nil
}
