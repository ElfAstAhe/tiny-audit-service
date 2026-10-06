package worker

import (
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

// Константы со значениями по умолчанию для конфигурации планировщика TailCutter.
const (
	defaultTailCutterWorkerCount      int           = 1
	defaultTailCutterDataCapacity     int           = 32
	defaultTailCutterCompleteProcess  bool          = false
	defaultTailCutterStartInterval    time.Duration = 30 * time.Second
	defaultTailCutterScheduleInterval time.Duration = 5 * time.Minute
	defaultTailCutterDataInterval     time.Duration = 14 * 24 * time.Hour
	defaultTailCutterCutEnabled       bool          = false
)

// TailCutterOption определяет сигнатуру функции для настройки параметров TailCutter.
type TailCutterOption func(*TailCutterOptions)

// TailCutterOptions содержит конфигурационные параметры для планировщика обрезки хвостов данных.
// Включает в себя базовые опции диспетчера через эмбеддинг BaseSchedulerDispatcherOptions.
type TailCutterOptions struct {
	*worker.BaseSchedulerDispatcherOptions[string]
	TailGetUC    usecase.TailGetUseCase[string]
	TailCutUC    usecase.TailCutUseCase[string]
	DataInterval time.Duration
	CutEnabled   bool
}

// NewTailCutterOptions создает и возвращает новый экземпляр TailCutterOptions
// с базовыми настройками диспетчера и значениями интервалов по умолчанию.
func NewTailCutterOptions() *TailCutterOptions {
	res := &TailCutterOptions{
		BaseSchedulerDispatcherOptions: tailCutterBuildDefaultSDOptions(),
		DataInterval:                   defaultTailCutterDataInterval,
		CutEnabled:                     defaultTailCutterCutEnabled,
	}

	return res
}

// Validate проверяет корректность заполнения всех обязательных полей конфигурации,
// включая встроенные параметры базового диспетчера и специфичные параметры use case.
//
//goland:noinspection DuplicatedCode
func (tco *TailCutterOptions) Validate() error {
	if strings.TrimSpace(tco.Name) == "" {
		return errs.NewTlCommonError("Validate", "name is required", nil)
	}
	if tco.WorkerCount <= 0 {
		return errs.NewTlCommonError("Validate", "worker count is required", nil)
	}
	if tco.DataCapacity <= 0 {
		return errs.NewTlCommonError("Validate", "data capacity is required", nil)
	}
	if tco.StopTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "stop timeout is required", nil)
	}
	//if utils.IsNil(tco.JobHandler) {
	//	return errs.NewTlCommonError("Validate", "job handler is required", nil)
	//}
	if utils.IsNil(tco.Logger) {
		return errs.NewTlCommonError("Validate", "logger is required", nil)
	}
	//if utils.IsNil(tco.DataProvider) {
	//	return errs.NewTlCommonError("Validate", "data provider is required", nil)
	//}
	if tco.StartInterval <= 0 {
		return errs.NewTlCommonError("Validate", "start interval is required", nil)
	}
	if tco.ScheduleInterval <= 0 {
		return errs.NewTlCommonError("Validate", "schedule interval is required", nil)
	}
	if utils.IsNil(tco.TailGetUC) {
		return errs.NewTlCommonError("Validate", "tail get use case is required", nil)
	}
	if utils.IsNil(tco.TailCutUC) {
		return errs.NewTlCommonError("Validate", "tail cut use case is required", nil)
	}
	if tco.DataInterval <= 0 {
		return errs.NewTlCommonError("Validate", "data interval is required", nil)
	}

	return nil
}

// WithTailCutName задает имя планировщика TailCutter.
func WithTailCutName(name string) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.Name = name
	}
}

// WithTailCutWorkerCount задает количество параллельных воркеров для пула обработки.
func WithTailCutWorkerCount(workerCount int) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.WorkerCount = workerCount
	}
}

// WithTailCutDataCapacity устанавливает емкость очереди данных для обработки.
func WithTailCutDataCapacity(capacity int) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.DataCapacity = capacity
	}
}

// WithTailCutCompleteProcess управляет флагом завершения всех задач в очереди при остановке.
func WithTailCutCompleteProcess(flag bool) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.CompleteProcess = flag
	}
}

// WithTailCutStartInterval определяет интервал задержки перед первым запуском планировщика.
func WithTailCutStartInterval(interval time.Duration) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.StartInterval = interval
	}
}

// WithTailCutScheduleInterval определяет интервал периодического запуска задач планировщиком.
func WithTailCutScheduleInterval(interval time.Duration) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.ScheduleInterval = interval
	}
}

// WithTailCutStopTimeout задает максимальное время ожидания завершения работы при остановке.
func WithTailCutStopTimeout(timeout time.Duration) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.StopTimeout = timeout
	}
}

// WithTailCutDataInterval устанавливает временной интервал (глубину), за который отсекаются данные.
func WithTailCutDataInterval(interval time.Duration) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.DataInterval = interval
	}
}

// WithTailCutCutEnabled включает или выключает физическое удаление/обрезание данных.
func WithTailCutCutEnabled(enabled bool) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.CutEnabled = enabled
	}
}

// WithTailCutTailGetUC конфигурирует Use Case для получения "хвостов" данных.
func WithTailCutTailGetUC(uc usecase.TailGetUseCase[string]) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.TailGetUC = uc
	}
}

// WithTailCutterTailCutUC конфигурирует Use Case для обрезания/удаления данных.
func WithTailCutterTailCutUC(uc usecase.TailCutUseCase[string]) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.TailCutUC = uc
	}
}

// WithTailCutLogger настраивает logger
func WithTailCutLogger(log logger.Logger) TailCutterOption {
	return func(options *TailCutterOptions) {
		options.Logger = log
	}
}

func tailCutterBuildDefaultSDOptions() *worker.BaseSchedulerDispatcherOptions[string] {
	res := worker.NewBaseSchedulerDispatcherOptions[string]()
	for _, option := range tailCutterDefaultSDOptionList() {
		option(res)
	}

	return res
}

func tailCutterDefaultSDOptionList() []worker.BaseSchedulerDispatcherOption[string] {
	return []worker.BaseSchedulerDispatcherOption[string]{
		worker.WithSchedulerDispatcherPoolWorkerCount[string](defaultTailCutterWorkerCount),
		worker.WithSchedulerDispatcherPoolDataCapacity[string](defaultTailCutterDataCapacity),
		worker.WithSchedulerDispatcherPoolCompleteProcess[string](defaultTailCutterCompleteProcess),
		worker.WithSchedulerDispatcherSchedulerStartInterval[string](defaultTailCutterStartInterval),
		worker.WithSchedulerDispatcherSchedulerScheduleInterval[string](defaultTailCutterScheduleInterval),
	}
}
