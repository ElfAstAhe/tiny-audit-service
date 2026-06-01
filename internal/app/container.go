package app

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
)

const ContainerName string = "app"

const (
	InstanceApplication       string = "application"
	InstanceApplicationReady  string = "application-ready"
	InstanceApplicationHealth string = "application-health"
	InstanceConfig            string = "config"
	InstanceLogger            string = "logger"
)

type Container struct {
	*container.BaseContainer
}

var _ container.Container = (*Container)(nil)

func NewContainer(
	orchestrator container.Orchestrator,
) *Container {
	return &Container{
		BaseContainer: container.NewBaseContainer(ContainerName, orchestrator),
	}
}

func (c *Container) Init(ctx context.Context) error {
	return nil
}
