APP_NAME := nagorneva_reviews

LOCAL_BIN := $(CURDIR)/bin
-include .env
export PATH := $(LOCAL_BIN):$(PATH)

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= reviews
POSTGRES_PASSWORD ?= reviews
POSTGRES_DB ?= reviews
DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: deps run build test test-race test-integration migrate-up migrate-down migration proto gen buf docker-up docker-down lint

deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v1.50.0
	GOBIN=$(LOCAL_BIN) go install github.com/not-for-prod/clay/cmd/protoc-gen-goclay@v0.4.0
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0

run:
	go run ./cmd/$(APP_NAME)

build:
	go build -o $(LOCAL_BIN)/app ./cmd/nagorneva_reviews

test:
	go test ./...

test-race:
	go test -race ./...

# Use a disposable database; tests create and remove their own schema.
test-integration:
	go test -race -tags=integration -count=1 ./e2e/...

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down


# Generate Go messages, gRPC, grpc-gateway, OpenAPI and Clay descriptors.
buf:
	buf build
	buf generate --template buf.gen.yaml


# Keep both familiar entry points.  Declaring them phony prevents the gen/
# directory from making `make gen` a silent no-op.
gen: buf

proto: buf

# Usage: make migration create_reviews_indexes
migration:
	goose -dir migrations create $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)) sql

lint:
	golangci-lint run --build-tags=integration ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down
