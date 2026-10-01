package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/infra/metrics"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// InfraContainer structures a lazy-loaded lifecycle dependency injection container managing cross-cutting telemetry metrics and broker notification routers.
type InfraContainer struct {
	*container.BaseLazyContainer // Generic framework-level baseline container orchestration handle
}

// Compile-time interface compliance verifications
var _ container.Container = (*InfraContainer)(nil)
var _ container.LazyContainer = (*InfraContainer)(nil)

// NewInfraContainer acts as a factory constructor deploying structural configuration, logger, and orchestrator boundaries.
func NewInfraContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *InfraContainer {
	return &InfraContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(InfraContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init registers event-driven notification providers, initializes cross-cutting Prometheus metrics, and dynamically binds active infrastructure observers.
//
//goland:noinspection GoUnusedParameter
func (ic *InfraContainer) Init(ctx context.Context) error {
	err := errors.Join(
	//ic.RegisterProvider(InstanceLoginAttemptsPublisher, ic.providerLoginAttemptsEventDispatcher),
	//ic.RegisterProvider(InstanceLoginAttemptsAMQPSubscriber, ic.providerLoginAttemptsAMQPObserver),
	//ic.RegisterProvider(InstanceLoginAttemptsKafkaSubscriber, ic.providerLoginAttemptsKafkaObserver),
	)
	if err != nil {
		return errs.NewContainerError(ic.GetName(), "container init: register providers failed", err)
	}

	// setup metrics
	metrics.InitHTTPMetrics()
	metrics.InitRepositoryMetrics()
	metrics.InitBrokerSenderMetrics()

	return nil
}
