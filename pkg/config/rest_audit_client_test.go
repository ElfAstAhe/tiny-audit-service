package config

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultRestAuditClientConfig(t *testing.T) {
	cfg := NewDefaultRestAuditClientConfig()

	require.NotNil(t, cfg.WorkerPoolConfig)
	assert.Equal(t, DefaultRestBaseURL, cfg.BaseURL)
	assert.Equal(t, DefaultRestReadTimeout, cfg.ReadTimeout)
}

func TestRestAuditClientConfig_Validate(t *testing.T) {
	// Вспомогательный генератор 100% валидного конфига шаблона
	// Предполагаем, что для валидности WorkerPoolConfig имя пула не должно быть пустым
	validPool := func() *config.WorkerPoolConfig {
		p := config.NewDefaultWorkerPoolConfig()
		// Если у вашего шаблона есть обязательные поля типа Name, задайте их тут:
		// p.Name = "test-pool"
		return p
	}

	tests := []struct {
		name        string
		modify      func(cfg *RestAuditClientConfig)
		wantErrText string
	}{
		{
			name:        "Success valid configuration",
			modify:      func(cfg *RestAuditClientConfig) {},
			wantErrText: "",
		},
		{
			name: "Error nil worker pool config (anti-panic check)",
			modify: func(cfg *RestAuditClientConfig) {
				cfg.WorkerPoolConfig = nil
			},
			wantErrText: "configuration is missing",
		},
		{
			name: "Error empty base url",
			modify: func(cfg *RestAuditClientConfig) {
				cfg.BaseURL = "   "
			},
			wantErrText: "base url is required",
		},
		{
			name: "Error invalid scheme in base url",
			modify: func(cfg *RestAuditClientConfig) {
				cfg.BaseURL = "ftp://localhost:8080"
			},
			wantErrText: "url scheme must be http or https",
		},
		{
			name: "Error negative read timeout",
			modify: func(cfg *RestAuditClientConfig) {
				cfg.ReadTimeout = -5 * time.Second
			},
			wantErrText: "validate failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewRestAuditClientConfig(
				validPool(),
				DefaultRestBaseURL,
				DefaultRestReadTimeout,
			)
			tt.modify(cfg)

			err := cfg.Validate()

			if tt.wantErrText != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
