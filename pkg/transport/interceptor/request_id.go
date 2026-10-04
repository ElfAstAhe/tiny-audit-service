package interceptor

import (
	"context"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// MDXRequestID defines the canonical case-insensitive metadata key used to extract upstream request sequences.
	MDXRequestID string = "x-request-id"
	// MDXCorrelationID defines a standard enterprise tracking metadata key for distributed cross-system operations.
	MDXCorrelationID string = "x-correlation-id"
	// MDRequestID duplicate target map key mapping diagnostic request sequences parameters indicators.
	MDRequestID string = "x-request-id"
)

// AuditRequestIDExtractorUnaryServerInterceptor creates a unary gRPC server interceptor that scans incoming metadata
// headers to extract, fall back, and inject an isolated transaction context correlation request identifier downstream.
func AuditRequestIDExtractorUnaryServerInterceptor(headers ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var requestID string
		for _, header := range headers {
			vals := metadata.ValueFromIncomingContext(ctx, header)
			if len(vals) > 0 {
				requestID = vals[0]
				break
			}
		}
		if requestID == "" {
			requestID = "unknown"
		}

		return handler(utils.WithRequestID(ctx, requestID), req)
	}
}

// AuditRequestIDExtractorStreamServerInterceptor creates a streaming gRPC server interceptor that extracts distributed correlation identifiers
// from metadata blocks and wraps server stream contexts before propagating the chunk execution loops down the stream line.
func AuditRequestIDExtractorStreamServerInterceptor(headers ...string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var requestID string
		for _, header := range headers {
			vals := metadata.ValueFromIncomingContext(ss.Context(), header)
			if len(vals) > 0 {
				requestID = vals[0]
				break
			}
		}
		if requestID == "" {
			requestID = "unknown"
		}

		return handler(srv, &serverStream{
			ServerStream: ss,
			ctx:          utils.WithRequestID(ss.Context(), requestID),
		})
	}
}
