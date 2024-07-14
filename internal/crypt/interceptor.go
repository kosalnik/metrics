package crypt

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func VerifyHashInterceptor(cfg Config, transport http.RoundTripper) *AddHash {
	return &AddHash{
		core: transport,
		key:  []byte(cfg.Key),
	}
}

type AddHash struct {
	core http.RoundTripper
	key  []byte
}

func (a *AddHash) RoundTrip(request *http.Request) (*http.Response, error) {
	if len(a.key) == 0 || request.Body == nil {
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

	h := GetSign(b, a.key)
	req := request.Clone(context.Background())
	req.Body = io.NopCloser(bytes.NewReader(b))
	ToSignRequest(req, h)

	return a.core.RoundTrip(req)
}

var _ http.RoundTripper = &AddHash{}
