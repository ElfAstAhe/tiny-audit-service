package container

import (
	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	libworker "github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

//goland:noinspection DuplicatedCode
func (wc *WorkerContainer) providerAuthAuditTailCutter() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditTailGetUCInst, err := container.GetInstance[usecase.TailGetUseCase[string]](InstanceAuthAuditTailGetUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditTailCutUCInst, err := container.GetInstance[usecase.TailCutUseCase[string]](InstanceAuthAuditTailCutUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}

	return worker.NewTailCutter(
		"auth-tail-cutter",
		worker.NewTailCutterOptions(
			libworker.NewBaseSchedulerDispatcherConfig(
				libworker.NewBaseSchedulerConfig(
					confInst.AuthTC.StartInterval,
					confInst.AuthTC.ScheduleInterval,
					confInst.AuthTC.ShutdownTimeout,
				),
				libworker.NewBasePoolConfig(
					confInst.AuthTC.WorkerCount,
					confInst.AuthTC.DataCapacity,
					confInst.AuthTC.CompleteProcessing,
					confInst.AuthTC.ShutdownTimeout,
				),
			),
			confInst.AuthTC.TailInterval,
			confInst.AuthTC.TailCut,
		),
		authAuditTailGetUCInst,
		authAuditTailCutUCInst,
		logInst,
	), nil
}

//goland:noinspection DuplicatedCode
func (wc *WorkerContainer) providerDataAuditTailCutter() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditTailGetUCInst, err := container.GetInstance[usecase.TailGetUseCase[string]](InstanceDataAuditTailGetUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditTailCutUCInst, err := container.GetInstance[usecase.TailCutUseCase[string]](InstanceDataAuditTailCutUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}

	return worker.NewTailCutter(
		"data-tail-cutter",
		worker.NewTailCutterOptions(
			libworker.NewBaseSchedulerDispatcherConfig(
				libworker.NewBaseSchedulerConfig(
					confInst.DataTC.StartInterval,
					confInst.DataTC.ScheduleInterval,
					confInst.DataTC.ShutdownTimeout,
				),
				libworker.NewBasePoolConfig(
					confInst.DataTC.WorkerCount,
					confInst.DataTC.DataCapacity,
					confInst.DataTC.CompleteProcessing,
					confInst.DataTC.ShutdownTimeout,
				),
			),
			confInst.DataTC.TailInterval,
			confInst.DataTC.TailCut,
		),
		dataAuditTailGetUCInst,
		dataAuditTailCutUCInst,
		logInst,
	), nil
}

//goland:noinspection DuplicatedCode
func (wc *WorkerContainer) providerLoginAttemptsAMQPListener() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditUCInst, err := container.GetInstance[usecase.AuthAuditUseCase](InstanceAuthAuditUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	receiverInst, err := container.GetInstance[libamqp.Receiver[*amqp.ReceiveOptions]](InstanceLoginAttemptsReceiver)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}

	return worker.NewLoginAttempts(
		worker.WithLAOName("login-attempts-amqp-listener"),
		worker.WithLAOReceiver(receiverInst),
		worker.WithLAOAuthAuditUseCase(authAuditUCInst),
		worker.WithLAOLogger(logInst),
		worker.WithLAOBatchSize(confInst.LoginAttempts.BatchSize),
		worker.WithLAOBatchReadTimeout(confInst.LoginAttempts.BatchReadTimeout),
		worker.WithLAOAcknowledgeTimeout(confInst.LoginAttempts.AcknowledgeTimeout),
		worker.WithLAOReceiveOpts(&amqp.ReceiveOptions{}),
		worker.WithLAODispatcherOpts(libworker.NewBaseSchedulerDispatcherConfig(
			libworker.NewBaseSchedulerConfig(
				confInst.LoginAttempts.StartInterval,
				confInst.LoginAttempts.ScheduleInterval,
				confInst.LoginAttempts.ShutdownTimeout,
			),
			libworker.NewBasePoolConfig(
				confInst.LoginAttempts.WorkerCount,
				confInst.LoginAttempts.DataCapacity,
				confInst.LoginAttempts.CompleteProcessing,
				confInst.LoginAttempts.ShutdownTimeout,
			),
		)),
	)
}
