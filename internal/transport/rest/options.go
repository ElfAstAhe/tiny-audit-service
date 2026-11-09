package rest

import (
	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libhttp "github.com/ElfAstAhe/go-service-template/pkg/transport/http"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/hellofresh/health-go/v5"
)

// AppRouterOptions группирует все зависимости, конфигурации и компоненты фасадов,
// необходимые для успешной сборки и инициализации HTTP-роутера приложения.
type AppRouterOptions struct {
	Conf            *config.Config
	Logger          logger.Logger
	JWTHelper       *helper.JWTHelper
	JWTHTTPHelper   *helper.JWTHTTPHelper
	AuthHelper      auth.Helper
	Health          *health.Health
	Livez           libhttp.LivezFunc
	Readyz          libhttp.ReadyzFunc
	AuthAuditFacade facade.AuthAudit
	DataAuditFacade facade.DataAudit
}

// Validate проверяет полноту и корректность заполнения всех обязательных зависимостей в структуре AppRouterOptions.
func (aro *AppRouterOptions) Validate() error {
	if utils.IsNil(aro.Conf) {
		return errs.NewTlCommonError("validate", "conf not applied", nil)
	}
	if utils.IsNil(aro.Logger) {
		return errs.NewTlCommonError("validate", "logger not applied", nil)
	}
	if utils.IsNil(aro.JWTHelper) {
		return errs.NewTlCommonError("validate", "jwt helper not applied", nil)
	}
	if utils.IsNil(aro.JWTHTTPHelper) {
		return errs.NewTlCommonError("validate", "jwt http helper not applied", nil)
	}
	if utils.IsNil(aro.AuthHelper) {
		return errs.NewTlCommonError("validate", "auth helper not applied", nil)
	}
	if utils.IsNil(aro.Health) {
		return errs.NewTlCommonError("validate", "health not applied", nil)
	}
	//if utils.IsNil(aro.Healthz) {
	//    return errs.NewTlCommonError("validate", "healthz not applied", nil)
	//}
	//if utils.IsNil(aro.Readyz) {
	//    return errs.NewTlCommonError("validate", "readyz not applied", nil)
	//}
	if utils.IsNil(aro.AuthAuditFacade) {
		return errs.NewTlCommonError("validate", "auth audit facade not applied", nil)
	}
	if utils.IsNil(aro.DataAuditFacade) {
		return errs.NewTlCommonError("validate", "data audit facade not applied", nil)
	}

	return nil
}

// Option определяет сигнатуру функции для конфигурирования структуры AppRouterOptions.
type Option func(options *AppRouterOptions)

// WithConfig возвращает опцию для настройки конфигурации приложения.
func WithConfig(conf *config.Config) Option {
	return func(o *AppRouterOptions) {
		o.Conf = conf
	}
}

// WithLogger возвращает опцию для настройки системного логгера.
func WithLogger(logger logger.Logger) Option {
	return func(o *AppRouterOptions) {
		o.Logger = logger
	}
}

// WithJwtHelper возвращает опцию для настройки утилиты JWTHelper.
func WithJwtHelper(helper *helper.JWTHelper) Option {
	return func(o *AppRouterOptions) {
		o.JWTHelper = helper
	}
}

// WithJwtHTTPHelper возвращает опцию для настройки HTTP-хелпера работы с JWT токенами в запросах.
func WithJwtHTTPHelper(helper *helper.JWTHTTPHelper) Option {
	return func(o *AppRouterOptions) {
		o.JWTHTTPHelper = helper
	}
}

// WithAuthHelper возвращает опцию для настройки хелпера контроля авторизации субъектов.
func WithAuthHelper(helper auth.Helper) Option {
	return func(o *AppRouterOptions) {
		o.AuthHelper = helper
	}
}

// WithHealth возвращает опцию для интеграции инстанса проверки здоровья (Health Check) сервиса.
func WithHealth(healthInst *health.Health) Option {
	return func(o *AppRouterOptions) {
		o.Health = healthInst
	}
}

// WithLivez возвращает опцию для настройки функции проверки Liveness probe (жизнеспособности) приложения.
func WithLivez(livez libhttp.LivezFunc) Option {
	return func(o *AppRouterOptions) {
		o.Livez = livez
	}
}

// WithReadyz возвращает опцию для настройки функции проверки Readiness probe (готовности к приему трафика).
func WithReadyz(readyz libhttp.ReadyzFunc) Option {
	return func(o *AppRouterOptions) {
		o.Readyz = readyz
	}
}

// WithAuthAuditFacade возвращает опцию для привязки фасада управления аудитом аутентификации.
func WithAuthAuditFacade(authAuditFacade facade.AuthAudit) Option {
	return func(o *AppRouterOptions) {
		o.AuthAuditFacade = authAuditFacade
	}
}

// WithDataAuditFacade возвращает опцию для привязки фасада управления аудитом изменений данных.
func WithDataAuditFacade(dataAuditFacade facade.DataAudit) Option {
	return func(o *AppRouterOptions) {
		o.DataAuditFacade = dataAuditFacade
	}
}
