package grpc

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
)

// MapAuditPeriodGRPCToDTO транслирует входящий gRPC-запрос pb.ListByPeriodRequest в объект параметров dto.AuditPeriodDTO, преобразуя protobuf-таймстампы в стандартный time.Time.
func MapAuditPeriodGRPCToDTO(auditPeriod *pb.ListByPeriodRequest) *dto.AuditPeriodDTO {
	if auditPeriod == nil {
		return nil
	}

	return &dto.AuditPeriodDTO{
		From:   auditPeriod.GetFrom().AsTime(),
		Till:   auditPeriod.GetTill().AsTime(),
		Limit:  (int)(auditPeriod.GetLimit()),
		Offset: (int)(auditPeriod.GetOffset()),
	}
}
