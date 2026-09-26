package container

import (
	"fmt"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/broker/amqp"
	libkafka "github.com/ElfAstAhe/go-service-template/pkg/transport/broker/kafka"
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

	receiver, err := libkafka.NewReceiver(
		libkafka.WithReceiverClientID(confInst.App.NodeName),
		libkafka.WithReceiverBrokers(confInst.LoginAttempts.KafkaConfig.Brokers),
		libkafka.WithReceiverTargetName(confInst.LoginAttempts.KafkaConfig.TargetName),
		libkafka.WithReceiverGroupID(confInst.LoginAttempts.KafkaConfig.GroupID),
		libkafka.WithReceiverPartition(confInst.LoginAttempts.KafkaConfig.Partition),
		libkafka.WithReceiverConnectTimeout(confInst.LoginAttempts.KafkaConfig.ConnectTimeout),
		libkafka.WithReceiverShutdownTimeout(confInst.LoginAttempts.KafkaConfig.ShutdownTimeout),
		libkafka.WithReceiverLogger(logInst),
		libkafka.WithReceiverGroupTimeouts(
			confInst.LoginAttempts.KafkaConfig.HeartbeatInterval,
			confInst.LoginAttempts.KafkaConfig.SessionTimeout,
			confInst.LoginAttempts.KafkaConfig.RebalanceTimeout,
			confInst.LoginAttempts.KafkaConfig.ReadTimeout,
		),
		libkafka.WithReceiverRuntimePerformance(
			confInst.LoginAttempts.KafkaConfig.MaxAttempts,
			confInst.LoginAttempts.KafkaConfig.QueueCapacity,
			confInst.LoginAttempts.KafkaConfig.StartOffset,
		),
		libkafka.WithReceiverSecurity(
			confInst.LoginAttempts.KafkaConfig.Username,
			confInst.LoginAttempts.KafkaConfig.Password,
		),
		libkafka.WithReceiverFetchBounds(
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
	connectorInst, err := container.GetInstance[broker.Connector[*amqp.Session]](InstanceAMQPConnector)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}
	receiverConfInst, err := container.GetInstance[*amqp.ReceiverOptions](InstanceLoginAttemptsAMQPReceiverReceiverOpts)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	receiver, err := libamqp.NewReceiver(
		libamqp.WithReceiverConnector(connectorInst),
		libamqp.WithReceiverTargetName(confInst.LoginAttempts.AMQPConfig.TargetName),
		libamqp.WithReceiverLinkCredit(int32(confInst.LoginAttempts.AMQPConfig.PrefetchCredit)),
		libamqp.WithReceiverConnectTimeout(confInst.LoginAttempts.AMQPConfig.ConnectTimeout),
		libamqp.WithReceiverShutdownTimeout(confInst.LoginAttempts.AMQPConfig.ShutdownTimeout),
		libamqp.WithReceiverOpts(receiverConfInst),
		libamqp.WithReceiverLogger(logInst),
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

	connectorInst, err := libamqp.NewConnector(
		libamqp.WithConnectorURL(confInst.AMQPConnector.URL),
		libamqp.WithConnectorConnectTimeout(confInst.AMQPConnector.ConnectTimeout),
		libamqp.WithConnectorShutdownTimeout(confInst.AMQPConnector.ShutdownTimeout),
		libamqp.WithConnectorConnOpts(connOpts),
		libamqp.WithConnectorSessionOpts(sessOpts),
		libamqp.WithConnectorLogger(logInst),
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
