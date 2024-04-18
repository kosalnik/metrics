package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/kosalnik/metrics/internal/infra/crypt"
)

type AddHash struct {
	core http.RoundTripper
}

func (a *AddHash) RoundTrip(request *http.Request) (*http.Response, error) {
	defer request.Body.Close()
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	h := crypt.GetSign(b)
	req := request.Clone(context.Background())
	req.Body = io.NopCloser(bytes.NewReader(b))
	crypt.ToSignRequest(req, h)

	return a.core.RoundTrip(req)
}

var _ http.RoundTripper = &AddHash{}
