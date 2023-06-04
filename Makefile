install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	@echo "go mod download..."
	@echo "downloading ./auth/go.mod ..."
	@cd auth && go mod download
	@echo "downloading ./configs/go.mod ..."
	@cd configs && go mod download
	@echo "downloading ./configs_builder/go.mod ..."
	@cd configs_builder && go mod download
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

update-pkgs:
	@echo "updating auth..."
	@cd auth && go get -u all && go mod tidy
	@echo "updating configs..."
	@cd configs && go get -u all && go mod tidy
	@echo "updating configs_builder..."
	@cd configs_builder && go get -u all && go mod tidy
	@echo "updating guid..."
	@cd guid && go get -u all && go mod tidy
	@echo "updating httpw..."
	@cd httpw && go get -u all && go mod tidy
	@echo "updating logging..."
	@cd logging && go get -u all && go mod tidy
	@echo "updating metrics..."
	@cd metrics && go get -u all && go mod tidy
	@echo "updating rabbitmq..."
	@cd rabbitmq && go get -u all && go mod tidy
	@echo "updating secrets_manager..."
	@cd secrets_manager && go get -u all && go mod tidy
	@echo "updating sql..."
	@cd sql && go get -u all && go mod tidy
	@echo "updating tracing..."
	@cd tracing && go get -u all && go mod tidy

tests:
	@cd configs && go test ./... -v
	@cd configs_builder && go test ./... -v
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