package container

import (
	"fmt"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp/azure"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp/kafka"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

func (cc *ClientContainer) providerLoginAttemptsKafkaReceiver() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	receiver, err := kafka.NewReceiver(
		kafka.WithReceiverClientID(confInst.App.NodeName),
		kafka.WithReceiverBrokers(confInst.LoginAttempts.KafkaConfig.Brokers),
		kafka.WithReceiverTargetName(confInst.LoginAttempts.KafkaConfig.TargetName),
		kafka.WithReceiverGroupID(confInst.LoginAttempts.KafkaConfig.GroupID),
		kafka.WithReceiverPartition(confInst.LoginAttempts.KafkaConfig.Partition),
		kafka.WithReceiverConnectTimeout(confInst.LoginAttempts.KafkaConfig.ConnectTimeout),
		kafka.WithReceiverShutdownTimeout(confInst.LoginAttempts.KafkaConfig.ShutdownTimeout),
		kafka.WithReceiverLogger(logInst),
		kafka.WithReceiverGroupTimeouts(
			confInst.LoginAttempts.KafkaConfig.HeartbeatInterval,
			confInst.LoginAttempts.KafkaConfig.SessionTimeout,
			confInst.LoginAttempts.KafkaConfig.RebalanceTimeout,
			confInst.LoginAttempts.KafkaConfig.ReadTimeout,
		),
		kafka.WithReceiverRuntimePerformance(
			confInst.LoginAttempts.KafkaConfig.MaxAttempts,
			confInst.LoginAttempts.KafkaConfig.QueueCapacity,
			confInst.LoginAttempts.KafkaConfig.StartOffset,
		),
		kafka.WithReceiverSecurity(
			confInst.LoginAttempts.KafkaConfig.Username,
			confInst.LoginAttempts.KafkaConfig.Password,
		),
		kafka.WithReceiverFetchBounds(
			confInst.LoginAttempts.KafkaConfig.MinBytes,
			confInst.LoginAttempts.KafkaConfig.MaxBytes,
			confInst.LoginAttempts.KafkaConfig.MaxWait,
		),
	)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceLoginAttemptsKafkaReceiver), err)
	}

	return receiver, nil
}

func (cc *ClientContainer) providerLoginAttemptsAMQPReceiver() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	connectorInst, err := container.GetInstance[libamqp.Connector[*amqp.Session]](InstanceAMQPConnector)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	receiverConfInst, err := container.GetInstance[*amqp.ReceiverOptions](InstanceLoginAttemptsAMQPReceiverReceiverOpts)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	receiver, err := azure.NewReceiver(
		azure.WithReceiverConnector(connectorInst),
		azure.WithReceiverTargetName(confInst.LoginAttempts.AMQPConfig.TargetName),
		azure.WithReceiverLinkCredit(int32(confInst.LoginAttempts.AMQPConfig.PrefetchCredit)),
		azure.WithReceiverConnectTimeout(confInst.LoginAttempts.AMQPConfig.ConnectTimeout),
		azure.WithReceiverShutdownTimeout(confInst.LoginAttempts.AMQPConfig.ShutdownTimeout),
		azure.WithReceiverOpts(receiverConfInst),
		azure.WithReceiverLogger(logInst),
	)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceLoginAttemptsAMQPReceiver), err)
	}

	return receiver, nil
}

func (cc *ClientContainer) providerLoginAttemptsAMQPReceiverReceiverOpts() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	return &amqp.ReceiverOptions{
		Name:   confInst.App.NodeName,
		Credit: int32(confInst.LoginAttempts.AMQPConfig.PrefetchCredit),
	}, nil
}

func (cc *ClientContainer) providerAMQPConnector() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	connOpts, err := container.GetInstance[*amqp.ConnOptions](InstanceAMQPConnectorConnOpts)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	sessOpts, err := container.GetInstance[*amqp.SessionOptions](InstanceAMQPConnectorSessOpts)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	connectorInst, err := azure.NewConnector(
		azure.WithConnectorURL(confInst.AMQPConnector.URL),
		azure.WithConnectorConnectTimeout(confInst.AMQPConnector.ConnectTimeout),
		azure.WithConnectorShutdownTimeout(confInst.AMQPConnector.ShutdownTimeout),
		azure.WithConnectorConnOpts(connOpts),
		azure.WithConnectorSessionOpts(sessOpts),
		azure.WithConnectorLogger(logInst),
	)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceAMQPConnector), err)
	}

	return connectorInst, nil
}

func (cc *ClientContainer) providerAMQPConnectorConnOpts() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	return &amqp.ConnOptions{
		HostName:     confInst.App.NodeName,
		IdleTimeout:  confInst.AMQPConnector.IdleTimeout,
		WriteTimeout: confInst.AMQPConnector.WriteTimeout,
		SASLType:     amqp.SASLTypePlain(confInst.AMQPConnector.Username, confInst.AMQPConnector.Password),
	}, nil
}

func (cc *ClientContainer) providerAMQPConnectorSessOpts() (any, error) {
	return &amqp.SessionOptions{
		MaxLinks: 4,
	}, nil
}
