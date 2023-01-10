install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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