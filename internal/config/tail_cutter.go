package config

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

type TailCutterConfig struct {
	*config.SchedulerDispatcherConfig `mapstructure:",squash"`
	TailInterval                      time.Duration `mapstructure:"tail_interval" json:"tail_interval,omitempty" yaml:"tail_interval,omitempty"`
	TailCut                           bool          `mapstructure:"tail_cut" json:"tail_cut,omitempty" yaml:"tail_cut,omitempty"`
}

func NewTailCutterConfig(
	startInterval time.Duration,
	scheduleInterval time.Duration,
	workerCount int,
	dataCapacity int,
	completeProcessing bool,
	shutdownTimeout time.Duration,
	tailInterval time.Duration,
	tailCut bool,
) *TailCutterConfig {
	return &TailCutterConfig{
		SchedulerDispatcherConfig: config.NewSchedulerDispatcherConfig(
			workerCount,
			dataCapacity,
			completeProcessing,
			startInterval,
			scheduleInterval,
			shutdownTimeout,
		),
		TailInterval: tailInterval,
		TailCut:      tailCut,
	}
}

func NewDefaultTailCutterConfig() *TailCutterConfig {
	return NewTailCutterConfig(
		defaultAuthTCStartInterval,
		defaultAuthTCScheduleInterval,
		defaultAuthTCWorkerCount,
		defaultAuthTCDataCapacity,
		defaultAuthTCCompleteProcessing,
		defaultAuthTCStopTimeout,
		defaultAuthTCTailInterval,
		defaultAuthTCTailCut,
	)
}

func (tcc *TailCutterConfig) Validate() error {
	if tcc.StartInterval <= 0 {
		return errs.NewConfigValidateError("tail cutter", "start interval", "must be greater zero", nil)
	}
	if tcc.ScheduleInterval <= 0 {
		return errs.NewConfigValidateError("tail cutter", "schedule interval", "must be greater zero", nil)
	}
	if tcc.WorkerCount <= 0 {
		return errs.NewConfigValidateError("tail cutter", "worker count", "must be greater zero", nil)
	}
	if tcc.DataCapacity <= 0 {
		return errs.NewConfigValidateError("tail cutter", "data capacity", "must be greater zero", nil)
	}
	if tcc.StopTimeout <= 0 {
		return errs.NewConfigValidateError("tail cutter", "stop timeout", "must be greater zero", nil)
	}
	if tcc.TailInterval <= 0 {
		return errs.NewConfigValidateError("tail cutter", "tail interval", "must be greater zero", nil)
	}

	return nil
}
