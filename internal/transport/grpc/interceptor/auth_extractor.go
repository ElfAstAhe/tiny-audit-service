package interceptor

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Вспомогательная структура для подмены контекста в стриме
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

type AuthExtractor struct {
	authHelper auth.Helper
	log        logger.Logger
	nonSecure  map[string]struct{}
}

func NewAuthExtractor(authHelper auth.Helper, logger logger.Logger) *AuthExtractor {
	return &AuthExtractor{
		authHelper: authHelper,
		log:        logger.GetLogger("gRPC-Auth-Extractor"),
		nonSecure:  map[string]struct{}{},
	}
}

func (ae *AuthExtractor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ae.log.Debugf("UnaryServerInterceptor start with req: [%v]", req)
	defer ae.log.Debug("UnaryServerInterceptor finish")

	// ignorance
	if ae.isNonSecure(info.FullMethod) {
		return handler(ctx, req)
	}

	subj, err := ae.authHelper.SubjectFromGRPCContext(ctx)
	if err != nil {
		ae.log.Errorf("AuthExtractor.UnaryServerInterceptor failed with error: [%v]", err)

		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	secureCtx := auth.WithSubject(ctx, subj)

	return handler(secureCtx, req)
}

func (ae *AuthExtractor) StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ae.log.Debugf("StreamServerInterceptor start: method [%v]", info.FullMethod)
	defer ae.log.Debugf("StreamServerInterceptor finish: method [%v]", info.FullMethod)

	// ignorance
	if ae.isNonSecure(info.FullMethod) {
		return handler(srv, stream)
	}

	subj, err := ae.authHelper.SubjectFromGRPCContext(stream.Context())
	if err != nil {
		ae.log.Errorf("AuthExtractor.StreamServerInterceptor failed with error: [%v]", err)

		return status.Error(codes.Unauthenticated, err.Error())
	}

	wrapped := &wrappedStream{
		ServerStream: stream,
		ctx:          auth.WithSubject(stream.Context(), subj),
	}

	return handler(srv, wrapped)
}

func (ae *AuthExtractor) isNonSecure(method string) bool {
	_, ok := ae.nonSecure[method]

	return ok
}
