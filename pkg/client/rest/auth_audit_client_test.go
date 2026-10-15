package rest

import (
	"testing"

	"github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	mocks2 "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthAuditClient_GetName(t *testing.T) {
	mockLogger := mocks.NewMockLogger(t)
	// Настраиваем поведение суб-логгера, так как в конструкторе вызывается GetLogger(...)
	mockLogger.On("GetLogger", mock.Anything).Return(mockLogger)

	mockTokenProvider := mocks2.NewMockTokenProvider(t)

	opts := []AuditClientOption[*dto.AuthAuditDTO]{
		WithName[*dto.AuthAuditDTO]("test-pool"),
		WithBaseURL[*dto.AuthAuditDTO]("http://localhost:8080"),
		WithLogger[*dto.AuthAuditDTO](mockLogger),
		WithTokenProvider[*dto.AuthAuditDTO](mockTokenProvider),
	}

	aac, err := NewAuthAuditClient(opts...)
	require.NoError(t, err)
	require.NotNil(t, aac)

	assert.Equal(t, "rest-auth-audit-client-test-pool", aac.GetName())
}
