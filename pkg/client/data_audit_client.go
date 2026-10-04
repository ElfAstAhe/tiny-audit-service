package client

import (
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

// DataAuditClient defines the dedicated high-performance domain integration contract
// managing data mutation logging operations utilizing strongly-typed DataAuditDTO payloads.
type DataAuditClient interface {
	AuditClient[*dto.DataAuditDTO]
}
