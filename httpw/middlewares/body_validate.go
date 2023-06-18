package middlewares

import (
	"encoding/json"
	"io"

	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

func BodyValidator[T any](reader io.ReadCloser, body *T) *viewmodels.HTTPError {
	err := json.NewDecoder(reader).Decode(body)

	if err != nil {

		// return viewmodels.BadRe
	}

	return nil
}
