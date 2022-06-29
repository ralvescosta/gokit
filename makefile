lint-tool:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	cd ./env && go mod download && go mod tidy
	cd ..
	cd ./logger && go mod download && go mod tidy
	cd ..
	cd ./sql && go mod download && go mod tidy
	cd ..

tests:
	go test ./env/... -v
	go test ./logger/... -v
	go test ./sql/... -v

test-ci:
	golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./env/... ./logger/... ./sql/... > golanci-report.xml
	go test ./env/... ./logger/... ./sql/... -race -covermode atomic -coverprofile=coverage.out -json > report.json