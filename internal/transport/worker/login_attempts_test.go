package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	mocks2 "github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker/amqp/azure"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker/mocks"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
	mocks3 "github.com/ElfAstAhe/tiny-audit-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupMockLogger настраивает мок логгера для подавления шума в выводах тестов.
func setupMockLogger(t *testing.T) *mocks2.MockLogger {
	loggerMock := mocks2.NewMockLogger(t)
	loggerMock.On("GetLogger", mock.Anything).Return(loggerMock).Maybe()
	loggerMock.EXPECT().Debugf(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	loggerMock.EXPECT().Debugf(mock.Anything, mock.Anything, mock.Anything).Maybe()
	loggerMock.EXPECT().Debugf(mock.Anything, mock.Anything).Maybe()
	loggerMock.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()
	return loggerMock
}

// createTestWorker собирает воркер через официальный конструктор с правильной передачей опций диспетчера.
func createTestWorker(
	mockReceiver *mocks.MockReceiver,
	mockUC *mocks3.MockAuthAuditUseCase,
	mockLog *mocks2.MockLogger,
	batchSize int,
) (*LoginAttempts, error) {
	return NewLoginAttempts(
		WithLAOName("test-login-attempts"),
		WithLAOLogger(mockLog),
		// Передаем опции диспетчера в полном соответствии с сигнатурой WithLAODispatcherOpts
		WithLAODispatcherOpts(
			worker.WithSchedulerDispatcherPoolWorkerCount[*dto.LoginAttemptWorkerJob](1),
			worker.WithSchedulerDispatcherPoolDataCapacity[*dto.LoginAttemptWorkerJob](10),
			worker.WithSchedulerDispatcherPoolCompleteProcess[*dto.LoginAttemptWorkerJob](true),
			worker.WithSchedulerDispatcherSchedulerStartInterval[*dto.LoginAttemptWorkerJob](time.Second),
			worker.WithSchedulerDispatcherSchedulerScheduleInterval[*dto.LoginAttemptWorkerJob](time.Minute),
			worker.WithSchedulerDispatcherStopTimeout[*dto.LoginAttemptWorkerJob](time.Second),
		),
		WithLAOReceiver(mockReceiver),
		WithLAOAuthAuditUseCase(mockUC),
		WithLAOBatchSize(batchSize),
		WithLAOBatchReadTimeout(50*time.Millisecond),
		WithLAOAcknowledgeTimeout(1*time.Second),
	)
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
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 10)

	require.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, la.batchSize)
}

// ============================================================================
// 2. ТЕСТЫ ДЛЯ DATA PROVIDER
// ============================================================================

func TestLoginAttempts_DataProvider_SuccessBatchSize(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 2)
	require.NoError(t, err)

	// Передаем валидный JSON-объект, чтобы внутренний mapper успешно распарсил структуру
	fakeMsg := azure.NewMessage([]byte(`{"username":"john_doe","event":"login"}`), nil)

	// Исправлено: Метод Receive принимает ровно 1 аргумент context.Context
	mockReceiver.On("Receive", mock.Anything).Return(fakeMsg, nil).Times(2)

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestLoginAttempts_DataProvider_ReadTimeout(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 5)
	require.NoError(t, err)

	fakeMsg := azure.NewMessage([]byte(`{"username":"john_doe"}`), nil)

	// Исправлено: Синхронизация по числу аргументов Receive
	mockReceiver.On("Receive", mock.Anything).Return(fakeMsg, nil).Once()
	mockReceiver.On("Receive", mock.Anything).Return(nil, context.DeadlineExceeded).Once()

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err)
	assert.Len(t, res, 1, "Должны вернуть то, что успели вычитать из брокера до таймаута")
}

func TestLoginAttempts_DataProvider_ContextCanceled(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 5)
	require.NoError(t, err)

	mockReceiver.On("Receive", mock.Anything).Return(nil, context.Canceled).Once()

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err, "При отмене контекста ошибка гасится внутри функции")
	assert.Empty(t, res, "Пачка должна сброситься и вернуться пустой")
}

func TestLoginAttempts_DataProvider_MapperError_RejectSuccess(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 2)
	require.NoError(t, err)

	badMsg := azure.NewMessage([]byte(`{ битый json }`), nil)

	// Исправлено: Соответствие аргументов для Receive и Reject
	mockReceiver.On("Receive", mock.Anything).Return(badMsg, nil).Once()
	mockReceiver.On("Reject", mock.Anything, badMsg, mock.Anything).Return(nil).Once()
	mockReceiver.On("Receive", mock.Anything).Return(nil, context.DeadlineExceeded).Once()

	res, err := la.dataProvider(context.Background(), time.Now())

	require.NoError(t, err)
	assert.Empty(t, res, "Сообщения с ошибкой маппинга отсекаются и уходят в DLQ")
}

// ============================================================================
// 3. ТЕСТЫ ДЛЯ STORE AUTH AUDIT
// ============================================================================

func TestLoginAttempts_StoreAuthAudit_Success(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 2)
	require.NoError(t, err)

	// Инициализируем внутренний DTO, чтобы маппер внутри storeAuthAudit не падал с nil pointer
	testDTO := &dto.LoginAttemptWorkerJob{
		Message: azure.NewMessage([]byte(`{}`), nil),
		Data:    &dto.LoginAttemptEventDTO{},
	}

	mockUC.On("Audit", mock.Anything, mock.Anything).Return(nil).Once()
	// Исправлено: Accept принимает строго 2 аргумента (ctx, message)
	mockReceiver.On("Accept", mock.Anything, testDTO.Message).Return(nil).Once()

	err = la.storeAuthAudit(context.Background(), 1, testDTO)

	assert.NoError(t, err)
}

func TestLoginAttempts_StoreAuthAudit_UniqueViolation_Accept(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 2)
	require.NoError(t, err)

	testDTO := &dto.LoginAttemptWorkerJob{
		Message: azure.NewMessage([]byte(`{}`), nil),
		Data:    &dto.LoginAttemptEventDTO{},
	}

	uniqueErr := &errs.BllUniqueError{}
	mockUC.On("Audit", mock.Anything, mock.Anything).Return(uniqueErr).Once()

	// Исправлено: Accept принимает 2 аргумента
	mockReceiver.On("Accept", mock.Anything, testDTO.Message).Return(nil).Once()

	err = la.storeAuthAudit(context.Background(), 1, testDTO)

	assert.NoError(t, err, "Ошибка дубликата должна успешно обрабатываться без прерывания")
}

func TestLoginAttempts_StoreAuthAudit_DatabaseError_Release(t *testing.T) {
	mockReceiver := mocks.NewMockReceiver(t)
	mockUC := mocks3.NewMockAuthAuditUseCase(t)
	mockLog := setupMockLogger(t)

	la, err := createTestWorker(mockReceiver, mockUC, mockLog, 2)
	require.NoError(t, err)

	testDTO := &dto.LoginAttemptWorkerJob{
		Message: azure.NewMessage([]byte(`{}`), nil),
		Data:    &dto.LoginAttemptEventDTO{},
	}

	criticalDBErr := errors.New("postgres connection timeout")
	mockUC.On("Audit", mock.Anything, mock.Anything).Return(criticalDBErr).Once()

	// Исправлено: Release принимает 2 аргумента
	mockReceiver.On("Release", mock.Anything, testDTO.Message).Return(nil).Once()

	err = la.storeAuthAudit(context.Background(), 1, testDTO)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "store auth audit failed")
}
