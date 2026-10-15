package client

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	mocks2 "github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Вспомогательная функция для настройки базовых моков логгера,
// так как BaseAuditClient активно пишет логи.
func setupMockLogger(t *testing.T) logger.Logger {
	mockLogger := mocks.NewMockLogger(t)
	// Конструктор делает GetLogger, возвращаем этот же мок
	mockLogger.On("GetLogger", mock.Anything).Return(mockLogger).Maybe()
	// Игнорируем логирование, чтобы не засорять тест-кейс
	mockLogger.On("Debug", mock.Anything).Return().Maybe()
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return().Maybe()
	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLogger.On("Warn", mock.Anything).Return().Maybe()
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Return().Maybe()

	return mockLogger
}

// 1. Тестируем успешный путь прохождения данных через пул воркеров
func TestBaseAuditClient_Success(t *testing.T) {
	mockLogger := setupMockLogger(t)

	mockTokenProvider := mocks2.NewMockTokenProvider(t)
	mockTokenProvider.On("GetAccessToken").Return("super-secret-token", nil)

	// Сигнализируем основному потоку теста, что асинхронный экшен выполнился
	actionDone := make(chan struct{})

	var dummyAction AuditAction[*dto.AuthAuditDTO] = func(ctx context.Context, workerIndex int, data *dto.AuthAuditDTO, token string) error {
		assert.Equal(t, "super-secret-token", token)
		assert.Equal(t, "user-login-event", data.Event) // пример поля DTO
		close(actionDone)
		return nil
	}

	opts := []BaseAuditClientOption[*dto.AuthAuditDTO]{
		WithName[*dto.AuthAuditDTO]("test-pool"),
		WithLogger[*dto.AuthAuditDTO](mockLogger),
		WithTokenProvider[*dto.AuthAuditDTO](mockTokenProvider),
		WithAuditAction[*dto.AuthAuditDTO](dummyAction),
		WithPoolWorkerCount[*dto.AuthAuditDTO](1),
		WithPoolDataCapacity[*dto.AuthAuditDTO](10),
		WithPoolStopTimeout[*dto.AuthAuditDTO](time.Second),
	}

	bac, err := NewBaseAuditClient(opts...)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем пул воркеров
	err = bac.Start(ctx)
	require.NoError(t, err)
	defer func() { _ = bac.Stop(ctx) }()

	// Отправляем тестовые данные в аудит
	testData := &dto.AuthAuditDTO{Event: "user-login-event"}
	err = bac.Audit(testData)
	assert.NoError(t, err)

	// Ждем выполнения асинхронного экшена или таймаута
	select {
	case <-actionDone:
		// Успех
	case <-ctx.Done():
		t.Fatal("timeout waiting for audit action execution")
	}

	assert.Equal(t, int32(0), bac.TotalLost())
}

