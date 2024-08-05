package crypt

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func NewCipherRequestMutator(encoder Encoder) func(req *http.Request) (*http.Request, error) {
	return func(request *http.Request) (*http.Request, error) {
		if encoder == nil || request.Body == nil {
			return request, nil
		}

		b, err := io.ReadAll(request.Body)
		if err != nil {
			return request, err
		}
		defer request.Body.Close()

		ciphertext, err := encoder.Encode(b)
		if err != nil {
			return request, err
		}

		req := request.Clone(context.Background())
		req.Body = io.NopCloser(bytes.NewReader(ciphertext))
		req.ContentLength = int64(len(ciphertext))
		return req, nil
	}
}
