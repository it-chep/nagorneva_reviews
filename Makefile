APP_NAME := nagorneva_reviews

LOCAL_BIN := $(CURDIR)/bin
-include .env
export PATH := $(PATH):$(LOCAL_BIN)

.PHONY: deps run build test migrate-up migrate-down migration proto docker-up docker-down lint

deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@latest

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

proto:
	buf dep update
	buf lint
	buf generate

# Usage: make migration create_reviews_indexes
migration:
	goose -dir migrations create $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)) sql

lint:
	golangci-lint run ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down
