.PHONY: up down nuke run build test lint lint-fix clean goose-up migrate-up

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

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

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
