package app

import (
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

type Option func(*Application)

func WithConfig(conf *config.Config) Option {
	return func(app *Application) {
		app.conf = conf
	}
}

func WithLogger(log logger.Logger) Option {
	return func(app *Application) {
		app.log = log
	}
}
