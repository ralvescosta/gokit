install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

tests:
	@go test ./... -v

lint:
	@golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./... > golanci-report.xml

test-cov:
	@go test ./... -v -covermode atomic -coverprofile=coverage.out