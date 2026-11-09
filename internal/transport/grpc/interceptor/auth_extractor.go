package interceptor

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// wrappedStream переопределяет интерфейс grpc.ServerStream, позволяя
// безопасно подменять встроенный context.Context дочерних gRPC-потоков.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context возвращает текущий (в том числе кастомный или обогащенный) контекст gRPC-стрима.
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// AuthExtractor представляет собой интерцептор для gRPC-сервера, выполняющий
// извлечение, валидацию JWT-токенов из метаданных контекста и авторизацию субъектов.
type AuthExtractor struct {
	jwtHelper     *helper.JWTHelper
	jwtGRPCHelper *helper.JWTGRPCHelper
	authHelper    auth.Helper
	log           logger.Logger
	nonSecure     map[string]struct{}
	acceptIssuers map[string]struct{}
}

// NewAuthExtractor создает и инициализирует новый экземпляр gRPC-интерцептора AuthExtractor,
// формируя хэш-карты для быстрой проверки белых списков публичных методов и доверенных эмитентов.
func NewAuthExtractor(
	jwtHelper *helper.JWTHelper,
	jwtGRPCHelper *helper.JWTGRPCHelper,
	authHelper auth.Helper,
	logger logger.Logger,
	nonSecureMethods []string,
	acceptIssuers []string,
) *AuthExtractor {
	nonSecureMap := make(map[string]struct{}, len(nonSecureMethods))
	for _, nonSecureMethod := range nonSecureMethods {
		nonSecureMap[nonSecureMethod] = struct{}{}
	}
	acceptIssuersMap := make(map[string]struct{}, len(acceptIssuers))
	for _, acceptIssuer := range acceptIssuers {
		acceptIssuersMap[acceptIssuer] = struct{}{}
	}

	return &AuthExtractor{
		jwtHelper:     jwtHelper,
		jwtGRPCHelper: jwtGRPCHelper,
		authHelper:    authHelper,
		log:           logger.GetLogger("gRPC-Auth-Extractor"),
		nonSecure:     nonSecureMap,
		acceptIssuers: acceptIssuersMap,
	}
}

// UnaryServerInterceptor осуществляет перехват унарных gRPC-вызовов для извлечения JWT-токена,
// проверки его эмитента и последующего внедрения авторизованного субъекта в контекст запроса.
func (ae *AuthExtractor) UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ae.log.Debugf("UnaryServerInterceptor start with req: [%v]", req)
	defer ae.log.Debug("UnaryServerInterceptor finish")

	// check nonsecure methods
	if ae.isNonSecure(info.FullMethod) {
		return handler(ctx, req)
	}

	// extract token
	token, err := ae.jwtGRPCHelper.ExtractTokenFromContext(ctx, auth.DefaultMetadataName)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	// convert into claims
	claims, err := ae.jwtHelper.ExtractClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	// check issuers via white list
	if !ae.isAcceptIssuer(claims.Issuer) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid issuer: %v", claims.Issuer)
	}

	subj, err := ae.authHelper.SubjectFromToken(token)
	if err != nil {
		ae.log.Errorf("AuthExtractor.UnaryServerInterceptor failed with error: [%v]", err)

		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	secureCtx := auth.WithSubject(ctx, subj)

	return handler(secureCtx, req)
}

// StreamServerInterceptor осуществляет перехват потоковых (streaming) gRPC-вызовов, проверяя
// авторизационные метаданные и оборачивая исходный поток для подмены контекста авторизованным субъектом.
func (ae *AuthExtractor) StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ae.log.Debugf("StreamServerInterceptor start: method [%v]", info.FullMethod)
	defer ae.log.Debugf("StreamServerInterceptor finish: method [%v]", info.FullMethod)

	// ignorance
	if ae.isNonSecure(info.FullMethod) {
		return handler(srv, stream)
	}

	// extract token
	token, err := ae.jwtGRPCHelper.ExtractTokenFromContext(stream.Context(), auth.DefaultMetadataName)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	// convert into claims
	claims, err := ae.jwtHelper.ExtractClaims(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	// check issuers via white list
	if !ae.isAcceptIssuer(claims.Issuer) {
		return status.Errorf(codes.Unauthenticated, "invalid issuer: %v", claims.Issuer)
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

func (ae *AuthExtractor) isAcceptIssuer(issuer string) bool {
	_, ok := ae.acceptIssuers[issuer]

	return ok
}
