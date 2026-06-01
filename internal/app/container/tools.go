package container

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
)

type ToolsContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*ToolsContainer)(nil)
var _ container.LazyContainer = (*ToolsContainer)(nil)

func NewToolsContainer(orchestrator container.Orchestrator) *ToolsContainer {
	return &ToolsContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(ToolsContainerName, orchestrator),
	}
}

func (tc *ToolsContainer) Init(ctx context.Context) error {
	// ToDo: implement

	return nil
}
