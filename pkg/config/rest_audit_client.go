package config

import (
	"net/url"
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// RestAuditClientConfig содержит параметры конфигурации REST-клиента аудита.
// Включает настройки пула воркеров через встраивание структуры WorkerPoolConfig.
type RestAuditClientConfig struct {
	*config.WorkerPoolConfig `mapstructure:",squash"`
	BaseURL                  string        `mapstructure:"base_url" json:"base_url,omitempty" yaml:"base_url,omitempty"`
	ReadTimeout              time.Duration `mapstructure:"read_timeout" json:"read_timeout,omitempty" yaml:"read_timeout,omitempty"`
}

// NewRestAuditClientConfig конструирует новый экземпляр RestAuditClientConfig с переданными параметрами.
func NewRestAuditClientConfig(
	workerPool *config.WorkerPoolConfig,
	baseURL string,
	readTimeout time.Duration,
) *RestAuditClientConfig {
	return &RestAuditClientConfig{
		WorkerPoolConfig: workerPool,
		BaseURL:          baseURL,
		ReadTimeout:      readTimeout,
	}
}

// NewDefaultRestAuditClientConfig создает конфигурацию REST-клиента аудита
// с дефолтными параметрами пула воркеров, сетевого адреса и таймаутов.
func NewDefaultRestAuditClientConfig() *RestAuditClientConfig {
	return NewRestAuditClientConfig(
		config.NewDefaultWorkerPoolConfig(),
		DefaultRestBaseURL,
		DefaultRestReadTimeout,
	)
}

// Validate выполняет глубокую проверку всех параметров конфигурации REST-клиента.
// Защищен от nil-указателей встроенного пула и проверяет корректность URL-схемы.
func (rac *RestAuditClientConfig) Validate() error {
	// Защита от Nil Pointer паники
	if rac.WorkerPoolConfig == nil {
		return errs.NewConfigValidateError("rest audit client", "pool", "worker pool configuration is missing", nil)
	}
	if err := rac.WorkerPoolConfig.Validate(); err != nil {
		return errs.NewConfigValidateError("rest audit client", "pool", "validate failed", err)
	}
	if strings.TrimSpace(rac.BaseURL) == "" {
		return errs.NewConfigValidateError("rest audit client", "base url", "base url is required", nil)
	}
	u, err := url.Parse(rac.BaseURL)
	if err != nil {
		return errs.NewConfigValidateError("rest audit client", "base url", "validate failed", err)
	}
	// Строгая валидация сетевого адреса REST API
	if u.Scheme != "http" && u.Scheme != "https" {
		return errs.NewConfigValidateError("rest audit client", "base url", "url scheme must be http or https", nil)
	}
	if u.Host == "" {
		return errs.NewConfigValidateError("rest audit client", "base url", "url host is empty", nil)
	}
	if rac.ReadTimeout <= 0 {
		return errs.NewConfigValidateError("rest audit client", "read timeout", "validate failed", err)
	}

	return nil
}
