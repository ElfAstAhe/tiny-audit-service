package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// serverStream переопределяет интерфейс grpc.ServerStream, позволяя
// безопасно подменять встроенный context.Context дочерних потоков.
type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context возвращает текущий (в том числе кастомный или обогащенный) контекст gRPC-стрима.
func (s *serverStream) Context() context.Context {
	return s.ctx
}
