package middlewares

import (
	"encoding/json"
	"io"

	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

func BodyValidator[B any](reader io.ReadCloser, body *B) *viewmodels.HTTPError {
	err := json.NewDecoder(reader).Decode(body)

	if err != nil {
		return viewmodels.BadRequest("wrong body")
	}

	return nil
}
