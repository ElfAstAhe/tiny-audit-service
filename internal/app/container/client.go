package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

const (
	InstanceAMQPConnector                         string = "amqp-connector"
	InstanceAMQPConnectorConnOpts                 string = "amqp-connector-conn-opts"
	InstanceAMQPConnectorSessOpts                 string = "amqp-connector-sess-opts"
	InstanceLoginAttemptsAMQPReceiver             string = "login-attempts-amqp-receiver"
	InstanceLoginAttemptsAMQPReceiverReceiverOpts string = "login-attempts-amqp-receiver-receiver-opts"
	InstanceLoginAttemptsKafkaReceiver            string = "login-attempts-kafka-receiver"
)

type ClientContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*ClientContainer)(nil)
var _ container.LazyContainer = (*ClientContainer)(nil)

func NewClientContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *ClientContainer {
	return &ClientContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(ClientContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

func (cc *ClientContainer) Init(ctx context.Context) error {
	err := errors.Join(
		cc.RegisterProvider(InstanceLoginAttemptsKafkaReceiver, cc.providerLoginAttemptsKafkaReceiver),
		cc.RegisterProvider(InstanceLoginAttemptsAMQPReceiver, cc.providerLoginAttemptsAMQPReceiver),
		cc.RegisterProvider(InstanceLoginAttemptsAMQPReceiverReceiverOpts, cc.providerLoginAttemptsAMQPReceiverReceiverOpts),
		cc.RegisterProvider(InstanceAMQPConnector, cc.providerAMQPConnector),
		cc.RegisterProvider(InstanceAMQPConnectorConnOpts, cc.providerAMQPConnectorConnOpts),
		cc.RegisterProvider(InstanceAMQPConnectorSessOpts, cc.providerAMQPConnectorSessOpts),
	)
	if err != nil {
		return errs.NewContainerError(cc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
