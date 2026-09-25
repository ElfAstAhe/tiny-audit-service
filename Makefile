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
KAFKA_DIR   = /opt/kafka_2.13-4.3.1
ARTEMIS_RUN = /var/lib/artemis-test-cluster/bin/artemis

.PHONY: gen-proto gen-swagger gen-http-client gen-mocks build run run-amqp run-kafka test static-check lint clean update-deps artemis-local-start artemis-local-stop kafka-local-start kafka-local-stop brokers-all-start kafka-docker-start kafka-docker-stop kafka-docker-logs

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

# запуск проекта с выводом информации о параметрах
run-help: build ## Собрать проект и запустить бинарник с информацией о параметрах
	./bin/$(SERVER_BINARY_NAME) --help

# Запуск проекта (сначала соберет, потом запустит)
run: build ## Собрать проект и запустить бинарник с локальными флагами (AMQP, БД, логи)
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
		--login-attempts-receiver-kind "kafka" \
		--login-attempts-worker-start-interval "3s" \
		--login-attempts-worker-schedule-interval "10s" \
		--login-attempts-worker-worker-count "3" \
		--login-attempts-worker-data-capacity "128" \
		--login-attempts-worker-complete-processing \
		--login-attempts-worker-shutdown-timeout "3s" \
		--login-attempts-worker-batch-size "25" \
		--login-attempts-worker-batch-read-timeout "10s" \
		--login-attempts-worker-ack-timeout "3s" \
		--login-attempts-receiver-amqp-target-name "tiny.auth::login.attempts" \
		--login-attempts-receiver-amqp-connect-timeout "2s" \
		--login-attempts-receiver-amqp-shutdown-timeout "3s" \
		--login-attempts-receiver-amqp-prefetch-credit "50" \
		--login-attempts-receiver-kafka-brokers "localhost:9092" \
		--login-attempts-receiver-kafka-target-name "tiny.auth.login.attempts" \
		--login-attempts-receiver-kafka-group-id "tiny-audit" \
		--login-attempts-receiver-kafka-partition "-1" \
		--login-attempts-receiver-kafka-connect-timeout "10s" \
		--login-attempts-receiver-kafka-shutdown-timeout "3s" \
		--login-attempts-receiver-kafka-min-bytes "1024" \
		--login-attempts-receiver-kafka-max-bytes "10000000" \
		--login-attempts-receiver-kafka-max-wait "500ms" \
		--login-attempts-receiver-kafka-heartbeat-interval "3s" \
		--login-attempts-receiver-kafka-session-timeout "30s" \
		--login-attempts-receiver-kafka-rebalance-timeout "45s" \
		--login-attempts-receiver-kafka-read-timeout "10s" \
		--login-attempts-receiver-kafka-max-attempts "2" \
		--login-attempts-receiver-kafka-queue-capacity "100" \
		--login-attempts-receiver-kafka-start-offset "last"

# Запуск проекта (сначала соберет, потом запустит)
run-amqp: build ## Собрать проект и запустить бинарник с локальными флагами (AMQP, БД, логи)
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
		--login-attempts-receiver-kind "amqp" \
		--login-attempts-receiver-amqp-target-name "tiny.auth::login.attempts" \
		--login-attempts-receiver-amqp-connect-timeout "2s" \
		--login-attempts-receiver-amqp-shutdown-timeout "3s" \
		--login-attempts-receiver-amqp-prefetch-credit "50" \
		--login-attempts-worker-start-interval "3s" \
		--login-attempts-worker-schedule-interval "10s" \
		--login-attempts-worker-worker-count "3" \
		--login-attempts-worker-data-capacity "128" \
		--login-attempts-worker-complete-processing \
		--login-attempts-worker-shutdown-timeout "3s" \
		--login-attempts-worker-batch-size "25" \
		--login-attempts-worker-batch-read-timeout "3s" \
		--login-attempts-worker-ack-timeout "3s"

