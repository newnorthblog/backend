ifeq ($(DATABASE_HOST),)
include .env
endif

LOCAL_BIN:=$(CURDIR)/bin
BUILD_DIR:=$(CURDIR)/cmd

# install goose migrator
install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

# linter
install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml

# compose deps
compose:
	@echo 'compose deps'
	docker compose -f docker-compose.yaml up -d

# down deps
compose-down:
	@echo 'compose deps'
	docker compose -f docker-compose.yaml down 

# build binary
build: deps build-binary

build-binary:
	@echo 'build backend binary'
	go build -tags netgo -ldflags '-s -w' -o $(LOCAL_BIN) $(BUILD_DIR)

deps:
	@echo 'install dependencies'
	go mod tidy -v

# run app
run: deps run-app

run-app:
	@echo "\033[32mrun backend\033[0m"
	go run $(BUILD_DIR)/main.go

# generate swagger
swag:
	@echo 'generation swagger docs'
	swag init --parseDependency -g handler.go -dir internal/api/http/blog/v1 --instanceName blogV1

# migrations
LOCAL_MIGRATION_DIR=$(CURDIR)/$(MIGRATION_DIR)

LOCAL_MIGRATION_DSN="postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_DB_NAME}?sslmode=disable"

migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

migration-create:
	@echo "Migration name:"
	@read migration_name; \
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create $$migration_name sql