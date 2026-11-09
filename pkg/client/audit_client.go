package client

import (
	"context"
)

// AuditAction определяет сигнатуру функции-обработчика, которая будет выполняться
// воркером асинхронного пула при отправке очередной записи аудит-лога.
type AuditAction[D any] func(ctx context.Context, workerIndex int, data D, token string) error

// AuditClient описывает контракт для асинхронного клиента отправки аудит-логов,
// инкапсулирующего в себе управление жизненным циклом фоновых воркеров и подсчет потерь.
type AuditClient[A any] interface {
	// Start запускает внутренний пул воркеров для асинхронной обработки и отправки логов.
	Start(ctx context.Context) error

	// Stop корректно завершает работу пула воркеров, ожидая сброса оставшихся в очереди логов.
	Stop(ctx context.Context) error

	// Audit ставит запись аудита в очередь на асинхронную отправку.
	Audit(ctx context.Context, data A) error

	// TotalLost возвращает суммарное количество логов, которые были утеряны (например, из-за переполнения очереди).
	TotalLost() int32
}
