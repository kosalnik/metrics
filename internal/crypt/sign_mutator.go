package crypt

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func NewSignMutator(key []byte) func(*http.Request) (*http.Request, error) {
	return func(request *http.Request) (*http.Request, error) {
		if len(key) == 0 || request.Body == nil {
			return request, nil
		}

		b, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		defer request.Body.Close()

		h := GetSign(b, key)
		req := request.Clone(context.Background())
		req.Body = io.NopCloser(bytes.NewReader(b))
		ToSignRequest(req, h)
		return req, nil
	}
}
