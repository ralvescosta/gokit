install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	@echo "1 - 7 :: download::env"
	@cd ./env && go mod download && go mod tidy

	@echo "2 - 7 :: download::logging"
	@cd ./logging && go mod download && go mod tidy

	@echo "3 - 7 :: download::sql"
	@cd ./sql && go mod download && go mod tidy

	@echo "4 - 7 :: download::guid"
	@cd ./guid && go mod download && go mod tidy

	@echo "5 - 7 :: download::messaging"
	@cd ./messaging && go mod download && go mod tidy

	@echo "6 - 7 :: download::telemetry"
	@cd ./telemetry && go mod download && go mod tidy

t-env:
	go test github.com/ralvescosta/gokit/env/... -v

t-logging:
	go test github.com/ralvescosta/gokit/logging/... -v

t-sql:
	go test github.com/ralvescosta/gokit/sql/... -v

t-messaging:
	go test github.com/ralvescosta/gokit/messaging/... -v	

t-guid:
	go test github.com/ralvescosta/gokit/guid/... -v

tests:
	@go test github.com/ralvescosta/gokit/env/... -v
	@go test github.com/ralvescosta/gokit/logging/... -v
	@go test github.com/ralvescosta/gokit/sql/... -v
	@go test github.com/ralvescosta/gokit/messaging/... -v
	@go test github.com/ralvescosta/gokit/guid/... -v

lint:
	@golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- \
		./env/... \
		./logging/... \
		./sql/... \
		./messaging/... \
		./guid/... > golanci-report.xml

test-cov:
	@go test \
		github.com/ralvescosta/gokit/env/... \
		github.com/ralvescosta/gokit/logging/... \
		github.com/ralvescosta/gokit/sql/... \
		github.com/ralvescosta/gokit/messaging/... \
		github.com/ralvescosta/gokit/guid/... \
		-v -covermode atomic -coverprofile=coverage.out