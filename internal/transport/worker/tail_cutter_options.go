package worker

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
)

type TailCutterOptions struct {
	*worker.BaseSchedulerDispatcherConfig
	dataInterval time.Duration
	cutEnabled   bool
}

func NewTailCutterOptions(
	schedulerDispatcherConf *worker.BaseSchedulerDispatcherConfig,
	dataInterval time.Duration,
	cutEnabled bool,
) *TailCutterOptions {
	return &TailCutterOptions{
		BaseSchedulerDispatcherConfig: worker.NewBaseSchedulerDispatcherConfig(
			schedulerDispatcherConf.SchedulerConfig,
			schedulerDispatcherConf.PoolConfig,
		),
		dataInterval: dataInterval,
		cutEnabled:   cutEnabled,
	}
}
