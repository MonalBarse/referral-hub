.PHONY: help db-up db-down db-reset run test fmt vet check

help:
	@echo "db-up     start Postgres in Docker"
	@echo "db-down   stop Postgres"
	@echo "db-reset  drop and recreate the schema"
	@echo "run       start the API on :8090"
	@echo "test      run the test suite"
	@echo "check     fmt + vet + test"

db-up:
	cd .. && docker compose up -d

db-down:
	cd .. && docker compose down

db-reset:
	cd .. && docker compose exec -T postgres psql -U referhub -d referhub \
		-c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

check: fmt vet test
