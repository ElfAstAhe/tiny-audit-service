package container

import (
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/broker/amqp"
	libworker "github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
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

	res, err := worker.NewTailCutter(
		worker.WithTailCutName("auth"),
		worker.WithTailCutCutEnabled(confInst.AuthTC.TailCut),
		worker.WithTailCutDataInterval(confInst.AuthTC.TailInterval),
		worker.WithTailCutTailGetUC(authAuditTailGetUCInst),
		worker.WithTailCutterTailCutUC(authAuditTailCutUCInst),
		worker.WithTailCutStopTimeout(confInst.AuthTC.StopTimeout),
		worker.WithTailCutWorkerCount(confInst.AuthTC.WorkerCount),
		worker.WithTailCutDataCapacity(confInst.AuthTC.DataCapacity),
		worker.WithTailCutCompleteProcess(confInst.AuthTC.CompleteProcess),
		worker.WithTailCutStartInterval(confInst.AuthTC.StartInterval),
		worker.WithTailCutScheduleInterval(confInst.AuthTC.ScheduleInterval),
		worker.WithTailCutLogger(logInst),
	)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceAuthAuditTailCutter), err)
	}

	return res, nil
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

	res, err := worker.NewTailCutter(
		worker.WithTailCutName("data"),
		worker.WithTailCutCutEnabled(confInst.DataTC.TailCut),
		worker.WithTailCutDataInterval(confInst.DataTC.TailInterval),
		worker.WithTailCutTailGetUC(dataAuditTailGetUCInst),
		worker.WithTailCutterTailCutUC(dataAuditTailCutUCInst),
		worker.WithTailCutStopTimeout(confInst.DataTC.StopTimeout),
		worker.WithTailCutWorkerCount(confInst.DataTC.WorkerCount),
		worker.WithTailCutDataCapacity(confInst.DataTC.DataCapacity),
		worker.WithTailCutCompleteProcess(confInst.DataTC.CompleteProcess),
		worker.WithTailCutStartInterval(confInst.DataTC.StartInterval),
		worker.WithTailCutScheduleInterval(confInst.DataTC.ScheduleInterval),
		worker.WithTailCutLogger(logInst),
	)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceDataAuditTailCutter), err)
	}

	return res, nil
}

//goland:noinspection DuplicatedCode
func (wc *WorkerContainer) providerLoginAttemptsListener() (any, error) {
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
	receiverInst, err := wc.getLoginAttemptsReceiver(confInst.LoginAttempts.ReceiverKind)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}

	res, err := worker.NewLoginAttempts(
		worker.WithLAOName("login-attempts-listener"),
		worker.WithLAOReceiver(receiverInst),
		worker.WithLAOAuthAuditUseCase(authAuditUCInst),
		worker.WithLAOLogger(logInst),
		worker.WithLAOBatchSize(confInst.LoginAttempts.BatchSize),
		worker.WithLAOBatchReadTimeout(confInst.LoginAttempts.BatchReadTimeout),
		worker.WithLAOAcknowledgeTimeout(confInst.LoginAttempts.AcknowledgeTimeout),
		worker.WithLAODispatcherOpts(
			libworker.WithSchedulerDispatcherStopTimeout[*dto.LoginAttemptWorkerJob](confInst.LoginAttempts.ShutdownTimeout),
			libworker.WithSchedulerDispatcherPoolWorkerCount[*dto.LoginAttemptWorkerJob](confInst.LoginAttempts.WorkerCount),
			libworker.WithSchedulerDispatcherPoolDataCapacity[*dto.LoginAttemptWorkerJob](confInst.LoginAttempts.DataCapacity),
			libworker.WithSchedulerDispatcherPoolCompleteProcess[*dto.LoginAttemptWorkerJob](confInst.LoginAttempts.CompleteProcessing),
			libworker.WithSchedulerDispatcherSchedulerStartInterval[*dto.LoginAttemptWorkerJob](confInst.LoginAttempts.StartInterval),
			libworker.WithSchedulerDispatcherSchedulerScheduleInterval[*dto.LoginAttemptWorkerJob](confInst.LoginAttempts.ScheduleInterval),
		),
	)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceLoginAttemptsListener), err)
	}

	return res, nil
}

func (wc *WorkerContainer) getLoginAttemptsReceiver(receiverKind string) (broker.Receiver, error) {
	switch receiverKind {
	case "amqp":
		return container.GetInstance[libamqp.Receiver](InstanceLoginAttemptsAMQPReceiver)
	case "kafka":
		return container.GetInstance[broker.Receiver](InstanceLoginAttemptsKafkaReceiver)
	default:
		return nil, errs.NewContainerError(wc.GetName(), fmt.Sprintf("provider: unknown receiver kind %s", receiverKind), nil)
	}
}
