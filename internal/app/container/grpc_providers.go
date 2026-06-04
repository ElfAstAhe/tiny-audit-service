package container

import (
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/grpc"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	grpcsvc "github.com/ElfAstAhe/tiny-audit-service/internal/transport/grpc"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"

	libgrpc "google.golang.org/grpc"
)

func (gc *GRPCContainer) serviceRegister(server *libgrpc.Server) error {
	authAuditServiceInst, err := container.GetInstance[*grpcsvc.AuthAuditGRPCService](InstanceAuthAuditGRPCService)
	if err != nil {
		return errs.NewContainerError(gc.GetName(), "service register: retrieve instance failed", err)
	}
	dataAuditServiceInst, err := container.GetInstance[*grpcsvc.DataAuditGRPCService](InstanceDataAuditGRPCService)
	if err != nil {
		return errs.NewContainerError(gc.GetName(), "service register: retrieve instance failed", err)
	}

	pb.RegisterAuthAuditServiceServer(server, authAuditServiceInst)
	pb.RegisterDataAuditServiceServer(server, dataAuditServiceInst)

	return nil
}

func (gc *GRPCContainer) providerAuthAuditGRPCService() (any, error) {
	authFacadeInst, err := container.GetInstance[facade.AuthAuditFacade](InstanceAuthFacade)
	if err != nil {
		return nil, errs.NewContainerError(gc.GetName(), "service register: retrieve instance failed", err)
	}

	return grpcsvc.NewAuthAuditGRPCService(authFacadeInst), nil
}

func (gc *GRPCContainer) providerDataAuditGRPCService() (any, error) {
	dataFacadeInst, err := container.GetInstance[facade.DataAuditFacade](InstanceDataFacade)
	if err != nil {
		return nil, errs.NewContainerError(gc.GetName(), "service register: retrieve instance failed", err)
	}

	return grpcsvc.NewDataAuditGRPCService(dataFacadeInst), nil
}

func (gc *GRPCContainer) providerGRPCRunner() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(gc.GetName(), "provider: retrieve instance failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(gc.GetName(), "provider: retrieve instance failed", err)
	}

	runner, err := grpc.NewRunner(
		grpc.WithName("main-grpc-server"),
		grpc.WithConfig(confInst.GRPC),
		grpc.WithLogger("grpc_server", logInst),
		grpc.WithServiceRegister(gc.serviceRegister),
	)
	if err != nil {
		return nil, errs.NewContainerError(gc.GetName(), fmt.Sprintf("provider: create %s failed", InstanceGRPCRunner), err)
	}

	return runner, nil
}
