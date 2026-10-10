package config

import (
	"time"
)

// Блок констант со значениями по умолчанию для REST-клиента аудита.
const (
	// DefaultRestBaseURL указывает базовый URL-адрес удаленного REST API сервиса аудита.
	DefaultRestBaseURL string = "http://localhost:8080/"

	// DefaultRestReadTimeout определяет максимальное время ожидания ответа на HTTP-запрос REST-клиента.
	DefaultRestReadTimeout time.Duration = 5 * time.Second
)

// Блок констант со значениями по умолчанию для gRPC-клиента аудита.
const (
	// DefaultGRPCTarget задает целевой сетевой адрес gRPC-сервера аудита для локальной среды.
	DefaultGRPCTarget string = "localhost:51052"

	// DefaultGRPCConnectionTimeout определяет таймаут на установку сетевого соединения gRPC-клиентом.
	DefaultGRPCConnectionTimeout time.Duration = 5 * time.Second
)
