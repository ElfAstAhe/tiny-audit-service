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
	worker2 "github.com/ElfAstAhe/tiny-audit-service/pkg/transport/worker"
	usecase2 "github.com/ElfAstAhe/tiny-audit-service/pkg/usecase"
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
	authAuditTailGetUCInst, err := container.GetInstance[usecase2.TailGetUseCase[string]](InstanceAuthAuditTailGetUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditTailCutUCInst, err := container.GetInstance[usecase2.TailCutUseCase[string]](InstanceAuthAuditTailCutUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}

	res, err := worker2.NewTailCutter(
		worker2.WithTailCutName("auth"),
		worker2.WithTailCutCutEnabled(confInst.AuthTC.TailCut),
		worker2.WithTailCutDataInterval(confInst.AuthTC.TailInterval),
		worker2.WithTailCutTailGetUC(authAuditTailGetUCInst),
		worker2.WithTailCutterTailCutUC(authAuditTailCutUCInst),
		worker2.WithTailCutStopTimeout(confInst.AuthTC.StopTimeout),
		worker2.WithTailCutWorkerCount(confInst.AuthTC.WorkerCount),
		worker2.WithTailCutDataCapacity(confInst.AuthTC.DataCapacity),
		worker2.WithTailCutCompleteProcess(confInst.AuthTC.CompleteProcess),
		worker2.WithTailCutStartInterval(confInst.AuthTC.StartInterval),
		worker2.WithTailCutScheduleInterval(confInst.AuthTC.ScheduleInterval),
		worker2.WithTailCutLogger(logInst),
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
	dataAuditTailGetUCInst, err := container.GetInstance[usecase2.TailGetUseCase[string]](InstanceDataAuditTailGetUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditTailCutUCInst, err := container.GetInstance[usecase2.TailCutUseCase[string]](InstanceDataAuditTailCutUC)
	if err != nil {
		return nil, errs.NewContainerError(wc.GetName(), "provider: retrieve instance failed", err)
	}

	res, err := worker2.NewTailCutter(
		worker2.WithTailCutName("data"),
		worker2.WithTailCutCutEnabled(confInst.DataTC.TailCut),
		worker2.WithTailCutDataInterval(confInst.DataTC.TailInterval),
		worker2.WithTailCutTailGetUC(dataAuditTailGetUCInst),
		worker2.WithTailCutterTailCutUC(dataAuditTailCutUCInst),
		worker2.WithTailCutStopTimeout(confInst.DataTC.StopTimeout),
		worker2.WithTailCutWorkerCount(confInst.DataTC.WorkerCount),
		worker2.WithTailCutDataCapacity(confInst.DataTC.DataCapacity),
		worker2.WithTailCutCompleteProcess(confInst.DataTC.CompleteProcess),
		worker2.WithTailCutStartInterval(confInst.DataTC.StartInterval),
		worker2.WithTailCutScheduleInterval(confInst.DataTC.ScheduleInterval),
		worker2.WithTailCutLogger(logInst),
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
