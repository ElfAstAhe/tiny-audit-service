package rest

import (
	"net/http"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	libconf "github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libhttp "github.com/ElfAstAhe/go-service-template/pkg/transport/http"
	libmware "github.com/ElfAstAhe/go-service-template/pkg/transport/http/middleware"
	_ "github.com/ElfAstAhe/tiny-audit-service/docs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	appmware "github.com/ElfAstAhe/tiny-audit-service/internal/transport/rest/middleware"
	pkgmware "github.com/ElfAstAhe/tiny-audit-service/pkg/transport/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellofresh/health-go/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	swagh "github.com/swaggo/http-swagger"
)

// AppChiRouter реализует интерфейс libhttp.Router, инкапсулируя мультиплексор chi.Mux,
// системные мидлвари, эндпоинты проверки здоровья, метрики, Swagger и доменные роуты.
type AppChiRouter struct {
	router          *chi.Mux
	log             logger.Logger
	config          *config.Config
	health          *health.Health
	livez           libhttp.LivezFunc
	readyz          libhttp.ReadyzFunc
	authAuditFacade facade.AuthAudit
	dataAuditFacade facade.DataAudit
}

var _ libhttp.Router = (*AppChiRouter)(nil)

// NewAppRouter инициализирует конфигурационные опции роутера, проводит валидацию
// зависимостей и возвращает настроенный экземпляр структуры AppChiRouter.
func NewAppRouter(opts ...Option) (*AppChiRouter, error) {
	options := &AppRouterOptions{}

	for _, opt := range opts {
		opt(options)
	}

	err := options.Validate()
	if err != nil {
		return nil, err
	}

	return newAppChiRouter(
		options.Conf,
		options.Logger,
		options.JWTHelper,
		options.JWTHTTPHelper,
		options.AuthHelper,
		options.Health,
		options.Livez,
		options.Readyz,
		options.AuthAuditFacade,
		options.DataAuditFacade,
	), err
}

func newAppChiRouter(
	config *config.Config,
	logger logger.Logger,
	jwtHelper *helper.JWTHelper,
	jwtHTTPHelper *helper.JWTHTTPHelper,
	authHelper auth.Helper,
	health *health.Health,
	livez libhttp.LivezFunc,
	readyz libhttp.ReadyzFunc,
	authAuditFacade facade.AuthAudit,
	dataAuditFacade facade.DataAudit,
) *AppChiRouter {
	res := &AppChiRouter{
		router:          chi.NewRouter(),
		log:             logger.GetLogger("app-chi-router"),
		config:          config,
		health:          health,
		livez:           livez,
		readyz:          readyz,
		authAuditFacade: authAuditFacade,
		dataAuditFacade: dataAuditFacade,
	}

	// setup middleware
	res.setupMiddleware(
		jwtHelper,
		jwtHTTPHelper,
		authHelper,
		logger,
	)

	// mount debug
	res.router.Mount("/debug", middleware.Profiler())
	// mount swagger
	res.router.Mount("/swagger/", swagh.WrapHandler)
	// mount status
	res.router.Mount("/status", res.health.Handler())
	// mount metrics
	res.router.Mount("/metrics", promhttp.Handler())

	// setup routes
	res.setupRoutes()

	return res
}

// GetRouter возвращает нижележащий сконфигурированный обработчик http.Handler (chi.Mux) для запуска HTTP-сервера.
func (cr *AppChiRouter) GetRouter() http.Handler {
	return cr.router
}

