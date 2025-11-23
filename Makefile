.PHONY: build run migrate down clean
.DEFAULT_GOAL := run

build:
	docker compose build app

run: build
	docker compose up -d

migrate:
	docker compose run --rm migrations

linter:
	golangci-lint run

logs:
	docker compose logs -f app

test-e2e:
	go test -v -tags=e2e ./tests/...

load-test:
	k6 run load_test.js

down:
	docker compose down

clean:
	docker compose down -v
	rm -rf postgres_data

test: linter test-e2e load-test
