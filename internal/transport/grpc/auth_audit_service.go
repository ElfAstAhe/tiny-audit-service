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
	// ToDo: implement

	return &emptypb.Empty{}, nil
}

func (aas *AuthAuditGRPCService) ListByPeriod(ctx context.Context, req *pb.AuthListByPeriodRequest) (*pb.AuthAuditInstances, error) {
	// ToDo: implement

	return nil, nil
}

func (aas *AuthAuditGRPCService) ListByUsername(ctx context.Context, req *pb.AuthListByUsernameRequest) (*pb.AuthAuditInstances, error) {
	// ToDo: implement

	return nil, nil
}
