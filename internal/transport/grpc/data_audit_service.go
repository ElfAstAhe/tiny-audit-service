package grpc

import (
	"context"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DataAuditGRPCService реализует сгенерированный gRPC-сервер pb.DataAuditServiceServer,
// обрабатывающий запросы регистрации и архивных выборок логов аудита изменений данных.
type DataAuditGRPCService struct {
	pb.UnimplementedDataAuditServiceServer
	dataAuditFacade facade.DataAudit
}

// NewDataAuditGRPCService создает и инициализирует новый экземпляр gRPC-сервиса DataAuditGRPCService.
func NewDataAuditGRPCService(dataAuditFacade facade.DataAudit) *DataAuditGRPCService {
	return &DataAuditGRPCService{
		dataAuditFacade: dataAuditFacade,
	}
}

// Audit принимает gRPC-запрос на регистрацию аудит-лога данных, выполняет маппинг в DTO и передает управление на слой фасада.
func (das *DataAuditGRPCService) Audit(ctx context.Context, req *pb.DataAuditRequest) (*emptypb.Empty, error) {
	err := das.dataAuditFacade.Audit(ctx, MapDataAuditGRPCToDTO(req.GetData()))
	if err != nil {
		return nil, mapToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

// ListByPeriod возвращает список записей аудита изменений данных за указанный в gRPC-запросе временной интервал.
func (das *DataAuditGRPCService) ListByPeriod(ctx context.Context, req *pb.ListByPeriodRequest) (*pb.DataAuditInstances, error) {
	res, err := das.dataAuditFacade.ListByPeriod(ctx, MapAuditPeriodGRPCToDTO(req))
	if err != nil {
		return nil, mapToGrpcError(err)
	}

	return pb.DataAuditInstances_builder{
		Instances: MapDataAuditsDTOToGRPC(res),
	}.Build(), nil
}

// ListByInstance возвращает список записей аудита для конкретного экземпляра сущности по её типу и уникальному идентификатору.
func (das *DataAuditGRPCService) ListByInstance(ctx context.Context, req *pb.ListByInstanceRequest) (*pb.DataAuditInstances, error) {
	res, err := das.dataAuditFacade.ListByInstance(ctx, MapAuditInstanceGRPCToDTO(req))
	if err != nil {
		return nil, mapToGrpcError(err)
	}

	return pb.DataAuditInstances_builder{
		Instances: MapDataAuditsDTOToGRPC(res),
	}.Build(), nil
}
