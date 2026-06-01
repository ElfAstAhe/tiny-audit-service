package app

import (
	"context"
	"errors"

	_ "expvar"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

type Orchestrator struct {
	*container.BaseOrchestrator
	conf   *config.Config
	logger logger.Logger
}

var _ container.Orchestrator = (*Orchestrator)(nil)

func NewOrchestrator(conf *config.Config, log logger.Logger) *Orchestrator {
	return &Orchestrator{
		BaseOrchestrator: container.NewBaseOrchestrator(log),
		conf:             conf,
		logger:           log,
	}
}

func (o *Orchestrator) Init(ctx context.Context) error {
	appCnt, err := o.GetContainer(AppContainerName)
	if err != nil {
		return errs.NewContainerError(OrchestratorName, "init failed", err)
	}
	err = errors.Join(
		appCnt.RegisterInstance(InstanceConfig, o.conf),
		appCnt.RegisterInstance(InstanceLogger, o.logger),
	)
	if err != nil {
		return errs.NewContainerError(OrchestratorName, "init failed", err)
	}

	return o.BaseOrchestrator.Init(ctx)
}
