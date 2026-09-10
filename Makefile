# Переменные для сборки
PROTO_ROOT=api/proto
PROTO_PATH=api/proto/tiny-audit-service/v1
PROTO_OUT=pkg/api/grpc/
OPEN_API_OUT=pkg/api/http/audit/v1
MODULE_NAME=github.com/ElfAstAhe/tiny-audit-service
SERVER_BINARY_NAME=tiny-audit-service
SERVER_BUILD_DIR=./cmd/tiny-audit-service
VERSION=1.0.0
BUILD_TIME=$(shell date +'%Y/%m/%d_%H:%M:%S')
STAGE=DEV

.PHONY: gen-proto gen-swagger gen-http-client gen-mocks build run test static-check lint clean update-deps

help:
	@echo "Доступные команды для сборки и тестирования:"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo "------------------------------------------------------------------------"

# Генерация gRPC кода
gen-proto: ## Сгенерировать gRPC код (Go & gRPC) из Protobuf файлов
	mkdir -p $(PROTO_OUT)
	protoc \
        -I $(PROTO_ROOT) \
		--proto_path=$(PROTO_PATH) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		--go_opt=default_api_level=API_OPAQUE \
		$(PROTO_PATH)/*.proto

# Генерация swagger
gen-swagger: ## Сгенерировать Swagger-документацию (swag init)
	swag init \
		-g $(SERVER_BUILD_DIR)/main.go \
		--parseDependency \
		--parseInternal \
		--exclude ./pkg/api \
		-o docs \
		--parseDepth 3

# Генерация http client
gen-http-client: ## Сгенерировать HTTP-клиент на основе swagger.json
#	oapi-codegen -package client -generate client docs/swagger.json > pkg/client/rest/api_client.gen.go
	mkdir -p $(OPEN_API_OUT)
	swagger generate client -f ./docs/swagger.json -A tiny-audit-service -t $(OPEN_API_OUT)

# Генерирует моки для интерфейсов в указанной папке, см. {project_root}/.mockery.yml конфиг
gen-mocks: ## Сгенерировать моки для интерфейсов (mockery)
	mockery

# Сборка проекта с прокидыванием переменных
build: gen-proto gen-swagger gen-http-client gen-mocks ## Полная сборка: генерация всего кода + компиляция бинарника
	go build -ldflags "-X '$(MODULE_NAME)/internal/config.AppVersion=$(VERSION)' \
	-X '$(MODULE_NAME)/internal/config.AppBuildTime=$(BUILD_TIME)'" \
	-o ./bin/$(SERVER_BINARY_NAME) $(SERVER_BUILD_DIR)/main.go

# Запуск проекта (сначала соберет, потом запустит)
run: build ## Собрать проект и запустить бинарник с локальными флагами (БД, логи)
	./bin/$(SERVER_BINARY_NAME) \
		--http-address "localhost:8081" \
		--grpc-address "localhost:51052" \
		--log-level "debug" \
		--db-driver "postgres" \
		--db-dsn "postgres://svc_audit:password@localhost:5432/test?sslmode=disable&search_path=audit_db" \
		--auth-jwt-secret "jwt-key" \
		--app-cipher-key "12345" \
		--app-max-list-limit 500 \
		--app-accept-token-issuers "tiny-auth-service,test-issuer" \
		--auth-tc-start-interval "4s" \
		--auth-tc-schedule-interval "77s" \
		--auth-tc-worker-count "2" \
		--auth-tc-data-capacity "128" \
		--auth-tc-shutdown-timeout "15s" \
		--auth-tc-tail-interval "4368h" \
		--auth-tc-tail-cut \
		--data-tc-start-interval "5s" \
		--data-tc-schedule-interval "66s" \
		--data-tc-worker-count "2" \
		--data-tc-data-capacity "128" \
		--data-tc-shutdown-timeout "15s" \
		--data-tc-tail-interval "8760h" \
		--data-tc-tail-cut \
		--amqp-connector-url "amqp://localhost:5672" \
		--amqp-connector-username "svc-audit" \
		--amqp-connector-password "test" \
		--amqp-connector-connect-timeout "2s" \
		--amqp-connector-write-timeout "2s" \
		--amqp-connector-idle-timeout "30s" \
		--amqp-connector-shutdown-timeout "3s" \
		--login-attempts-receiver-target-name "tiny.auth::login.attempts" \
		--login-attempts-receiver-connect-timeout "2s" \
		--login-attempts-receiver-shutdown-timeout "3s" \
		--login-attempts-receiver-prefetch-credit "50" \
		--login-attempts-start-interval "3s" \
		--login-attempts-schedule-interval "10s" \
		--login-attempts-worker-count "3" \
		--login-attempts-data-capacity "128" \
		--login-attempts-complete-processing \
		--login-attempts-shutdown-timeout "3s" \
		--login-attempts-batch-size "25" \
		--login-attempts-batch-read-timeout "3s" \
		--login-attempts-ack-timeout "3s"

# Запуск тестов
test: gen-proto gen-mocks ## Запустить модульные и интеграционные тесты проекта
	go test -v ./...

# Запуск бенчмарков (сюда добавляем все вызовы) или разные параметры под один пакет
bench: gen-proto gen-mocks ## Запустить утилиты с замером памяти
#	go test -bench=BenchmarkManager_FullCycle -benchmem ./pkg/infra/cache/test/...
	go test -bench=. -benchmem ./...

# Запуск static check
static-check: ## Запустить статический анализ кода (пропуская автогенерируемый pkg/api)
	staticcheck $$(go list ./... | grep -vE "pkg/api")

# Запуск линтера
lint: ## Запустить линтер
	revive ./...

# Очистка бинарников
clean: ## Очистить скомпилированные файлы из папки ./bin
	rm -rf ./bin/*

# обновление зависимостей
update-deps: ## Принудительно обновить и скачать все Go-зависимости проекта
	go get -u -x all

#
