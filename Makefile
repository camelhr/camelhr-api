.PHONY: run build test up down nuke

run:
	go run cmd/api/main.go

build:
	go build -o bin/ ./...

test:
	go test -v -cover ./...

up:
	docker-compose -f docker-compose.yml up -d --build

down:
	docker-compose -f docker-compose.yml down --remove-orphans

nuke:
	docker-compose -f docker-compose.yml down --volumes --remove-orphans
