.PHONY: up down nuke run build test lint clean goose-up migrate-up

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
	go test -v -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf bin

migrate-up:
	go run ./cmd/dbmigrator/main.go -db_conn=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/${PGDATABASE}?sslmode=${PGSSLMODE} up
