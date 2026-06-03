package app

import (
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/app"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/app/container"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

type Application struct {
	*app.BaseApplication
	conf *config.Config
	log  logger.Logger
}

var _ app.Application = (*Application)(nil)

func NewApplication(opts ...Option) (*Application, error) {
	// create instance
	res := &Application{}
	// setup
	for _, opt := range opts {
		opt(res)
	}
	// embed
	res.BaseApplication = app.NewBaseApplication(
		app.WithOrchestrator(container.NewOrchestrator(res.conf, res.log)),
		app.WithLogger(res.log),
		app.WithCloseTimeout(res.conf.App.CloseTimeout),
		app.WithStopTimeout(res.conf.App.StopTimeout),
	)
	// orchestrator and containers
	err := errors.Join(
		// app container
		res.GetOrchestrator().Register(container.NewAppContainer(res.GetOrchestrator())),
		// tools container
		res.GetOrchestrator().Register(container.NewToolsContainer(res.GetOrchestrator())),
		// postgres container
		res.GetOrchestrator().Register(container.NewPgContainer(res.GetOrchestrator())),
		// repository container
		res.GetOrchestrator().Register(container.NewRepositoryContainer(res.GetOrchestrator())),
		// use case container
		res.GetOrchestrator().Register(container.NewUseCaseContainer(res.GetOrchestrator())),
		// facade container
		res.GetOrchestrator().Register(container.NewFacadeContainer(res.GetOrchestrator())),
		// services container (inner kitchen)
		res.GetOrchestrator().Register(container.NewServiceContainer(res.GetOrchestrator())),
		// worker container
		res.GetOrchestrator().Register(container.NewWorkerContainer(res.GetOrchestrator())),
		// http container
		res.GetOrchestrator().Register(container.NewHTTPContainer(res.GetOrchestrator())),
		// gRPC container
		res.GetOrchestrator().Register(container.NewGRPCContainer(res.GetOrchestrator())),
	)
	if err != nil {
		return nil, errs.NewCommonError("application create failed", err)
	}

	return res, nil
}

func (app *Application) Init() error {
	appCnt, err := app.GetOrchestrator().GetContainer(container.AppContainerName)
	if err != nil {
		return errs.NewCommonError("app init", err)
	}
	err = errors.Join(
		appCnt.RegisterInstance(container.InstanceApplication, app),
		appCnt.RegisterInstance(container.InstanceApplicationReady, app.IsReady),
	)
	if err != nil {
		return errs.NewCommonError("app init", err)
	}

	return app.BaseApplication.Init()
}
