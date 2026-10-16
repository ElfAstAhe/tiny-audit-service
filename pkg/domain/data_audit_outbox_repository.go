package domain

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

// DataAuditOutboxRepository определяет специализированный интерфейс доменного уровня
// для управления персистентным состоянием записей транзакционного журнала Outbox.
//
// Расширяет базовый контракт CRUDRepository фреймворка, связывая указатель на доменную
// сущность *DataAuditOutbox с текстовым типом первичного ключа (string).
type DataAuditOutboxRepository interface {
	domain.CRUDRepository[*DataAuditOutbox, string]

	// ListPending извлекает из хранилища пачку неотправленных логов аудита (пакетный режим Batching).
	ListPending(ctx context.Context, limit int) ([]*DataAuditOutbox, error)
}
