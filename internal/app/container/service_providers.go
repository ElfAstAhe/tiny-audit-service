package container

import (
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/hellofresh/health-go/v5"
)

func (sc *ServiceContainer) providerHealthStatus() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(sc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := health.New(health.WithComponent(health.Component{
		Name:    confInst.App.NodeName,
		Version: config.AppVersion,
	}))
	if err != nil {
		return nil, errs.NewContainerError(sc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceHealthStatus), err)
	}

	return res, nil
}
