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
	schedulerDispatcherConf *worker.BaseSchedulerDispatcherConfig,
	dataInterval time.Duration,
	cutEnabled bool,
) *TailCutterConfig {
	return &TailCutterConfig{
		BaseSchedulerDispatcherConfig: worker.NewBaseSchedulerDispatcherConfig(
			schedulerDispatcherConf.SchedulerConfig,
			schedulerDispatcherConf.PoolConfig,
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

var _ worker.Scheduler = (*TailCutter)(nil)
var _ worker.CommonWorker = (*TailCutter)(nil)

func NewTailCutter(
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
		config.BaseSchedulerDispatcherConfig,
		res.dataProvider,
		res.cutTail,
		log,
	)

	res.BaseSchedulerDispatcher = base

	return res
}

func (tc *TailCutter) dataProvider(ctx context.Context, eventTime time.Time) ([]string, error) {
	tc.GetLogger().Debugf("tail cutter %s time event %s data provider start", tc.GetName(), eventTime.Format(time.DateTime))
	defer tc.GetLogger().Debugf("tail cutter %s time event %s data provider finish", tc.GetName(), eventTime.Format(time.DateTime))

	tailCutTime := eventTime.Add(-tc.config.dataInterval)
	tc.GetLogger().Debugf("tail cutter %s time event %s data provider tail time %s ", tc.GetName(), eventTime.Format(time.DateTime), tailCutTime.Format(time.DateTime))

	if !tc.config.cutEnabled {
		tc.GetLogger().Debugf("tail cutter %s time event %s tail cut enabled [%v] pass iteration", tc.GetName(), eventTime.Format(time.DateTime), tc.config.cutEnabled)

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

func (tc *TailCutter) GetConfig() *TailCutterConfig {
	return tc.config
}
