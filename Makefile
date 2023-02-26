install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	@echo "go mod download"
	@cd env && go mod download
	@echo "downloading ./guid/go.mod ..."
	@cd guid && go mod download
	@echo "downloading ./httpw/go.mod ..."
	@cd httpw && go mod download
	@echo "downloading ./logging/go.mod ..."
	@cd logging && go mod download
	@echo "downloading ./metrics/go.mod ..."
	@cd metrics && go mod download
	@echo "downloading ./rabbitmq/go.mod ..."
	@cd rabbitmq && go mod download
	@echo "downloading ./secrets_manager/go.mod ..."
	@cd secrets_manager && go mod download
	@echo "downloading ./sql/go.mod ..."
	@cd sql && go mod download
	@echo "downloading ./tracing/go.mod ..."
	@cd tracing && go mod download
	@echo "modules downloded"

tests:
	@cd env && go test ./... -v
	@cd guid && go test ./... -v
	@cd httpw && go test ./... -v
	@cd logging && go test ./... -v
	@cd metrics && go test ./... -v
	@cd rabbitmq && go test ./... -v
	@cd secrets_manager && go test ./... -v
	@cd sql && go test ./... -v
	@cd tracing && go test ./... -v

lint:
	@golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./... > golanci-report.xml

test-cov:
	@go test ./... -v -covermode atomic -coverprofile=coverage.out