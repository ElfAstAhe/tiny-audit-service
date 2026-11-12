package utils

import (
	"encoding/json"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/pelletier/go-toml/v2"
	"go.yaml.in/yaml/v3"
)

// BuildJSONPayload сериализует переданные данные в формат JSON с использованием encoding/json/v2.
func BuildJSONPayload(data any) ([]byte, error) {
	return buildPayload(data, func(v any) ([]byte, error) {
		return json.Marshal(v)
	})
}

// BuildYAMLPayload сериализует переданные данные в формат YAML (v3).
func BuildYAMLPayload(data any) ([]byte, error) {
	return buildPayload(data, yaml.Marshal)
}

// BuildTOMLPayload сериализует переданную структуру или карту в формат TOML (v2).
func BuildTOMLPayload(data any) ([]byte, error) {
	return buildPayload(data, toml.Marshal)
}

// buildPayload проверяет входные данные на nil и выполняет их сериализацию с помощью переданной функции marshaller.
func buildPayload(data any, marshaller func(any) ([]byte, error)) ([]byte, error) {
	if utils.IsNil(data) {
		return nil, errs.NewInvalidArgumentError("data", "is required")
	}

	res, err := marshaller(data)
	if err != nil {
		return nil, errs.NewCommonError("payload encoding failed", err)
	}

	return res, nil
}
