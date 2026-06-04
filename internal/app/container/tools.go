package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

const (
	InstanceHashCipher       string = "sha256-cipher"
	InstanceDataCipher       string = "aes-gcm-cipher"
	InstanceDataCipherHelper string = "data-cipher-helper"
	InstanceKeysHelper       string = "rsa-2048-keys-helper"
	InstanceJWTHelper        string = "jwt-helper"
	InstanceJWTHTTPHelper    string = "jwt-http-helper"
	InstanceJWTGRPCHelper    string = "jwt-grpc-helper"
	InstanceAuthHelper       string = "auth-helper"
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
	err := errors.Join(
		tc.RegisterProvider(InstanceHashCipher, tc.providerHashCipher),
		tc.RegisterProvider(InstanceDataCipher, tc.providerDataCipher),
		tc.RegisterProvider(InstanceDataCipherHelper, tc.providerDataCipherHelper),
		tc.RegisterProvider(InstanceKeysHelper, tc.providerKeysHelper),
		tc.RegisterProvider(InstanceJWTHelper, tc.providerJWTHelper),
		tc.RegisterProvider(InstanceJWTHTTPHelper, tc.providerJWTHTTPHelper),
		tc.RegisterProvider(InstanceJWTGRPCHelper, tc.providerJWTGRPCHelper),
		tc.RegisterProvider(InstanceAuthHelper, tc.providerAuthHelper),
	)
	if err != nil {
		return errs.NewContainerError(tc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
