package config

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

type TailCutterConfig struct {
	StartInterval      time.Duration `mapstructure:"start_interval" json:"start_interval,omitempty" yaml:"start_interval,omitempty"`
	ScheduleInterval   time.Duration `mapstructure:"schedule_interval" json:"schedule_interval,omitempty" yaml:"schedule_interval,omitempty"`
	WorkerCount        int           `mapstructure:"worker_count" json:"worker_count,omitempty" yaml:"worker_count,omitempty"`
	DataCapacity       int           `mapstructure:"data_capacity" json:"data_capacity,omitempty" yaml:"data_capacity,omitempty"`
	CompleteProcessing bool          `mapstructure:"complete_processing" json:"complete_processing,omitempty" yaml:"complete_processing,omitempty"`
	ShutdownTimeout    time.Duration `mapstructure:"shutdown_timeout" json:"shutdown_timeout,omitempty" yaml:"shutdown_timeout,omitempty"`
	TailInterval       time.Duration `mapstructure:"tail_interval" json:"tail_interval,omitempty" yaml:"tail_interval,omitempty"`
	TailCut            bool          `mapstructure:"tail_cut" json:"tail_cut,omitempty" yaml:"tail_cut,omitempty"`
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
		StartInterval:      startInterval,
		ScheduleInterval:   scheduleInterval,
		WorkerCount:        workerCount,
		DataCapacity:       dataCapacity,
		CompleteProcessing: completeProcessing,
		ShutdownTimeout:    shutdownTimeout,
		TailInterval:       tailInterval,
		TailCut:            tailCut,
	}
}

func NewDefaultTailCutterConfig() *TailCutterConfig {
	return NewTailCutterConfig(
		defaultAuthTCStartInterval,
		defaultAuthTCScheduleInterval,
		defaultAuthTCWorkerCount,
		defaultAuthTCDataCapacity,
		defaultAuthTCCompleteProcessing,
		defaultAuthTCShutdownTimeout,
		defaultAuthTCTailInterval,
		defaultAuthTCTailCut,
	)
}

func (tcc *TailCutterConfig) Validate() error {
	if !(tcc.StartInterval > 0) {
		return errs.NewConfigValidateError("_tc", "start_interval", "must be greater zero", nil)
	}
	if !(tcc.ScheduleInterval > 0) {
		return errs.NewConfigValidateError("_tc", "schedule_interval", "must be greater zero", nil)
	}
	if !(tcc.WorkerCount > 0) {
		return errs.NewConfigValidateError("_tc", "worker_count", "must be greater zero", nil)
	}
	if !(tcc.DataCapacity > 0) {
		return errs.NewConfigValidateError("_tc", "data_capacity", "must be greater zero", nil)
	}
	if !(tcc.ShutdownTimeout > 0) {
		return errs.NewConfigValidateError("_tc", "shutdown_timeout", "must be greater zero", nil)
	}
	if !(tcc.TailInterval > 0) {
		return errs.NewConfigValidateError("_tc", "tail_interval", "must be greater zero", nil)
	}

	return nil
}
