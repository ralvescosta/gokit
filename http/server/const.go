package server

import "errors"

var (
	ErrorInvalidHttpMethod = errors.New("invalid http method")
)

func LogMessage(msg string) string {
	return "[Pkg::HttpServer] " + msg
}
