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

type AppRouterOptions struct {
	Conf            *config.Config
	Logger          logger.Logger
	JWTHelper       *helper.JWTHelper
	JWTHTTPHelper   *helper.JWTHTTPHelper
	AuthHelper      auth.Helper
	Health          *health.Health
	Healthz         libhttp.HealthzFunc
	Readyz          libhttp.ReadyzFunc
	AuthAuditFacade facade.AuthAuditFacade
	DataAuditFacade facade.DataAuditFacade
}

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

type Option func(options *AppRouterOptions)

func WithConfig(conf *config.Config) Option {
	return func(o *AppRouterOptions) {
		o.Conf = conf
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(o *AppRouterOptions) {
		o.Logger = logger
	}
}

func WithJwtHelper(helper *helper.JWTHelper) Option {
	return func(o *AppRouterOptions) {
		o.JWTHelper = helper
	}
}

func WithJwtHTTPHelper(helper *helper.JWTHTTPHelper) Option {
	return func(o *AppRouterOptions) {
		o.JWTHTTPHelper = helper
	}
}

func WithAuthHelper(helper auth.Helper) Option {
	return func(o *AppRouterOptions) {
		o.AuthHelper = helper
	}
}

func WithHealth(healthInst *health.Health) Option {
	return func(o *AppRouterOptions) {
		o.Health = healthInst
	}
}

func WithHealthz(healthz libhttp.HealthzFunc) Option {
	return func(o *AppRouterOptions) {
		o.Healthz = healthz
	}
}

func WithReadyz(readyz libhttp.ReadyzFunc) Option {
	return func(o *AppRouterOptions) {
		o.Readyz = readyz
	}
}

func WithAuthAuditFacade(authAuditFacade facade.AuthAuditFacade) Option {
	return func(o *AppRouterOptions) {
		o.AuthAuditFacade = authAuditFacade
	}
}

func WithDataAuditFacade(dataAuditFacade facade.DataAuditFacade) Option {
	return func(o *AppRouterOptions) {
		o.DataAuditFacade = dataAuditFacade
	}
}
