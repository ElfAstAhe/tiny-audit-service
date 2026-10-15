package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuditorService представляет собой обобщенный (Generic) интерфейс верхнего уровня
// для абстрагирования автогенерированных gRPC-клиентов службы аудита.
//
// Данный интерфейс позволяет изолировать прикладной слой SDK от конкретных реализаций
// Protobuf-сообщений (таких как AuthAuditRequest или DataAuditRequest) и упрощает процесс
// подмены сетевого gRPC-транспорта на mock-заглушки в изолированных юнит-тестах.
type AuditorService[D any] interface {
	// Audit осуществляет блокирующий удаленный вызов RPC-процедуры на сервере аудита.
	// Принимает типизированную структуру сообщения и опциональный вариативный срез опций grpc.CallOption.
	Audit(ctx context.Context, data D, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
