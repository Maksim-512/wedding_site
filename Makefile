.PHONY: build run docker-build docker-up docker-down docker-logs

build:
	go build -o wedding-website ./cmd/wedding_website/http

run:
	go run ./cmd/wedding_website

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

migrate:
	docker-compose exec postgres psql -U postgres -d wedding_website -f /docker-entrypoint-initdb.d/init.sql