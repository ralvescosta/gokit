module github.com/ralvescosta/gokit/mqtt

go 1.22.2

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/ralvescosta/gokit/configs v1.20.0
	github.com/ralvescosta/gokit/logging v1.20.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ralvescosta/gokit/configs => ../configs

replace github.com/ralvescosta/gokit/logging => ../logging
