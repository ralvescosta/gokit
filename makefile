environment:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

download:
	cd ./env & go mod download
	cd ./logger & go mod download

test:
	go test ./env/... -v
	go test ./logger/... -v

test-ci:
	golangci-lint run --out-format=github-actions --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- env/... logger/... > golanci-report.xml
	go test ./env/... ./logger/...  -race -covermode atomic -coverprofile=coverage.out -json > report.json