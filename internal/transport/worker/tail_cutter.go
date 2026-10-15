package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

const (
	tailCutterNameTemplate string = "tail-cutter-%s"
)

type TailCutter struct {
	*worker.BaseSchedulerDispatcher[string]
	name      string
	opts      *TailCutterOptions
	tailGetUC usecase.TailGetUseCase[string]
	tailCutUC usecase.TailCutUseCase[string]
	logger    logger.Logger
}

var _ worker.Scheduler = (*TailCutter)(nil)
var _ worker.CommonWorker = (*TailCutter)(nil)
var _ container.Runner = (*TailCutter)(nil)

func NewTailCutter(options ...TailCutterOption) (*TailCutter, error) {
	opts := NewTailCutterOptions()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, errs.NewTlCommonError("NewTailCutter", "tail cutter options validation failed", err)
	}
	// instance
	res := &TailCutter{
		name:      fmt.Sprintf(tailCutterNameTemplate, opts.Name),
		opts:      opts,
		tailGetUC: opts.TailGetUC,
		tailCutUC: opts.TailCutUC,
		logger:    opts.Logger.GetLogger(fmt.Sprintf(tailCutterNameTemplate, opts.Name)),
	}
	// scheduler dispatcher
	dispatcher, err := worker.NewBaseSchedulerDispatcher[string](
		worker.WithSchedulerDispatcherName[string](opts.Name),
		worker.WithSchedulerDispatcherStopTimeout[string](opts.StopTimeout),
		worker.WithSchedulerDispatcherDataProvider[string](res.dataProvider),
		worker.WithSchedulerDispatcherLogger[string](opts.Logger),
		worker.WithSchedulerDispatcherPoolWorkerCount[string](opts.WorkerCount),
		worker.WithSchedulerDispatcherPoolDataCapacity[string](opts.DataCapacity),
		worker.WithSchedulerDispatcherPoolCompleteProcess[string](opts.CompleteProcess),
		worker.WithSchedulerDispatcherPoolJobHandler[string](res.cutTail),
		worker.WithSchedulerDispatcherSchedulerStartInterval[string](opts.StartInterval),
		worker.WithSchedulerDispatcherSchedulerScheduleInterval[string](opts.ScheduleInterval),
	)
	if err != nil {
		return nil, errs.NewTlCommonError("NewTailCutter", "tail cutter scheduler dispatcher creation failed", err)
	}
	// setup
	res.BaseSchedulerDispatcher = dispatcher

	return res, nil
}

func (tc *TailCutter) dataProvider(ctx context.Context, eventTime time.Time) ([]string, error) {
	tc.GetLogger().Debugf("tail cutter %s time event %s data provider start", tc.GetName(), eventTime.Format(time.DateTime))
	defer tc.GetLogger().Debugf("tail cutter %s time event %s data provider finish", tc.GetName(), eventTime.Format(time.DateTime))

	tailCutTime := eventTime.Add(-tc.opts.DataInterval)
	tc.GetLogger().Debugf("tail cutter %s time event %s data provider tail time %s ", tc.GetName(), eventTime.Format(time.DateTime), tailCutTime.Format(time.DateTime))

	if !tc.opts.CutEnabled {
		tc.GetLogger().Debugf("tail cutter %s time event %s tail cut enabled [%v] pass iteration", tc.GetName(), eventTime.Format(time.DateTime), tc.opts.CutEnabled)

		return []string{}, nil
	}

	res, err := tc.tailGetUC.GetTail(ctx, tailCutTime)
	if err != nil {
		return nil, err
	}

	tc.GetLogger().Debugf("tail cutter %s time event %s total data records to dispatch [%v]", tc.GetName(), eventTime.Format(time.DateTime), len(res))

	return res, nil
}

func (tc *TailCutter) cutTail(ctx context.Context, workerIndex int, data string) error {
	tc.GetLogger().Debugf("tail cutter %s worker %v cut tail start", tc.GetName(), workerIndex)
	defer tc.GetLogger().Debugf("tail cutter %s worker %v cut tail finish", tc.GetName(), workerIndex)

	return tc.tailCutUC.Cut(ctx, data)
}

func (tc *TailCutter) GetName() string {
	return tc.name
}

func (tc *TailCutter) GetLogger() logger.Logger {
	return tc.logger
}

func (tc *TailCutter) GetOpts() *TailCutterOptions {
	return tc.opts
}
