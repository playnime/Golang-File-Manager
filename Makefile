include .env
export

# Docker 
build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f server

restart:
	docker compose restart server

# Clean  чтобы почистить данные бд например (весь volume чистит)
clean:
	docker compose down -v --remove-orphans

# Database migrations 
migrate-up:
	migrate -path migrations -database ${CONN_STR} up

migrate-down:
	migrate -path migrations -database ${CONN_STR} down

# local
service-run:
	go run server/cmd/server/main.go



.PHONY: build up down logs restart migrate-up migrate-down service-run clean