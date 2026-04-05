package rest

import (
	"net/http"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport"
	libmware "github.com/ElfAstAhe/go-service-template/pkg/transport/middleware"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	trmware "github.com/ElfAstAhe/tiny-audit-service/internal/transport/rest/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellofresh/health-go/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	swagh "github.com/swaggo/http-swagger"
)

type AppChiRouter struct {
	router          *chi.Mux
	log             logger.Logger
	config          *config.Config
	health          *health.Health
	healthz         transport.HealthzFunc
	readyz          transport.ReadyzFunc
	authAuditFacade facade.AuthAuditFacade
	dataAuditFacade facade.DataAuditFacade
}

var _ transport.HTTPRouter = (*AppChiRouter)(nil)

func NewAppChiRouter(
	config *config.Config,
	logger logger.Logger,
	jwtHelper *helper.JWTHelper,
	jwtHTTPHelper *helper.JWTHTTPHelper,
	authHelper auth.Helper,
	health *health.Health,
	healthz transport.HealthzFunc,
	readyz transport.ReadyzFunc,
	authAuditFacade facade.AuthAuditFacade,
	dataAuditFacade facade.DataAuditFacade,
) *AppChiRouter {
	res := &AppChiRouter{
		router:          chi.NewRouter(),
		log:             logger.GetLogger("app-chi-router"),
		config:          config,
		health:          health,
		healthz:         healthz,
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

func (cr *AppChiRouter) GetRouter() http.Handler {
	return cr.router
}

func (cr *AppChiRouter) setupMiddleware(
	jwtHelper *helper.JWTHelper,
	jwtHTTPHelper *helper.JWTHTTPHelper,
	authHelper auth.Helper,
	logger logger.Logger,
) {
	// tracing
	cr.router.Use(otelchi.Middleware(cr.config.Telemetry.ServiceName, otelchi.WithChiRoutes(cr.router)))
	// metrics
	cr.router.Use(libmware.HTTPMetricsMiddleware)
	// requestID
	cr.router.Use(middleware.RequestID)
	// realIP
	cr.router.Use(middleware.RealIP)
	// recoverer
	cr.router.Use(middleware.Recoverer)
	// timeout
	cr.router.Use(middleware.Timeout(cr.config.HTTP.ReadTimeout))
	// compress (add any content-types)
	cr.router.Use(libmware.NewHTTPCompress(logger,
		"application/json", "plain/text",
	).Handle)
	// decompress
	cr.router.Use(libmware.NewHTTPDecompress(int64(cr.config.HTTP.MaxRequestBodySize), logger).Handle)
	// jwt auth extractor - extract user info from token
	cr.router.Use(trmware.NewAuthExtractor(
		jwtHelper,
		jwtHTTPHelper,
		authHelper,
		logger,
		transport.NewHTTPPathMatchers([]*transport.HTTPPathMatcher{
			transport.NewHTTPPathMatcher(http.MethodGet, "/metrics", "^/metrics*$"),
			transport.NewHTTPPathMatcher(http.MethodGet, "/swagger", "^/swagger*$"),
			transport.NewHTTPPathMatcher(http.MethodGet, "/status", "^/status*$"),
			transport.NewHTTPPathMatcher(http.MethodGet, "/healthz", "^/healthz*$"),
			transport.NewHTTPPathMatcher(http.MethodGet, "/readyz", "^/readyz*$"),
			transport.NewHTTPPathMatcher(http.MethodGet, "/debug", "^/debug*$"),
		}),
		cr.config.App.AcceptTokenIssuers,
	).Handle)
	// income/outcome logger
	cr.router.Use(libmware.NewHTTPRequestLogger(logger).Handle)
}

func (cr *AppChiRouter) setupRoutes() {
	// health check
	cr.router.Get("/healthz", cr.getHealthz)
	// readiness check
	cr.router.Get("/readyz", cr.getReadyz)
	// config (debug)
	if cr.config.App.Env != config.AppEnvProduction {
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
			/*
			   // crud, может и не нужен...
			   // /auth
			   r.Route("/auth", func(r chi.Router) {
			       r.Get("/{id}", cr.getAPIV1AuthAudit)
			       r.Get("/", cr.getAPIV1AuthAudits)
			       r.Post("/", cr.postAPIV1Auth)
			       r.Put("/{id}", cr.putAPIV1Auth)
			       r.Delete("/{id}", cr.deleteAPIV1Auth)
			   })
			   // /data
			   r.Route("/data", func(r chi.Router) {
			       r.Get("/{id}", cr.getAPIV1DataAudit)
			       r.Get("/", cr.getAPIV1DataAudits)
			       r.Post("/", cr.postAPIV1DataAudit)
			       r.Put("/{id}", cr.putAPIV1DataAudit)
			       r.Delete("/{id}", cr.deleteAPIV1DataAudit)
			   })
			*/
		})
	})
}
