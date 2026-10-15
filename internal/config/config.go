package config

import (
	"log"
	"strings"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/spf13/viper"
)

const (
	// ApplicationName определяет константное системное имя текущего микросервиса.
	ApplicationName = "tiny-audit-service"
)

// Config является центральным агрегатором (композитным паспортом) всей конфигурации микросервиса.
//
// Сводит воедино параметры ядра, настройки среза хвостов очередей (TailCutter),
// конфигурации сетевых протоколов (HTTP/gRPC), базы данных, телеметрии брокера AMQP и лимитов попыток входа.
type Config struct {
	App           *AppConfig                  `mapstructure:"app" json:"app,omitempty" yaml:"app,omitempty"`
	AuthTC        *TailCutterConfig           `mapstructure:"auth_tc" json:"auth_tc,omitempty" yaml:"auth_tc,omitempty"`
	DataTC        *TailCutterConfig           `mapstructure:"data_tc" json:"data_tc,omitempty" yaml:"data_tc,omitempty"`
	Auth          *config.AuthConfig          `mapstructure:"auth" json:"auth,omitempty" yaml:"auth,omitempty"`
	HTTP          *config.HTTPConfig          `mapstructure:"http" json:"http,omitempty" yaml:"http,omitempty"`
	GRPC          *config.GRPCConfig          `mapstructure:"grpc" json:"grpc,omitempty" yaml:"grpc,omitempty"`
	Log           *config.LogConfig           `mapstructure:"log" json:"log,omitempty" yaml:"log,omitempty"`
	DB            *config.DBConfig            `mapstructure:"db" json:"db,omitempty" yaml:"db,omitempty"`
	Telemetry     *config.TelemetryConfig     `mapstructure:"telemetry" json:"telemetry,omitempty" yaml:"telemetry,omitempty"`
	AMQPConnector *config.AMQPConnectorConfig `mapstructure:"amqp_connector" json:"amqp_connector,omitempty" yaml:"amqp_connector,omitempty"`
	LoginAttempts *LoginAttemptsConfig        `mapstructure:"login_attempts_receiver" json:"login_attempts,omitempty" yaml:"login_attempts,omitempty"`
}

// Переменные для динамической инжекции метаданных сборки компилятором через linker flags (-ldflags)
var (
	AppName      string
	AppVersion   string
	AppBuildTime string
)

// NewConfig — параметризованный конструктор для явного формирования агрегированного объекта конфигурации.
func NewConfig(
	app *AppConfig,
	authTC *TailCutterConfig,
	dataTC *TailCutterConfig,
	auth *config.AuthConfig,
	HTTP *config.HTTPConfig,
	GRPC *config.GRPCConfig,
	log *config.LogConfig,
	db *config.DBConfig,
	telemetry *config.TelemetryConfig,
	amqpConnector *config.AMQPConnectorConfig,
	loginAttempts *LoginAttemptsConfig,
) *Config {
	return &Config{
		App:           app,
		AuthTC:        authTC,
		DataTC:        dataTC,
		Auth:          auth,
		HTTP:          HTTP,
		GRPC:          GRPC,
		Log:           log,
		DB:            db,
		Telemetry:     telemetry,
		AMQPConnector: amqpConnector,
		LoginAttempts: loginAttempts,
	}
}

// NewDefaultConfig инициализирует объект конфигурации со стандартными значениями таймаутов и лимитов по умолчанию.
func NewDefaultConfig() *Config {
	return NewConfig(
		NewDefaultAppConfig(),
		NewDefaultTailCutterConfig(),
		NewDefaultTailCutterConfig(),
		config.NewDefaultAuthConfig(),
		config.NewDefaultHTTPConfig(),
		config.NewDefaultGRPCConfig(),
		config.NewDefaultLogConfig(),
		config.NewDefaultDBConfig(),
		config.NewDefaultTelemetryConfig(),
		config.NewDefaultAMQPConnectorConfig(),
		NewDefaultLoginAttemptsConfig(),
	)
}

// NewEmptyConfig выделяет память под пустые указатели внутренних структур для последующего заполнения парсером Viper.
func NewEmptyConfig() *Config {
	return &Config{
		App:           &AppConfig{},
		AuthTC:        &TailCutterConfig{},
		DataTC:        &TailCutterConfig{},
		Auth:          &config.AuthConfig{},
		HTTP:          &config.HTTPConfig{},
		GRPC:          &config.GRPCConfig{},
		Log:           &config.LogConfig{},
		DB:            &config.DBConfig{},
		Telemetry:     &config.TelemetryConfig{},
		AMQPConnector: &config.AMQPConnectorConfig{},
		LoginAttempts: &LoginAttemptsConfig{},
	}
}

// Validate итерируется по вееру зарегистрированных вложенных конфигураций и каскадно запускает их проверку целостности.
func (c *Config) Validate() error {
	validators := []interface {
		Validate() error
	}{
		c.App,
		c.AuthTC,
		c.DataTC,
		c.Auth,
		c.HTTP,
		c.GRPC,
		c.Log,
		c.DB,
		c.Telemetry,
		c.AMQPConnector,
		c.LoginAttempts,
	}

	for _, validator := range validators {
		if err := validator.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Load поэтапно собирает итоговую конфигурацию в соответствии с приоритетом: Flags -> ENV -> YAML -> Defaults.
func Load() (*Config, error) {
	v := viper.New()

	// 1. Установка значений по умолчанию (Defaults)
	applyDefaults(v)

	// 2. Настройка Флагов (pflag для длинных имен --port)
	pFlagSet, err := initFLags()
	if err != nil {
		return nil, errs.NewConfigError("failed init flags", err)
	}

	// 3. Привязываем все флаги к Viper
	if err := bindFlags(pFlagSet, v); err != nil {
		return nil, err
	}

	// 4. Настройка Переменных окружения (ENV)
	// Используем твой механизм AutomaticEnv
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// 5. Поддержка ENV для пути к конфигу
	err = v.BindEnv(FlagConfig, EnvConfig)
	if err != nil {
		return nil, errs.NewConfigError("failed to bind env", err)
	}

	// 6. Чтение файла конфигурации
	cfgFile := v.GetString(FlagConfig)
	v.SetConfigFile(cfgFile)

	if err := v.ReadInConfig(); err != nil {
		// Если файла нет — это предупреждение, но не фатальная ошибка (могут быть ENV)
		log.Printf("WARN: config file not found [%s]\n", cfgFile)
	}

	// 7. Маппинг в структуру
	var cfg = NewEmptyConfig()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, errs.NewConfigError("failed to unmarshal config struct", err)
	}

	// 8. Итоговая валидация по всем слоям
	if err := cfg.Validate(); err != nil {
		return nil, errs.NewConfigError("config validation failed", err)
	}

	return cfg, nil
}
