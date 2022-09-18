package httpw

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	PrometheusMetricKind MetricKind = 1
	OtelMetricKind       MetricKind = 2
)

var (
	ErrorInvalidHttpMethod = errors.New("invalid http method")
	allowedHTTPMethods     = map[string]bool{http.MethodGet: true, http.MethodPost: true, http.MethodPut: true, http.MethodPatch: true, http.MethodDelete: true}
)

func OTLPOperationName(method, path string) string {
	return method + " " + path
}

func LogMessage(msg string) string {
	return "[gokit::httpserver] " + msg
}

func LogRouterRegister(method, path string) string {
	return LogMessage(fmt.Sprintf("registering route: %s %s", method, path))
}
