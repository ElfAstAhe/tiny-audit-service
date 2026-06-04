package container

import (
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (tc *ToolsContainer) providerHashCipher() (any, error) {
	return utils.NewSHA256Hash(), nil
}

func (tc *ToolsContainer) providerDataCipher() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}
	hashCipherInst, err := container.GetInstance[utils.Cipher](InstanceHashCipher)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}

	var dataCipherKey []byte
	dataCipherKey, err = hashCipherInst.Encrypt([]byte(confInst.App.CipherKey))
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: build correct data cipher key", err)
	}

	res, err := utils.NewAesGcmCipher(dataCipherKey)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceDataCipher), err)
	}

	return res, nil
}

func (tc *ToolsContainer) providerDataCipherHelper() (any, error) {
	dataCipherInst, err := container.GetInstance[utils.Cipher](InstanceDataCipher)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}

	return helper.NewCipherHelper(dataCipherInst), nil
}

func (tc *ToolsContainer) providerKeysHelper() (any, error) {
	return helper.NewRSAKeysHelper(helper.RSAKey2048), nil
}

func (tc *ToolsContainer) providerJWTHelper() (any, error) {
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}
	return helper.NewJWTHelper(
		confInst.App.TokenIssuer,
		jwt.GetSigningMethod(confInst.Auth.JWTSigningMethod),
		confInst.Auth.JWTSecret,
		confInst.Auth.AccessTokenTTL,
		func() string {
			template := "tas-%v"
			rnd, err := uuid.NewRandom()
			if err != nil {
				return fmt.Sprintf(template, time.Now().Nanosecond())
			}

			return fmt.Sprintf(template, rnd.String())
		},
	), nil
}

func (tc *ToolsContainer) providerJWTHTTPHelper() (any, error) {
	jwtHelperInst, err := container.GetInstance[*helper.JWTHelper](InstanceJWTHelper)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}

	return helper.NewJWTHTTPHelper(jwtHelperInst), nil
}

func (tc *ToolsContainer) providerJWTGRPCHelper() (any, error) {
	jwtHelperInst, err := container.GetInstance[*helper.JWTHelper](InstanceJWTHelper)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}

	return helper.NewJWTGRPCHelper(jwtHelperInst), nil
}

func (tc *ToolsContainer) providerAuthHelper() (any, error) {
	jwtHelperInst, err := container.GetInstance[*helper.JWTHelper](InstanceJWTHelper)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}
	jwtHTTPHelperInst, err := container.GetInstance[*helper.JWTHTTPHelper](InstanceJWTHTTPHelper)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}
	jwtGRPCHelperInst, err := container.GetInstance[*helper.JWTGRPCHelper](InstanceJWTGRPCHelper)
	if err != nil {
		return nil, errs.NewContainerError(tc.GetName(), "provider: retrieve instance failed", err)
	}
	return auth.NewHelper(
		auth.DefaultHeaderName,
		auth.DefaultCookieName,
		auth.DefaultMetadataName,
		jwtHelperInst,
		jwtHTTPHelperInst,
		jwtGRPCHelperInst,
	), nil
}
