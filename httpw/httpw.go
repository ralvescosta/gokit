package httpw

import (
	"fmt"
	"net/http"
)

// HTTPError
type HTTPError struct {
	StatusCode int    `json:"statusCode" example:"400"`
	Message    string `json:"message" example:"bad request"`
	Details    any    `json:"details"`
}

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
