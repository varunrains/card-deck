API_BINARY=apiApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## run: stops docker-compose (if running), builds all projects and starts docker compose
run: build
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build: builds the api binary
build:
	@echo Building api binary...
	chdir .\ && set CGO_ENABLED=0&& go build -o ${API_BINARY} ./cmd/api
	@echo Done!

test_cmd_api:
	@echo testing inside cmd/api ...
	go test ./cmd/api -v
	@echo "Finished testing"

test_integration:
	@echo testing integration tests with DB ...
	go test ./internal/repository/dbrepo -v
	@echo "Finished testing"

test_coverage:
	go test ./... -coverprofile=coverage.out
