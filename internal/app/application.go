package app

import (
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/app"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	"github.com/ElfAstAhe/tiny-audit-service/internal/tools"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/grpc"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/rest"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
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
		app.WithOrchestrator(NewOrchestrator(res.conf, res.log)),
		app.WithLogger(res.log),
		app.WithCloseTimeout(res.conf.App.CloseTimeout),
		app.WithStopTimeout(res.conf.App.StopTimeout),
	)
	// orchestrator and containers
	err := errors.Join(
		// app container
		res.GetOrchestrator().Register(NewContainer(res.GetOrchestrator())),
		// tools container
		res.GetOrchestrator().Register(tools.NewContainer(res.GetOrchestrator())),
		// postgres container
		res.GetOrchestrator().Register(postgres.NewContainer(res.GetOrchestrator())),
		// repository container
		res.GetOrchestrator().Register(repository.NewContainer(res.GetOrchestrator())),
		// use case container
		res.GetOrchestrator().Register(usecase.NewContainer(res.GetOrchestrator())),
		// facade container
		res.GetOrchestrator().Register(facade.NewContainer(res.GetOrchestrator())),
		// services container (inner kitchen)
		// res.GetOrchestrator().Register(container.NewServiceContainer(res.GetOrchestrator())),
		// worker container
		res.GetOrchestrator().Register(worker.NewContainer(res.GetOrchestrator())),
		// http container
		res.GetOrchestrator().Register(rest.NewContainer(res.GetOrchestrator())),
		// gRPC container
		res.GetOrchestrator().Register(grpc.NewContainer(res.GetOrchestrator())),
	)
	if err != nil {
		return nil, errs.NewCommonError("application create failed", err)
	}

	return res, nil
}

func (app *Application) Init() error {
	appCnt, err := app.GetOrchestrator().GetContainer(ContainerName)
	if err != nil {
		return errs.NewCommonError("app init", err)
	}
	err = errors.Join(
		appCnt.RegisterInstance(InstanceApplication, app),
		appCnt.RegisterInstance(InstanceApplicationReady, app.IsReady),
	)
	if err != nil {
		return errs.NewCommonError("app init", err)
	}

	return app.BaseApplication.Init()
}
