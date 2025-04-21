install:
	@go install github.com/swaggo/swag/cmd/swag@latest

swagger-gen:
	@swag init

download:
	go mod download

run: download
	GO_ENV=local go run main.go api
