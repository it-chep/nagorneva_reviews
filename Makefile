APP_NAME := nagorneva_reviews

LOCAL_BIN := $(CURDIR)/bin
-include .env
export PATH := $(PATH):$(LOCAL_BIN)

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= reviews
POSTGRES_PASSWORD ?= reviews
POSTGRES_DB ?= reviews
DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: deps run build test migrate-up migrate-down migration proto gen buf docker-up docker-down lint

deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@latest
	GOBIN=$(LOCAL_BIN) go install github.com/not-for-prod/clay/cmd/protoc-gen-goclay@v0.4.0
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@v0.36.0

run:
	go run ./cmd/$(APP_NAME)

build:
	go build -o $(LOCAL_BIN)/app ./cmd/nagorneva_reviews

test:
	go test ./...

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down


# buf mirrors medblogers_base: it validates the proto module and regenerates
# the Go, grpc-gateway, OpenAPI and Clay service descriptors together.
buf:
	buf dep update
	buf dep prune
	buf build
	buf generate --template buf.gen.yaml
	goimports -w $$(find gen/go -type f -name '*.go')

# Keep both familiar entry points.  Declaring them phony prevents the gen/
# directory from making `make gen` a silent no-op.
gen: buf

proto: buf

# Usage: make migration create_reviews_indexes
migration:
	goose -dir migrations create $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)) sql

lint:
	golangci-lint run ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down