// 2. Тестируем переполнение очереди: когда TryPush возвращает false
func TestBaseAuditClient_QueueOverflow(t *testing.T) {
	mockLogger := setupMockLogger(t)
	mockTokenProvider := mocks2.NewMockTokenProvider(t)

	// Каналы для контролируемой блокировки воркера
	blockWorker := make(chan struct{})
	workerStarted := make(chan struct{})

	// Используем sync.Once, чтобы закрыть канал workerStarted только один раз
	var signalOnce sync.Once

	// Экшен, который зависнет при обработке первого элемента
	var blockingAction AuditAction[*dto.AuthAuditDTO] = func(ctx context.Context, workerIndex int, data *dto.AuthAuditDTO, token string) error {
		signalOnce.Do(func() {
			close(workerStarted) // Сигнализируем только при первом запуске воркера
		})
		<-blockWorker // Блокируем воркер навсегда, пока тест не разрешит выйти
		return nil
	}

	opts := []BaseAuditClientOption[*dto.AuthAuditDTO]{
		WithName[*dto.AuthAuditDTO]("overflow-pool"),
		WithLogger[*dto.AuthAuditDTO](mockLogger),
		WithTokenProvider[*dto.AuthAuditDTO](mockTokenProvider),
		WithAuditAction[*dto.AuthAuditDTO](blockingAction),
		WithPoolWorkerCount[*dto.AuthAuditDTO](1),
		WithPoolDataCapacity[*dto.AuthAuditDTO](1), // Буфер канала = 1
		WithPoolStopTimeout[*dto.AuthAuditDTO](time.Second),
	}

	bac, err := NewBaseAuditClient(opts...)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Запускаем пул, чтобы воркер начал слушать канал
	err = bac.Start(ctx)
	require.NoError(t, err)

	// Настраиваем токен, так как воркер будет его запрашивать для каждого элемента
	mockTokenProvider.On("GetAccessToken").Return("token", nil).Maybe()

	// 1. Этот элемент сразу заберет запущенный воркер в обработку и зависнет
	err = bac.Audit(&dto.AuthAuditDTO{Event: "task-1"})
	assert.NoError(t, err)

	// Ждем, пока воркер гарантированно займет задачу 1
	select {
	case <-workerStarted:
	case <-ctx.Done():
		t.Fatal("timeout waiting for worker to start blocking action")
	}

	// 2. Этот элемент полностью забьет буфер канала (ёмкость 1)
	err = bac.Audit(&dto.AuthAuditDTO{Event: "task-2"})
	assert.NoError(t, err)

	// 3. Третий элемент уже не поместится ни в воркер, ни в буфер канала.
	// TryPush вернет false, зафиксировав утерю данных.
	err = bac.Audit(&dto.AuthAuditDTO{Event: "task-3-overflow"})

	// Проверяем, что зафиксирована ошибка переполнения очереди
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to push auth audit data")

	// Счетчик должен показать ровно 1 гарантированно потерянный элемент (task-3)
	assert.Equal(t, int32(1), bac.TotalLost())

	// Разблокируем воркер, чтобы пул мог без паники переварить оставшийся в буфере task-2
	close(blockWorker)

	// Даем воркеру немного времени завершить работу перед Stop
	time.Sleep(20 * time.Millisecond)
	_ = bac.Stop(ctx)
}

// 3. Тестируем сценарий, когда TokenProvider падает с ошибкой
func TestBaseAuditClient_TokenProviderError(t *testing.T) {
	mockLogger := setupMockLogger(t)

	mockTokenProvider := mocks2.NewMockTokenProvider(t)
	mockTokenProvider.On("GetAccessToken").Return("", errors.New("vault connection refused"))

	// Если токен упал, jobHandler сразу выходит, увеличивая счетчик
	opts := []BaseAuditClientOption[*dto.AuthAuditDTO]{
		WithName[*dto.AuthAuditDTO]("token-error-pool"),
		WithLogger[*dto.AuthAuditDTO](mockLogger),
		WithTokenProvider[*dto.AuthAuditDTO](mockTokenProvider),
		WithAuditAction[*dto.AuthAuditDTO](func(ctx context.Context, wIdx int, d *dto.AuthAuditDTO, tkn string) error { return nil }),
		WithPoolWorkerCount[*dto.AuthAuditDTO](1),
		WithPoolDataCapacity[*dto.AuthAuditDTO](10),
		WithPoolStopTimeout[*dto.AuthAuditDTO](time.Second),
	}

	bac, err := NewBaseAuditClient(opts...)
	require.NoError(t, err)

	// Мы можем вызвать внутренний jobHandler вручную (unit-тест),
	// чтобы не плодить асинхронное состояние
	ctx := context.Background()
	err = bac.jobHandler(ctx, 0, &dto.AuthAuditDTO{}) // Предполагается, что вы экспортировали JobHandler для тестов или тестируете через пул

	// Проверяем инкремент счетчика потерянных данных
	assert.Equal(t, int32(1), bac.TotalLost())
}
