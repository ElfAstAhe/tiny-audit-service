package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов коннекторов, ресиверов и их опций для AMQP и Kafka в DI-контейнере.
const (
	// InstanceAMQPConnector определяет имя инстанса общего AMQP-коннектора.
	InstanceAMQPConnector string = "amqp-connector"

	// InstanceAMQPConnectorConnOpts определяет имя инстанса параметров подключения (connection options) для AMQP-коннектора.
	InstanceAMQPConnectorConnOpts string = "amqp-connector-conn-opts"

	// InstanceAMQPConnectorSessOpts определяет имя инстанса параметров сессии (session options) для AMQP-коннектора.
	InstanceAMQPConnectorSessOpts string = "amqp-connector-sess-opts"

	// InstanceLoginAttemptsAMQPReceiver определяет имя инстанса AMQP-ресивера для получения сообщений о попытках входа.
	InstanceLoginAttemptsAMQPReceiver string = "login-attempts-amqp-receiver"

	// InstanceLoginAttemptsAMQPReceiverReceiverOpts определяет имя инстанса специфичных опций AMQP-ресивера попыток входа.
	InstanceLoginAttemptsAMQPReceiverReceiverOpts string = "login-attempts-amqp-receiver-receiver-opts"

	// InstanceLoginAttemptsKafkaReceiver определяет имя инстанса Kafka-ресивера для получения сообщений о попытках входа.
	InstanceLoginAttemptsKafkaReceiver string = "login-attempts-kafka-receiver"
)

// ClientContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом и опциями внешних клиентов, коннекторов очередей сообщений и ресиверов (AMQP/Kafka).
type ClientContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*ClientContainer)(nil)
var _ container.LazyContainer = (*ClientContainer)(nil)

// NewClientContainer создает и инициализирует новый экземпляр ленивого контейнера клиентов ClientContainer.
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

// Init выполняет первичную инициализацию контейнера ClientContainer, регистрируя провайдеры для всех брокер-клиентов и их опций.
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
