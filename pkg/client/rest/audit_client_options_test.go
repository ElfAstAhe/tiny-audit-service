package rest

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	mocks2 "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditClientOptions_Defaults(t *testing.T) {
	opts := NewAuditClientOptions[any]()

	require.NotNil(t, opts.BaseAuditClientOptions)
	require.NotNil(t, opts.Pool)

	// Проверяем, что локальные дефолты перетерли дефолты базового клиента внутри Pool
	assert.Equal(t, 2, opts.Pool.WorkerCount)
	assert.Equal(t, 10000, opts.Pool.DataCapacity)
	assert.True(t, opts.Pool.CompleteProcess)
	assert.Equal(t, 10*time.Second, opts.Pool.StopTimeout)

	// Проверяем REST-специфичные опции
	assert.Equal(t, "http://localhost:8080/", opts.BaseURL)
	assert.Equal(t, 5*time.Second, opts.ReadTimeout)
}

func TestAuditClientOptions_WithFunctions(t *testing.T) {
	opts := NewAuditClientOptions[string]()

	mockLogger := mocks.NewMockLogger(t)
	mockTokenProvider := mocks2.NewMockTokenProvider(t)

	WithName[string]("rest-pool-client")(opts)
	WithWorkerCount[string](15)(opts)
	WithBaseURL[string]("https://service.internal")(opts)
	WithLogger[string](mockLogger)(opts)
	WithTokenProvider[string](mockTokenProvider)(opts)

	assert.Equal(t, "rest-pool-client", opts.Pool.Name)
	assert.Equal(t, 15, opts.Pool.WorkerCount)
	assert.Equal(t, "https://service.internal", opts.BaseURL)
	assert.Equal(t, mockLogger, opts.Pool.Logger)
	assert.Equal(t, mockTokenProvider, opts.TokenProvider)
}
