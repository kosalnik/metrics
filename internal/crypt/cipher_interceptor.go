package crypt

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

//go:generate mockgen -source=cipher_interceptor.go -destination=./mock/encoder.go -package=mock
type Encoder interface {
	Encode([]byte) ([]byte, error)
}

func NewCipherInterceptor(encoder Encoder, transport http.RoundTripper) *CipherRoundTripper {
	return &CipherRoundTripper{
		core:    transport,
		encoder: encoder,
	}
}

type CipherRoundTripper struct {
	core    http.RoundTripper
	encoder Encoder
}

func (a *CipherRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if a.encoder == nil || request.Body == nil {
		return a.core.RoundTrip(request)
	}

	b, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		if request.Body != nil {
			request.Body.Close()
		}
	}()

	ciphertext, err := a.encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	req := request.Clone(context.Background())
	req.Body = io.NopCloser(bytes.NewReader(ciphertext))
	req.ContentLength = int64(len(ciphertext))

	return a.core.RoundTrip(req)
}

var _ http.RoundTripper = &CipherRoundTripper{}
