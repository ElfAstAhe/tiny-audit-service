package container

import (
	"fmt"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp/azure"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

func (cc *ClientContainer) providerLoginAttemptsReceiver() (any, error) {
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
	receiverConfInst, err := container.GetInstance[*amqp.ReceiverOptions](InstanceLoginAttemptsReceiverReceiverOpts)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	receiver, err := azure.NewReceiver(
		azure.WithReceiverConnector(connectorInst),
		azure.WithReceiverTargetName(confInst.LoginAttempts.ReceiverConf.TargetName),
		azure.WithReceiverLinkCredit(int32(confInst.LoginAttempts.ReceiverConf.PrefetchCredit)),
		azure.WithReceiverConnectTimeout(confInst.LoginAttempts.ReceiverConf.ConnectTimeout),
		azure.WithReceiverShutdownTimeout(confInst.LoginAttempts.ReceiverConf.ShutdownTimeout),
		azure.WithReceiverOpts(receiverConfInst),
		azure.WithReceiverLogger(logInst),
	)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceLoginAttemptsReceiver), err)
	}

	return receiver, nil
}

func (cc *ClientContainer) providerLoginAttemptsReceiverReceiverOpts() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(cc.GetName(), "provider: retrieve instance failed", err)
	}

	return &amqp.ReceiverOptions{
		Name:   confInst.App.NodeName,
		Credit: int32(confInst.LoginAttempts.ReceiverConf.PrefetchCredit),
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
