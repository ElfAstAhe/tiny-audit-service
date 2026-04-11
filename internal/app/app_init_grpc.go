package app

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	grpcsvc "github.com/ElfAstAhe/tiny-audit-service/internal/transport/grpc"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/grpc/interceptor"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	pkggrpcintercept "github.com/ElfAstAhe/tiny-audit-service/pkg/transport/interceptor"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

func (app *App) initGRPCService() error {
	app.grpcAuthAuditService = grpcsvc.NewAuthAuditGRPCService(app.authFacade)
	app.grpcDataAuditService = grpcsvc.NewDataAuditGRPCService(app.dataFacade)

	return nil
}

//goland:noinspection DuplicatedCode
func (app *App) initGRPCServer() error {
	// Настраиваем KeepAlive на основе твоего GRPCConfig
	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     app.config.GRPC.MaxConnIdle,
		MaxConnectionAge:      app.config.GRPC.MaxConnAge,
		MaxConnectionAgeGrace: app.config.GRPC.MaxConnAgeGrace,
		Time:                  app.config.GRPC.KeepAliveTime,
		Timeout:               app.config.GRPC.KeepAliveTimeout,
	}
	// Метрики
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
		// Add tenant_name as a context label. This server option is necessary
		// to initialize the metrics with the labels that will be provided
		// dynamically from the context. This should be used in tandem with
		// WithLabelsFromContext in the interceptor options.
		grpcprom.WithContextLabels("tenant_name"),
	)
	if err := prometheus.Register(srvMetrics); err != nil {
		return errs.NewCommonError("failed to register grpc metrics", err)
	}
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}

		return nil
	}
	// Extract the tenant name value from gRPC metadata
	// and use it as a label on our metrics.
	labelsFromContext := func(ctx context.Context) prometheus.Labels {
		labels := prometheus.Labels{}

		md := metadata.ExtractIncoming(ctx)
		tenantName := md.Get("tenant-name")
		if tenantName == "" {
			tenantName = "unknown"
		}
		labels["tenant_name"] = tenantName

		return labels
	}
	// Setup metric for panic recoveries.
	panicsTotal := promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		//		rpcLogger.Error("recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	// real IP
	realIPOpts := []realip.Option{
		realip.WithHeaders([]string{
			realip.XRealIp,
			realip.XForwardedFor,
			realip.TrueClientIp,
		}),
	}

	authExtractor := interceptor.NewAuthExtractor(
		app.jwtHelper,
		app.jwtGRPCHelper,
		app.authHelper,
		app.logger,
		[]string{},
		app.config.App.AcceptTokenIssuers,
	)
	// Собираем опции сервера
	opts := []grpc.ServerOption{
		// keepalive
		grpc.KeepaliveParams(kasp),
		// timeout
		grpc.ConnectionTimeout(app.config.GRPC.Timeout),
		// tracing
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		// metrics
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(
				grpcprom.WithExemplarFromContext(exemplarFromContext),
				grpcprom.WithLabelsFromContext(labelsFromContext),
			),
			pkggrpcintercept.AuditRequestIDExtractorUnaryServerInterceptor([]string{
				pkggrpcintercept.MDXRequestID,
				pkggrpcintercept.MDXCorrelationID,
				pkggrpcintercept.MDRequestID,
			}),
			pkggrpcintercept.AuditTraceIDExtractorUnaryServerInterceptor([]string{
				pkggrpcintercept.MDXCloudTraceContext,
				pkggrpcintercept.MDTraceParent,
				pkggrpcintercept.MDXTraceID,
				pkggrpcintercept.MDTraceID,
			}),
			realip.UnaryServerInterceptorOpts(realIPOpts...),
			authExtractor.UnaryServerInterceptor,
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(
				grpcprom.WithExemplarFromContext(exemplarFromContext),
				grpcprom.WithLabelsFromContext(labelsFromContext),
			),
			pkggrpcintercept.AuditRequestIDExtractorStreamServerInterceptor([]string{
				pkggrpcintercept.MDXRequestID,
				pkggrpcintercept.MDXCorrelationID,
				pkggrpcintercept.MDRequestID,
			}),
			pkggrpcintercept.AuditTraceIDExtractorStreamServerInterceptor([]string{
				pkggrpcintercept.MDXCloudTraceContext,
				pkggrpcintercept.MDTraceParent,
				pkggrpcintercept.MDXTraceID,
				pkggrpcintercept.MDTraceID,
			}),
			realip.StreamServerInterceptorOpts(realIPOpts...),
			authExtractor.StreamServerInterceptor,
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	}

	app.grpcServer = grpc.NewServer(opts...)

	// Регистрация
	pb.RegisterAuthAuditServiceServer(app.grpcServer, app.grpcAuthAuditService)
	pb.RegisterDataAuditServiceServer(app.grpcServer, app.grpcDataAuditService)

	// Инициализация метрик с нулевыми рядами
	srvMetrics.InitializeMetrics(app.grpcServer)

	return nil
}
