package crypt

import (
	"net/http"
)

func VerifyHashInterceptor(cfg Config, transport http.RoundTripper) *AddHash {
	return &AddHash{
		core:   transport,
		signer: NewSignMutator([]byte(cfg.Key)),
	}
}

type AddHash struct {
	core   http.RoundTripper
	signer func(r *http.Request) (*http.Request, error)
}

func (a *AddHash) RoundTrip(request *http.Request) (*http.Response, error) {
	req, err := a.signer(request)
	if err != nil {
		return nil, err
	}
	return a.core.RoundTrip(req)
}

var _ http.RoundTripper = &AddHash{}
