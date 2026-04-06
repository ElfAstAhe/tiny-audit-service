package worker

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type TailCutterConfig struct {
	*worker.BaseSchedulerDispatcherConfig
	dataInterval time.Duration
	cutEnabled   bool
}

func NewTailCutterConfig(
	startInterval time.Duration,
	scheduleInterval time.Duration,
	workerCount int,
	dataCapacity int,
	dataInterval time.Duration,
	cutEnabled bool,
) *TailCutterConfig {
	return &TailCutterConfig{
		BaseSchedulerDispatcherConfig: worker.NewBaseSchedulerDispatcherConfig(
			worker.NewBaseSchedulerConfig(startInterval, scheduleInterval),
			workerCount,
			dataCapacity,
		),
		dataInterval: dataInterval,
		cutEnabled:   cutEnabled,
	}
}

type TailCutter struct {
	*worker.BaseSchedulerDispatcher[string]
	config    *TailCutterConfig
	tailGetUC usecase.TailGetUseCase[string]
	tailCutUC usecase.TailCutUseCase[string]
}

func NewTailCutter(
	parentCtx context.Context,
	name string,
	config *TailCutterConfig,
	tailGetUC usecase.TailGetUseCase[string],
	tailCutUC usecase.TailCutUseCase[string],
	log logger.Logger,
) *TailCutter {
	res := &TailCutter{
		config:    config,
		tailGetUC: tailGetUC,
		tailCutUC: tailCutUC,
	}

	base := worker.NewBaseSchedulerDispatcher[string](
		name,
		parentCtx,
		config.BaseSchedulerDispatcherConfig,
		res.dataProvider,
		res.cutTail,
		log,
	)

	res.BaseSchedulerDispatcher = base

	return res
}

func (tc *TailCutter) dataProvider(ctx context.Context, eventTime time.Time) ([]string, error) {
	tc.GetLogger().Debugf("start %s tail cutter data provider", tc.GetName())
	defer tc.GetLogger().Debugf("finish %s tail cutter data provider", tc.GetName())

	tailCutTime := eventTime.Add(-tc.config.dataInterval)
	tc.GetLogger().Debugf("tail cut time %s tail cutter data provider [%s]", tc.GetName(), tailCutTime.Format(time.DateTime))

	if !tc.config.cutEnabled {
		tc.GetLogger().Debugf("tail cut %s tail cut enabled [%v] pass iteration", tc.GetName(), tc.config.cutEnabled)

		return []string{}, nil
	}

	return tc.tailGetUC.GetTail(ctx, tailCutTime)
}

func (tc *TailCutter) cutTail(ctx context.Context, workerIndex int, data string) error {
	tc.GetLogger().Debugf("start %s tail cutter cutter worker %v job handler", tc.GetName(), workerIndex)
	defer tc.GetLogger().Debugf("finish %s tail cutter worker %v cutter job handler", tc.GetName(), workerIndex)

	return tc.tailCutUC.Cut(ctx, data)
}

func (tc *TailCutter) GetConfig() *TailCutterConfig {
	return tc.config
}
