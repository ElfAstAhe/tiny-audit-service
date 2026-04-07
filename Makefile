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

.PHONY: build run test clean

# Генерация gRPC кода
gen-proto:
	mkdir -p $(PROTO_OUT)
	protoc \
        -I $(PROTO_ROOT) \
		--proto_path=$(PROTO_PATH) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		--go_opt=default_api_level=API_OPAQUE \
		$(PROTO_PATH)/*.proto

# Генерация swagger
gen-swagger:
	swag init \
		-g $(SERVER_BUILD_DIR)/main.go \
		--parseDependency \
		--parseInternal \
		--exclude ./pkg/api \
		-o docs \
		--parseDepth 3

# Генерация http client
gen-http-client:
#	oapi-codegen -package client -generate client docs/swagger.json > pkg/client/rest/api_client.gen.go
	mkdir -p $(OPEN_API_OUT)
	swagger generate client -f ./docs/swagger.json -A tiny-audit-service -t $(OPEN_API_OUT)

# Генерирует моки для всех интерфейсов в указанной папке
gen-mocks:
	mockery

# Сборка проекта с прокидыванием переменных
build: gen-proto gen-swagger gen-http-client gen-mocks
	go build -ldflags "-X '$(MODULE_NAME)/internal/config.AppVersion=$(VERSION)' \
	-X '$(MODULE_NAME)/internal/config.AppBuildTime=$(BUILD_TIME)'" \
	-o ./bin/$(SERVER_BINARY_NAME) $(SERVER_BUILD_DIR)/main.go

# Запуск проекта (сборка, затем запуск)
run: build
	./bin/$(SERVER_BINARY_NAME) --http-address "localhost:8081" --grpc-address "localhost:51052" --log-level "debug" --db-driver "postgres" --db-dsn "postgres://svc_audit:password@localhost:5432/test?sslmode=disable&search_path=audit_db" --auth-jwt-secret "jwt-key" --app-cipher-key "12345" --app-max-list-limit 500 --app-accept-token-issuers "tiny-auth-service,test-issuer" --app-auth-tail-job-repeat-duration "90s" --app-auth-tail-cut --app-auth-tail-duration "48h" --app-data-tail-job-repeat-duration "85s" --app-data-tail-cut --app-data-tail-duration "48h"

# Запуск тестов
test:
	go test -v ./...

# Очистка бинарников
clean:
	rm -rf ./bin/*
