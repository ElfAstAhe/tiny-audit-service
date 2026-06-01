package postgres

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
)

const ContainerName string = "postgres"

type Container struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*Container)(nil)
var _ container.LazyContainer = (*Container)(nil)

func NewContainer(orchestrator container.Orchestrator) *Container {
	return &Container{
		BaseLazyContainer: container.NewBaseLazyContainer(ContainerName, orchestrator),
	}
}

func (c *Container) Init(ctx context.Context) error {
	// ToDo: implement

	return nil
}
