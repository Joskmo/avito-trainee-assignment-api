.PHONY: build run migrate down clean
.DEFAULT_GOAL := run

build:
	docker compose build app

run: build
	docker compose up -d

migrate:
	docker compose run --rm migrations

logs:
	docker compose logs -f app

down:
	docker compose down

clean:
	docker compose down -v
	rm -rf postgres_data
