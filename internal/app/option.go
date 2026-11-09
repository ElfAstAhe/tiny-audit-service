package app

import (
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

// Option определяет сигнатуру функции для конфигурирования структуры Application.
type Option func(*Application)

// WithConfig возвращает опцию для настройки корневой конфигурации приложения.
func WithConfig(conf *config.Config) Option {
	return func(app *Application) {
		app.conf = conf
	}
}

// WithLogger возвращает опцию для настройки системного логгера приложения.
func WithLogger(log logger.Logger) Option {
	return func(app *Application) {
		app.log = log
	}
}
