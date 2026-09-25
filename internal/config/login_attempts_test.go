package config

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/stretchr/testify/assert"
)

// 1. Тест успешной валидации (Happy Path) для обоих режимов брокеров
func TestLoginAttemptsConfig_Validate_Success(t *testing.T) {
	baseCfg := func(kind string) *LoginAttemptsConfig {
		return NewLoginAttemptsConfig(
			kind,
			1*time.Second,
			100*time.Millisecond,
			5,
			1000,
			true,
			15*time.Second,
			100,
			5*time.Second,
			10*time.Second,
			config.NewDefaultAMQPReceiverConfig(),
			config.NewDefaultKafkaReceiverConfig(),
		)
	}

	t.Run("kafka_mode_success", func(t *testing.T) {
		cfg := baseCfg("kafka")
		cfg.KafkaConfig.TargetName = "tiny.auth"
		cfg.KafkaConfig.GroupID = "audit-group"
		cfg.KafkaConfig.Partition = -1

		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("amqp_mode_success", func(t *testing.T) {
		cfg := baseCfg("amqp")
		cfg.AMQPConfig.TargetName = "tiny.auth::login.attempts"

		err := cfg.Validate()
		assert.NoError(t, err)
	})
}

// 2. Тест создания дефолтной конфигурации
func TestLoginAttemptsConfig_NewDefault(t *testing.T) {
	cfg := NewDefaultLoginAttemptsConfig()
	assert.NotNil(t, cfg)
}

// 3. Табличный тест на все негативные сценарии валидации воркера (Edge Cases)
func TestLoginAttemptsConfig_Validate_Failures(t *testing.T) {
	// Фабрика базового валидного конфига (по умолчанию берем kafka)
	validKafkaBase := func() *LoginAttemptsConfig {
		c := NewLoginAttemptsConfig(
			"kafka",
			1*time.Second,
			100*time.Millisecond,
			5,
			1000,
			true,
			15*time.Second,
			100,
			5*time.Second,
			10*time.Second,
			config.NewDefaultAMQPReceiverConfig(),
			config.NewDefaultKafkaReceiverConfig(),
		)
		c.KafkaConfig.TargetName = "tiny.auth"
		c.KafkaConfig.GroupID = "audit-group"
		c.KafkaConfig.Partition = -1
		return c
	}

	tests := []struct {
		name           string
		mutate         func(c *LoginAttemptsConfig) // ИСПРАВЛЕНО: Только один чистый мутатор
		expectedPhrase string
	}{
		{
			name: "unknown_receiver_kind",
			mutate: func(c *LoginAttemptsConfig) {
				c.ReceiverKind = "invalid-broker"
			},
			expectedPhrase: "unknown receiver kind",
		},
		{
			name: "invalid_schedule_interval",
			mutate: func(c *LoginAttemptsConfig) {
				c.ScheduleInterval = 0
			},
			expectedPhrase: "ScheduleInterval",
		},
		{
			name: "zero_worker_count",
			mutate: func(c *LoginAttemptsConfig) {
				c.WorkerCount = 0
			},
			expectedPhrase: "WorkerCount",
		},
		{
			name: "zero_data_capacity",
			mutate: func(c *LoginAttemptsConfig) {
				c.DataCapacity = -1
			},
			expectedPhrase: "DataCapacity",
		},
		{
			name: "invalid_batch_size",
			mutate: func(c *LoginAttemptsConfig) {
				c.BatchSize = 0
			},
			expectedPhrase: "BatchSize",
		},
		{
			name: "nil_kafka_config_when_kind_is_kafka",
			mutate: func(c *LoginAttemptsConfig) {
				c.KafkaConfig = nil
			},
			expectedPhrase: "is required",
		},
		{
			name: "deep_validation_failure_inside_kafka_config",
			mutate: func(c *LoginAttemptsConfig) {
				c.KafkaConfig.TargetName = "" // Опустошаем топик внутри дочернего конфига
			},
			expectedPhrase: "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validKafkaBase()
			if tt.mutate != nil {
				tt.mutate(cfg) // ИСПРАВЛЕНО: Прямой вызов без путаницы
			}

			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedPhrase)
		})
	}
}
