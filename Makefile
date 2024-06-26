.PHONY: up down nuke run build test unit-test lint lint-fix install-golangci-lint mock clean 
.PHONY: migrate-up migrate-down migrate-create-schema migrate-create-datafix migrate-version migrate-status

export PGHOST ?= localhost
export PGPORT ?= 5432
export PGUSER ?= postgres
export PGPASSWORD ?= postgres
export PGDATABASE ?= camelhr_dev_db
export PGSSLMODE ?= disable
export REDIS_HOST ?= localhost
export REDIS_PORT ?= 6379

up:
	docker-compose -f docker-compose.yml up -d --build

down:
	docker-compose -f docker-compose.yml down --remove-orphans

nuke:
	docker-compose -f docker-compose.yml down --volumes --remove-orphans

run:
	DB_CONN=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/${PGDATABASE}?sslmode=${PGSSLMODE} REDIS_CONN=redis://${REDIS_HOST}:${REDIS_PORT} go run cmd/api/main.go

build:
	go build -o bin/ ./...

test:
	go test -v -cover -race ./...

unit-test:
	go test -v -cover -race ./... -short

lint: install-golangci-lint
	golangci-lint run

lint-fix: install-golangci-lint
	golangci-lint run --fix

install-golangci-lint:
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(shell go env GOPATH)/bin" v1.58.2

mock:
	@which mockery > /dev/null || curl -sSfL "https://github.com/vektra/mockery/releases/download/v2.43.1/mockery_2.43.1_$(shell uname -s)_$(shell uname -m).tar.gz" | tar xz -C "$(shell go env GOPATH)/bin" mockery
	mockery

clean:
	rm -rf bin

migrate-up:
	@go run ./cmd/dbmigrator/main.go -db_conn=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/${PGDATABASE}?sslmode=${PGSSLMODE} up

migrate-down:
	@go run ./cmd/dbmigrator/main.go -db_conn=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/${PGDATABASE}?sslmode=${PGSSLMODE} down

migrate-create-schema:
	@go run ./cmd/dbmigrator/main.go -type=schema create $(name) sql

migrate-create-datafix:
	@go run ./cmd/dbmigrator/main.go -type=datafix create $(name) go

migrate-version:
	@go run ./cmd/dbmigrator/main.go -db_conn=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/${PGDATABASE}?sslmode=${PGSSLMODE} version

migrate-status:
	@go run ./cmd/dbmigrator/main.go -db_conn=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/${PGDATABASE}?sslmode=${PGSSLMODE} status
