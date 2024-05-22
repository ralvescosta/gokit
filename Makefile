LIBS := auth configs configs_builder guid httpw logging metrics mqtt rabbitmq secrets_manager sql tiny_http tracing

install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sudo sh -s -- -b /bin v2.19.0

download:
	@echo "Downloading external packages..."
	@for dir in ${LIBS}; do \
		echo "Downloading packages for: $$dir..."; \
		cd $$dir; \
		go mod download; \
		cd ..; \
	done
	@echo "External packages downloaded successfully!"

update:
	@echo "Updating external packages..."
	@for dir in ${LIBS}; do \
		echo "Updating packages for: $$dir..."; \
		cd $$dir; \
		go get -u all; \
		go mod tidy; \
		cd ..; \
	done
	@echo "External packages updated successfully!"

tests:
	@echo "Running unit tests..."
	@for dir in ${LIBS}; do \
		echo "Testing package: $$dir..."; \
		cd $$dir; \
		go test ./... -v -covermode atomic -coverprofile=coverage.out; \
		cd ..; \
	done
	@echo "All unit test runned successfully!"

lint:
	@echo "Running golangci-lint..."
	@for dir in ${LIBS}; do \
		echo "Testing package: $$dir..."; \
		cd $$dir; \
		golangci-lint run --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./...; \
		cd ..; \
	done

gosec:
	gosec -quiet ./...

push: lint gosec
	git push

test-cov:
	@go test ./... -v -covermode atomic -coverprofile=coverage.out
