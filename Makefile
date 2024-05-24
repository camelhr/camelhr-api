.PHONY: up down nuke run build test unit-test lint lint-fix install-golangci-lint mock clean 
.PHONY: migrate-up migrate-down migrate-create-schema migrate-create-datafix migrate-version migrate-status

export PGHOST ?= localhost
export PGPORT ?= 5432
export PGUSER ?= postgres
export PGPASSWORD ?= postgres
export PGDATABASE ?= camelhr_dev_db
export PGSSLMODE ?= disable

up:
	docker-compose -f docker-compose.yml up -d --build

down:
	docker-compose -f docker-compose.yml down --remove-orphans

nuke:
	docker-compose -f docker-compose.yml down --volumes --remove-orphans

run:
	go run cmd/api/main.go

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
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

mock:
	@which mockery > /dev/null || go install github.com/vektra/mockery/v2@v2.43.1
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
