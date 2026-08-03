package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	mocks2 "github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp/azure"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp/mocks"
	libworker "github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
	mocks3 "github.com/ElfAstAhe/tiny-audit-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Вспомогательный метод для быстрой инициализации мока логгера без шума в тестах
func setupMockLogger(t *testing.T) *mocks2.MockLogger {
	loggerMock := mocks2.NewMockLogger(t)
	loggerMock.On("GetLogger", mock.Anything).Return(loggerMock)
	// Разрешаем любые вызовы Debugf и Errorf с любыми аргументами, чтобы тесты не падали на логах
	loggerMock.EXPECT().Debugf(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	loggerMock.EXPECT().Debugf(mock.Anything, mock.Anything, mock.Anything).Maybe()
	loggerMock.EXPECT().Debugf(mock.Anything, mock.Anything).Maybe()
	loggerMock.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()
	return loggerMock
}

// ============================================================================
// 1. ТЕСТЫ КОНСТРУКТОРА И ВАЛИДАЦИИ ОПЦИЙ
// ============================================================================

func TestNewLoginAttempts_ValidationFailed(t *testing.T) {
	_, err := NewLoginAttempts(WithLAOName(""))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestNewLoginAttempts_Success(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := NewLoginAttempts(
		WithLAOName("test-login-attempts"),
		WithLAOLogger(mockLog),
		WithLAODispatcherOpts(&libworker.BaseSchedulerDispatcherConfig{}),
		WithLAOReceiver(mockReceiver),
		WithLAOAuthAuditUseCase(mockUC),
		WithLAOBatchSize(10),
		WithLAOBatchReadTimeout(1*time.Second),
		WithLAOAcknowledgeTimeout(1*time.Second),
	)

	require.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, la.batchSize)
}

// ============================================================================
// 2. ТЕСТЫ ДЛЯ DATA PROVIDER
// ============================================================================

func TestLoginAttempts_DataProvider_SuccessBatchSize(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		batchSize:        2,
		batchReadTimeout: 1 * time.Second,
		receiver:         mockReceiver,
		opts:             NewLoginAttemptsOptions(),
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	fakeMsg := azure.NewMessage([]byte(`{"id":"1","value":"test"}`), nil)

	// Настраиваем Mockery: Receive должен вернуть валидные данные 2 раза
	mockReceiver.On("Receive", mock.Anything, mock.Anything).Return(fakeMsg, nil).Times(2)

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestLoginAttempts_DataProvider_ReadTimeout(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		batchSize:        5,
		batchReadTimeout: 50 * time.Millisecond,
		receiver:         mockReceiver,
		opts:             NewLoginAttemptsOptions(),
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	fakeMsg := azure.NewMessage([]byte(`{"id":"1"}`), nil)

	mockReceiver.On("Receive", mock.Anything, mock.Anything).Return(fakeMsg, nil).Once()
	// На второй итерации имитируем, что брокер пуст и время вышло
	mockReceiver.On("Receive", mock.Anything, mock.Anything).Return(nil, context.DeadlineExceeded).Once()

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err)
	assert.Len(t, res, 1, "Должны вернуть то, что накопили до наступления таймаута")
}

func TestLoginAttempts_DataProvider_ContextCanceled(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		batchSize:        5,
		batchReadTimeout: 1 * time.Second,
		receiver:         mockReceiver,
		opts:             NewLoginAttemptsOptions(),
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	// Имитируем Graceful Shutdown (SIGTERM во время ожидания сообщения брокера)
	mockReceiver.On("Receive", mock.Anything, mock.Anything).Return(nil, context.Canceled).Once()

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err, "При отмене контекста ошибку наружу не выкидываем")
	assert.Empty(t, res, "Пачка должна сброситься и вернуться пустой")
}

func TestLoginAttempts_DataProvider_MapperError_RejectSuccess(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		batchSize:        2,
		batchReadTimeout: 1 * time.Second,
		receiver:         mockReceiver,
		opts:             NewLoginAttemptsOptions(),
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	badMsg := azure.NewMessage([]byte(`{ broken json }`), nil)

	mockReceiver.On("Receive", mock.Anything, mock.Anything).Return(badMsg, nil).Once()
	// Проверяем, что сработал Reject сообщения в DLQ
	mockReceiver.On("Reject", mock.Anything, badMsg, mock.Anything).Return(nil).Once()
	// Завершаем цикл по таймауту на второй круг
	mockReceiver.On("Receive", mock.Anything, mock.Anything).Return(nil, context.DeadlineExceeded).Once()

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err)
	assert.Empty(t, res, "Сломанные данные не должны попасть в пачку воркеров")
}

// ============================================================================
// 3. ТЕСТЫ ДЛЯ STORE AUTH AUDIT
// ============================================================================

func TestLoginAttempts_StoreAuthAudit_Success(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		receiver:           mockReceiver,
		authAuditUC:        mockUC,
		acknowledgeTimeout: 1 * time.Second,
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	testDTO := &dto.AuthAuditDTO{
		Message: azure.NewMessage([]byte(`{}`), nil),
	}

	// Успешный аудит в БД тянет за собой успешное подтверждение (Accept) в Azure
	mockUC.On("Audit", mock.Anything, mock.Anything).Return(nil).Once()
	mockReceiver.On("Accept", mock.Anything, testDTO.Message).Return(nil).Once()

	err := la.storeAuthAudit(context.Background(), 1, testDTO)

	assert.NoError(t, err)
}

func TestLoginAttempts_StoreAuthAudit_UniqueViolation_Accept(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		receiver:           mockReceiver,
		authAuditUC:        mockUC,
		acknowledgeTimeout: 1 * time.Second,
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	testDTO := &dto.AuthAuditDTO{
		Message: azure.NewMessage([]byte(`{}`), nil),
	}

	// База сообщает о дубликате записи (UniqueViolation)
	uniqueErr := &errs.BllUniqueError{}
	mockUC.On("Audit", mock.Anything, mock.Anything).Return(uniqueErr).Once()

	// Наш код на Go 1.26 должен проглотить ошибку дубликата и всё равно выполнить Accept
	mockReceiver.On("Accept", mock.Anything, testDTO.Message).Return(nil).Once()

	err := la.storeAuthAudit(context.Background(), 1, testDTO)

	assert.NoError(t, err, "Ошибка уникальности должна гаситься на месте")
}

func TestLoginAttempts_StoreAuthAudit_DatabaseError_Release(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver[*amqp.ReceiveOptions](t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la := &LoginAttempts{
		receiver:           mockReceiver,
		authAuditUC:        mockUC,
		acknowledgeTimeout: 1 * time.Second,
	}
	la.BaseSchedulerDispatcher = libworker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		"test", &libworker.BaseSchedulerDispatcherConfig{}, nil, nil, mockLog,
	)

	testDTO := &dto.AuthAuditDTO{
		Message: azure.NewMessage([]byte(`{}`), nil),
	}

	// Имитируем падение коннекта к СУБД Postgres
	criticalDbErr := errors.New("postgres connection timeout")
	mockUC.On("Audit", mock.Anything, mock.Anything).Return(criticalDbErr).Once()

	// Метод должен сделать Release, чтобы брокер вернул сообщение обратно в очередь
	mockReceiver.On("Release", mock.Anything, testDTO.Message).Return(nil).Once()

	err := la.storeAuthAudit(context.Background(), 1, testDTO)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "store auth audit failed")
}
