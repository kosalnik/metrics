package crypt

import (
	"net/http"
)

//go:generate mockgen -source=cipher_interceptor.go -destination=./mock/encoder.go -package=mock
type Encoder interface {
	Encode([]byte) ([]byte, error)
}

func NewCipherInterceptor(encoder Encoder, transport http.RoundTripper) *CipherRoundTripper {
	return &CipherRoundTripper{
		core:    transport,
		mutator: NewCipherRequestMutator(encoder),
	}
}

type CipherRoundTripper struct {
	core    http.RoundTripper
	mutator func(r *http.Request) (*http.Request, error)
}

func (a *CipherRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	req, err := a.mutator(request)
	if err != nil {
		return nil, err
	}
	return a.core.RoundTrip(req)
}

var _ http.RoundTripper = &CipherRoundTripper{}
