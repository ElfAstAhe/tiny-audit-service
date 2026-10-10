package grpc

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	mocks2 "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewAuditClientOptions_Defaults проверяет корректность инициализации
// структуры опций значениями по умолчанию.
func TestNewAuditClientOptions_Defaults(t *testing.T) {
	opts := NewAuditClientOptions[any]()

	require.NotNil(t, opts.BaseAuditClientOptions)
	require.NotNil(t, opts.GRPC)

	// Проверяем высокоуровневые gRPC-дефолты
	assert.Equal(t, "localhost:51052", opts.Target)
	assert.Equal(t, 2, opts.Pool.WorkerCount)
	assert.Equal(t, 10000, opts.Pool.DataCapacity)
	assert.Equal(t, 10*time.Second, opts.Pool.StopTimeout)

	// Проверяем низкоуровневые KeepAlive-дефолты
	assert.Equal(t, 15*time.Second, opts.GRPC.KATime)
	assert.True(t, opts.GRPC.KAPermitWithoutStream)
}

// TestAuditClientOptions_WithFunctions проверяет, что все Fluent API функции (опции)
// корректно перезаписывают конфигурационные поля.
func TestAuditClientOptions_WithFunctions(t *testing.T) {
	opts := NewAuditClientOptions[string]()

	mockLogger := mocks.NewMockLogger(t)
	mockTokenProvider := mocks2.NewMockTokenProvider(t)

	WithName[string]("grpc-audit-pool")(opts)
	WithTarget[string]("10.0.0.10:50051")(opts)
	WithWorkerCount[string](5)(opts)
	WithKeepAliveTime[string](30 * time.Second)(opts)
	WithSecure[string](true)(opts)
	WithLogger[string](mockLogger)(opts)
	WithTokenProvider[string](mockTokenProvider)(opts)

	assert.Equal(t, "grpc-audit-pool", opts.Pool.Name)
	assert.Equal(t, "10.0.0.10:50051", opts.Target)
	assert.Equal(t, 5, opts.Pool.WorkerCount)
	assert.Equal(t, 30*time.Second, opts.GRPC.KATime)
	assert.True(t, opts.GRPC.Secure)
	assert.Equal(t, mockLogger, opts.Pool.Logger)
	assert.Equal(t, mockTokenProvider, opts.TokenProvider)
}

// TestAuditClientOptions_Validate проверяет все сценарии валидации
// (успешные и ошибочные) с помощью табличных тестов.
func TestAuditClientOptions_Validate(t *testing.T) {
	mockLogger := mocks.NewMockLogger(t)
	mockTokenProvider := mocks2.NewMockTokenProvider(t)

	// Вспомогательная функция для создания 100% валидного объекта опций
	validOpts := func() *AuditClientOptions[string] {
		opts := NewAuditClientOptions[string]()
		WithName[string]("valid-grpc-client")(opts)
		WithLogger[string](mockLogger)(opts)
		WithTokenProvider[string](mockTokenProvider)(opts)
		return opts
	}

	tests := []struct {
		name        string
		modify      func(opts *AuditClientOptions[string])
		wantErrText string
	}{
		{
			name:        "Success valid options",
			modify:      func(opts *AuditClientOptions[string]) {},
			wantErrText: "",
		},
		{
			name: "Error empty worker pool name",
			modify: func(opts *AuditClientOptions[string]) {
				opts.Pool.Name = ""
			},
			wantErrText: "name is required",
		},
		{
			name: "Error invalid worker count",
			modify: func(opts *AuditClientOptions[string]) {
				opts.Pool.WorkerCount = 0
			},
			wantErrText: "worker count is required",
		},
		{
			name: "Error zero stop timeout",
			modify: func(opts *AuditClientOptions[string]) {
				opts.Pool.StopTimeout = 0
			},
			wantErrText: "stop timeout is required",
		},
		{
			name: "Error empty target address",
			modify: func(opts *AuditClientOptions[string]) {
				opts.Target = "   "
			},
			wantErrText: "target address is required",
		},
		{
			name: "Error nil GRPC sub-options (anti-panic check)",
			modify: func(opts *AuditClientOptions[string]) {
				opts.GRPC = nil
			},
			wantErrText: "grpc raw client options are required",
		},
		{
			name: "Error invalid low-level KeepAlive parameters",
			modify: func(opts *AuditClientOptions[string]) {
				opts.GRPC.ConnTimeout = 0 // спровоцирует ошибку внутри aco.GRPC.Validate()
			},
			wantErrText: "grpc client options validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := validOpts()
			tt.modify(opts)

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
