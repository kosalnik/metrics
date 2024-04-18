package crypt

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/kosalnik/metrics/internal/config"
)

func VerifyHashInterceptor(cfg config.Hash) *AddHash {
	return &AddHash{
		core: http.DefaultTransport,
		key:  []byte(cfg.Key),
	}
}

type AddHash struct {
	core http.RoundTripper
	key  []byte
}

func (a *AddHash) RoundTrip(request *http.Request) (*http.Response, error) {
	if len(a.key) == 0 {
		return a.core.RoundTrip(request)
	}

	defer request.Body.Close()

	b, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	h := GetSign(b, a.key)
	req := request.Clone(context.Background())
	req.Body = io.NopCloser(bytes.NewReader(b))
	ToSignRequest(req, h)

	return a.core.RoundTrip(req)
}

var _ http.RoundTripper = &AddHash{}
