package client

import (
	"context"
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	mocks2 "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестируем, что дефолтные опции создаются с правильными значениями
func TestNewBaseAuditClientOptions_Defaults(t *testing.T) {
	opts := NewBaseAuditClientOptions[any]()

	assert.Equal(t, DefaultWorkerCount, opts.WorkerCount)
	assert.Equal(t, DefaultDataCapacity, opts.DataCapacity)
	assert.Equal(t, DefaultCompleteProcess, opts.CompleteProcess)
	assert.Equal(t, worker.DefaultPoolStopTimeout, opts.StopTimeout)
	assert.Empty(t, opts.Name)
	assert.Nil(t, opts.Logger)
}

// Тестируем Fluent API (функции With...)
func TestWithFunctions(t *testing.T) {
	opts := NewBaseAuditClientOptions[string]()

	// Инициализируем моки
	mockLogger := mocks.NewMockLogger(t)
	mockTokenProvider := mocks2.NewMockTokenProvider(t)
	var dummyAction AuditAction[string] = func(ctx context.Context, workerIndex int, data string, token string) error { return nil }

	// Применяем опции
	WithName[string]("test-audit")(opts)
	WithPoolWorkerCount[string](10)(opts)
	WithPoolDataCapacity[string](128)(opts)
	WithPoolCompleteProcess[string](false)(opts)
	WithPoolStopTimeout[string](time.Second * 10)(opts)
	WithLogger[string](mockLogger)(opts)
	WithTokenProvider[string](mockTokenProvider)(opts)
	WithAuditAction[string](dummyAction)(opts)

	// Проверяем, что все поля перезаписались
	assert.Equal(t, "test-audit", opts.Name)
	assert.Equal(t, 10, opts.WorkerCount)
	assert.Equal(t, 128, opts.DataCapacity)
	assert.False(t, opts.CompleteProcess)
	assert.Equal(t, time.Second*10, opts.StopTimeout)
	assert.Equal(t, mockLogger, opts.Logger)
	assert.Equal(t, mockTokenProvider, opts.TokenProvider)
	assert.NotNil(t, opts.AuditAction)
}

// Тестируем валидацию (Table-Driven Tests)
func TestBaseAuditClientOptions_Validate(t *testing.T) {
	mockLogger := mocks.NewMockLogger(t)
	mockTokenProvider := mocks2.NewMockTokenProvider(t)
	var dummyAction AuditAction[string] = func(ctx context.Context, workerIndex int, data string, token string) error { return nil }

	// Вспомогательная функция для создания валидного набора опций
	validOpts := func() *BaseAuditClientOptions[string] {
		opts := NewBaseAuditClientOptions[string]()
		opts.Name = "valid-name"
		opts.Logger = mockLogger
		opts.TokenProvider = mockTokenProvider
		opts.AuditAction = dummyAction
		return opts
	}

	// Структура для тест-кейсов
	tests := []struct {
		name        string
		modify      func(opts *BaseAuditClientOptions[string])
		wantErrText string
	}{
		{
			name:        "Success valid options",
			modify:      func(opts *BaseAuditClientOptions[string]) {}, // ничего не меняем
			wantErrText: "",
		},
		{
			name: "Error empty name",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.Name = "   "
			},
			wantErrText: "name is required",
		},
		{
			name: "Error worker count zero",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.WorkerCount = 0
			},
			wantErrText: "worker count is required",
		},
		{
			name: "Error data capacity negative",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.DataCapacity = -1
			},
			wantErrText: "data capacity is required",
		},
		{
			name: "Error stop timeout zero",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.StopTimeout = 0
			},
			wantErrText: "stop timeout is required",
		},
		{
			name: "Error nil logger",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.Logger = nil // Внимание: если utils.IsNil проверяет интерфейс, это сработает
			},
			wantErrText: "logger is required",
		},
		{
			name: "Error nil token provider",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.TokenProvider = nil
			},
			wantErrText: "token provider is required",
		},
		{
			name: "Error nil audit action",
			modify: func(opts *BaseAuditClientOptions[string]) {
				opts.AuditAction = nil
			},
			wantErrText: "audit action is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := validOpts()
			tt.modify(opts) // вносим ломающие изменения для конкретного кейса

			err := opts.Validate()

			if tt.wantErrText != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
