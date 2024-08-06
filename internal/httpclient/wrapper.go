package httpclient

import (
	"io"
	"net/http"
)

type Wrapper struct {
	cl    http.Client
	hooks []*Hook
}

type Hook func(r *http.Request)

func NewWrapper(client http.Client, hooks ...*Hook) *Wrapper {
	return &Wrapper{
		cl:    client,
		hooks: hooks,
	}
}

func (w *Wrapper) Do(req *http.Request) (*http.Response, error) {
	for i := range w.hooks {
		(*w.hooks[i])(req)
	}
	return w.cl.Do(req)
}

func (w *Wrapper) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return w.cl.Do(req)
}
