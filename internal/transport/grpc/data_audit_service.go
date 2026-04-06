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
	err := das.dataAuditFacade.Audit(ctx, MapDataAuditGRPCToDTO(req.GetData()))
	if err != nil {
		return nil, MapToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (das *DataAuditGRPCService) ListByPeriod(ctx context.Context, req *pb.ListByPeriodRequest) (*pb.DataAuditInstances, error) {
	res, err := das.dataAuditFacade.ListByPeriod(ctx, MapAuditPeriodGRPCToDTO(req))
	if err != nil {
		return nil, MapToGrpcError(err)
	}

	return pb.DataAuditInstances_builder{
		Instances: MapDataAuditsDTOToGRPC(res),
	}.Build(), nil
}

func (das *DataAuditGRPCService) ListByInstance(ctx context.Context, req *pb.ListByInstanceRequest) (*pb.DataAuditInstances, error) {
	res, err := das.dataAuditFacade.ListByInstance(ctx, MapAuditInstanceGRPCToDTO(req))
	if err != nil {
		return nil, MapToGrpcError(err)
	}

	return pb.DataAuditInstances_builder{
		Instances: MapDataAuditsDTOToGRPC(res),
	}.Build(), nil
}
