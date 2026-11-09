package grpc

import (
	"context"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthAuditGRPCService реализует сгенерированный gRPC-сервер pb.AuthAuditServiceServer,
// обрабатывающий запросы регистрации и выборки архивных логов аудита аутентификации.
type AuthAuditGRPCService struct {
	pb.UnimplementedAuthAuditServiceServer
	authAuditFacade facade.AuthAudit
}

var _ pb.AuthAuditServiceServer = (*AuthAuditGRPCService)(nil)

// NewAuthAuditGRPCService создает и инициализирует новый экземпляр gRPC-сервиса AuthAuditGRPCService.
func NewAuthAuditGRPCService(authAuditFacade facade.AuthAudit) *AuthAuditGRPCService {
	return &AuthAuditGRPCService{
		authAuditFacade: authAuditFacade,
	}
}

// Audit принимает gRPC-запрос на регистрацию аудит-лога аутентификации, выполняет маппинг в DTO и передает управление на слой фасада.
func (aas *AuthAuditGRPCService) Audit(ctx context.Context, req *pb.AuthAuditRequest) (*emptypb.Empty, error) {
	err := aas.authAuditFacade.Audit(ctx, MapAuthAuditGRPCToDTO(req.GetData()))
	if err != nil {
		return nil, mapToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

// ListByPeriod возвращает список записей аудита событий аутентификации за указанный в gRPC-запросе временной интервал.
func (aas *AuthAuditGRPCService) ListByPeriod(ctx context.Context, req *pb.ListByPeriodRequest) (*pb.AuthAuditInstances, error) {
	dtoRes, err := aas.authAuditFacade.ListByPeriod(ctx, MapAuditPeriodGRPCToDTO(req))
	if err != nil {
		return nil, mapToGrpcError(err)
	}

	return pb.AuthAuditInstances_builder{
		Instances: MapAuthAuditDTOsToGRPCs(dtoRes),
	}.Build(), nil
}

// ListByUsername возвращает список записей аудита аутентификации для конкретного пользователя, переданного в gRPC-запросе.
func (aas *AuthAuditGRPCService) ListByUsername(ctx context.Context, req *pb.AuthListByUsernameRequest) (*pb.AuthAuditInstances, error) {
	dtoRes, err := aas.authAuditFacade.ListByUsername(ctx, MapAuthAuditUserGRPCToDTO(req))
	if err != nil {
		return nil, mapToGrpcError(err)
	}

	return pb.AuthAuditInstances_builder{
		Instances: MapAuthAuditDTOsToGRPCs(dtoRes),
	}.Build(), nil
}
