install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	@echo "1 - 7 :: download::env"
	@cd ./env && go mod download && go mod tidy

	@echo "2 - 7 :: download::logging"
	@cd ./logging && go mod download && go mod tidy

	@echo "3 - 7 :: download::sql"
	@cd ./sql && go mod download && go mod tidy

	@echo "4 - 7 :: download::uuid"
	@cd ./uuid && go mod download && go mod tidy

	@echo "5 - 7 :: download::messaging"
	@cd ./messaging && go mod download && go mod tidy

test-env:
	go test ./env/... -v

test-logging:
	go test ./logging/... -v

test-sql:
	go test ./sql/... -v

test-messaging:
	go test ./messaging/... -v	

test-uuid:
	go test ./uuid/... -v

tests:
	@go test ./env/... -v
	@go test ./logging/... -v
	@go test ./sql/... -v
	@go test ./messaging/... -v
	@go test ./uuid/... -v

lint:
	@golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./env/... ./logging/... ./sql/... ./messaging/... ./uuid/... > golanci-report.xml

test-cov:
# go test ./env/... ./logging/... ./sql/... ./messaging/... -v -race -covermode atomic -coverprofile=coverage.out -json > report.json
	@go test ./env/... ./logging/... ./sql/... ./messaging/... ./uuid/... -v -covermode atomic -coverprofile=coverage.out