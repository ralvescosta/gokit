package httpw

import (
	"fmt"
	"net/http"
)

var (
	allowedHTTPMethods = map[string]bool{http.MethodGet: true, http.MethodPost: true, http.MethodPut: true, http.MethodPatch: true, http.MethodDelete: true}
)

func OTLPOperationName(method, path string) string {
	return method + " " + path
}

func Message(msg string) string {
	return "[httpw] " + msg
}

func LogRouterRegister(method, path string) string {
	return Message(fmt.Sprintf("registering route: %s %s", method, path))
}
