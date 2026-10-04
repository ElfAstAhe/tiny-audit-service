package interceptor

import (
	"context"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// MDXCloudTraceContext defines the standard Google Cloud pseudo-header used for cross-cutting tracing propagation.
	MDXCloudTraceContext string = "x-cloud-trace-context"
	// MDTraceParent defines the W3C standard distributed tracing context header specification.
	MDTraceParent string = "traceparent"
	// MDXTraceID defines a common customized legacy corporate identification tracing span header.
	MDXTraceID string = "x-trace-id"
	// MDTraceID defines a flat simplified diagnostic sequence tracing key.
	MDTraceID string = "trace-id"
)

// AuditTraceIDExtractorUnaryServerInterceptor creates a unary gRPC server interceptor that scans incoming metadata
// headers to extract, fall back, and inject an isolated transaction context correlation trace identifier downstream.
func AuditTraceIDExtractorUnaryServerInterceptor(headers ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var traceID string
		for _, header := range headers {
			vals := metadata.ValueFromIncomingContext(ctx, header)
			if len(vals) > 0 {
				traceID = vals[0]
				break
			}
		}
		if traceID == "" {
			traceID = "unknown"
		}

		return handler(utils.WithTraceID(ctx, traceID), req)
	}
}

// AuditTraceIDExtractorStreamServerInterceptor creates a streaming gRPC server interceptor that extracts distributed correlation identifiers
// from metadata blocks and wraps server stream contexts before propagating the chunk execution loops down the stream line.
func AuditTraceIDExtractorStreamServerInterceptor(headers ...string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var traceID string
		for _, header := range headers {
			vals := metadata.ValueFromIncomingContext(ss.Context(), header)
			if len(vals) > 0 {
				traceID = vals[0]
				break
			}
		}
		if traceID == "" {
			traceID = "unknown"
		}

		return handler(srv, &serverStream{
			ServerStream: ss,
			ctx:          utils.WithTraceID(ss.Context(), traceID),
		})
	}
}