# Запуск проекта (сначала соберет, потом запустит)
run-kafka: build ## Собрать проект и запустить бинарник с локальными флагами (Kafka, БД, логи)
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
		--data-tc-start-interval "5s" \
		--data-tc-schedule-interval "66s" \
		--data-tc-worker-count "2" \
		--data-tc-data-capacity "128" \
		--data-tc-shutdown-timeout "15s" \
		--data-tc-tail-interval "8760h" \
		--data-tc-tail-cut \
		--login-attempts-receiver-kind "kafka" \
		--login-attempts-worker-start-interval "3s" \
		--login-attempts-worker-schedule-interval "10s" \
		--login-attempts-worker-worker-count "3" \
		--login-attempts-worker-data-capacity "128" \
		--login-attempts-worker-complete-processing \
		--login-attempts-worker-shutdown-timeout "3s" \
		--login-attempts-worker-batch-size "25" \
		--login-attempts-worker-batch-read-timeout "3s" \
		--login-attempts-worker-ack-timeout "3s" \
		--login-attempts-receiver-kafka-brokers "localhost:9092" \
		--login-attempts-receiver-kafka-target-name "tiny.auth.login.attempts" \
		--login-attempts-receiver-kafka-group-id "tiny-audit" \
		--login-attempts-receiver-kafka-partition "-1" \
		--login-attempts-receiver-kafka-connect-timeout "10s" \
		--login-attempts-receiver-kafka-shutdown-timeout "3s" \
		--login-attempts-receiver-kafka-min-bytes "1024" \
		--login-attempts-receiver-kafka-max-bytes "10000000" \
		--login-attempts-receiver-kafka-max-wait "500ms" \
		--login-attempts-receiver-kafka-heartbeat-interval "3s" \
		--login-attempts-receiver-kafka-session-timeout "30s" \
		--login-attempts-receiver-kafka-rebalance-timeout "45s" \
		--login-attempts-receiver-kafka-read-timeout "10s" \
		--login-attempts-receiver-kafka-max-attempts "2" \
		--login-attempts-receiver-kafka-queue-capacity "100" \
		--login-attempts-receiver-kafka-start-offset "last"

# Запуск тестов
test: gen-proto gen-mocks ## Запустить модульные и интеграционные тесты проекта
	go test -v $$(go list ./... | grep -vE "mocks")

# Запуск бенчмарков (сюда добавляем все вызовы) или разные параметры под один пакет
bench: gen-proto gen-mocks ## Запустить утилиты с замером памяти
#	go test -bench=BenchmarkManager_FullCycle -benchmem ./pkg/infra/cache/test/...
	go test -bench=. -benchmem $$(go list ./... | grep -vE "mocks")

# Запуск static check
static-check: ## Запустить статический анализ кода (пропуская автогенерируемый pkg/api)
	staticcheck $$(go list ./... | grep -vE "pkg/api|mocks")

# Запуск линтера
lint: ## Запустить линтер revive (пропуская автогенерируемый код)
	revive -exclude "_test\.go$$" $$(go list ./... | grep -vE "pkg/api|mocks")

# Очистка бинарников
clean: ## Очистить скомпилированные файлы из папки ./bin
	rm -rf ./bin/*

# обновление зависимостей
update-deps: ## Принудительно обновить и скачать все Go-зависимости проекта
	go get -u -x all

artemis-local-start: ## start artemis local (ubuntu, in separate terminal)
	gnome-terminal -- bash -c "sudo $(ARTEMIS_RUN) run; exec bash"

artemis-local-stop: ## stop artemis local (not implemented)
	echo "not implemented :-)"

kafka-local-start: ## start kafka local (ubuntu, in separate terminal)
	gnome-terminal -- bash -c "$(KAFKA_DIR)/bin/kafka-server-start.sh $(KAFKA_DIR)/config/server.properties; exec bash"

kafka-local-stop: ## stop kafka local (not implemented)
	echo "not implemented :-)"

brokers-all-start: kafka-local-start artemis-local-start ## start both brokers simultaneously in separate windows

# start kafka (docker compose)
kafka-docker-start: ## start kafka docker container (docker compose)
	docker compose -f kafka-docker-compose.yml up -d

# stop kafka (docker compose)
kafka-docker-stop: ## stop kafka docker container (docker compose)
	docker compose -f kafka-docker-compose.yml down

# watch kafka logs (docker compose)
kafka-docker-logs: ## watch kafka logs (docker compose)
	docker compose -f kafka-docker-compose.yml logs -f

#
