package grpc

import (
	"context"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataAuditGRPCService struct {
	pb.UnimplementedDataAuditServiceServer
	dataAuditFacade facade.DataAuditFacade
}

func NewDataAuditGRPCService(dataAuditFacade facade.DataAuditFacade) *DataAuditGRPCService {
	return &DataAuditGRPCService{
		dataAuditFacade: dataAuditFacade,
	}
}

func (das *DataAuditGRPCService) Audit(ctx context.Context, req *pb.DataAuditRequest) (*emptypb.Empty, error) {
	// ToDo: implement

	return &emptypb.Empty{}, nil
}

func (das *DataAuditGRPCService) ListByPeriod(ctx context.Context, req *pb.DataListByPeriod) (*pb.DataAuditInstances, error) {
	// ToDo: implement

	return nil, nil
}

func (das *DataAuditGRPCService) ListByInstance(ctx context.Context, req *pb.DataListByInstance) (*pb.DataAuditInstances, error) {
	// ToDo: implement

	return nil, nil
}