func (cr *AppChiRouter) setupMiddleware(
	jwtHelper *helper.JWTHelper,
	jwtHTTPHelper *helper.JWTHTTPHelper,
	authHelper auth.Helper,
	logger logger.Logger,
) {
	// 1. CRITICAL: Recoverer must be placed at the absolute top of the stack to intercept panics from all downstream handlers.
	cr.router.Use(middleware.Recoverer)

	// 2. LOGGING: HTTPRequestLogger records every single inbound request, including those that panic or fail auth validation.
	cr.router.Use(libmware.NewHTTPRequestLogger(logger).Handle)

	// 3. TELEMETRY & OBSERVABILITY: OpenTelemetry and Prometheus catch the final HTTP execution states and write precise metrics.
	// tracing
	cr.router.Use(otelchi.Middleware(cr.config.Telemetry.ServiceName, otelchi.WithChiRoutes(cr.router)))
	// metrics
	cr.router.Use(libmware.MetricsMiddleware)

	// 4. REQUEST IDENTIFICATION: Correlation tokens propagation for end-to-end distributed tracking.
	// requestID
	cr.router.Use(middleware.RequestID)
	// requestID (own implementation)
	cr.router.Use(libmware.NewDefaultRequestIDExtractor().Handler)
	// traceID (own implementation)
	cr.router.Use(libmware.NewDefaultTraceIDExtractor().Handler)
	// audit requestID
	cr.router.Use(pkgmware.NewDefaultAuditRequestIDExtractor().Handle)
	// audit trace id
	cr.router.Use(pkgmware.NewDefaultAuditTraceIDExtractor().Handle)

	// 5. NETWORKING CORE: Real IP extraction boundaries.
	cr.router.Use(libmware.NewDefaultRealIPExtractor().Handler)

	// 6. RESOURCE BUDGETING & PROTECTION: Network stream limits and defensive runtime constraints.
	// timeout
	cr.router.Use(middleware.Timeout(cr.config.HTTP.ReadTimeout))
	// compress (add any content-types)
	cr.router.Use(libmware.NewCompress(logger,
		"application/json", "plain/text",
	).Handle)
	// decompress
	cr.router.Use(libmware.NewDecompress(int64(cr.config.HTTP.MaxRequestBodySize), logger).Handle)

	// 7. SECURITY BOUNDARY: Authentication extraction executes right before routing to core application aggregates.
	// ROUTING OPTIMIZATION NOTICE: PathMatchers match regex rules against raw endpoint patterns securely.
	// jwt auth extractor - extract user info from token
	cr.router.Use(appmware.NewAuthExtractor(
		jwtHelper,
		jwtHTTPHelper,
		authHelper,
		logger,
		libhttp.NewHTTPPathMatchers([]*libhttp.PathMatcher{
			libhttp.NewPathMatcher(http.MethodGet, "/metrics", "^/metrics.*$"),
			libhttp.NewPathMatcher(http.MethodGet, "/swagger", "^/swagger.*$"),
			libhttp.NewPathMatcher(http.MethodGet, "/status", "^/status.*$"),
			libhttp.NewPathMatcher(http.MethodGet, "/livez", "^/livez.*$"),
			libhttp.NewPathMatcher(http.MethodGet, "/readyz", "^/readyz.*$"),
			libhttp.NewPathMatcher(http.MethodGet, "/debug", "^/debug.*$"),
			libhttp.NewPathMatcher(http.MethodGet, "/config", "^/config.*$"),
		}),
		cr.config.App.AcceptTokenIssuers,
	).Handle)
}

func (cr *AppChiRouter) setupRoutes() {
	// health check
	cr.router.Get("/livez", cr.getLivez)
	// readiness check
	cr.router.Get("/readyz", cr.getReadyz)
	// config (debug)
	if cr.config.App.Env != libconf.AppEnvProduction {
		cr.router.Get("/config", cr.getConfig)
	}

	// api
	cr.router.Route("/api", func(r chi.Router) {
		// /v1
		r.Route("/v1", func(r chi.Router) {
			// /audit
			r.Route("/audit", func(r chi.Router) {
				// /auth
				r.Route("/auth", func(r chi.Router) {
					r.Post("/", cr.postAPIV1AuditAuth)
					r.Post("/period", cr.postAPIV1AuditAuthPeriod)
					r.Post("/username", cr.postAPIV1AuditAuthUsername)
				})
				// /data
				r.Route("/data", func(r chi.Router) {
					r.Post("/", cr.postAPIV1AuditData)
					r.Post("/period", cr.postAPIV1AuditDataPeriod)
					r.Post("/instance", cr.postAPIV1AuditDataInstance)
				})
			})
		})
	})
}
