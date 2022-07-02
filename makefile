install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	@echo "1 - 5 :: download::env"
	@cd ./env && go mod download && go mod tidy

	@echo "2 - 5 :: download::logging"
	@cd ./logging && go mod download && go mod tidy

	@echo "3 - 5 :: download::sql"
	@cd ./sql && go mod download && go mod tidy

	@echo "4 - 5 :: download::uuid"
	@cd ./uuid && go mod download && go mod tidy

tests:
	go test ./env/... -v
	go test ./logging/... -v
	go test ./sql/... -v

test-ci:
	golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./env/... ./logging/... ./sql/... > golanci-report.xml
	go test ./env/... ./logging/... ./sql/... -race -covermode atomic -coverprofile=coverage.out -json > report.json