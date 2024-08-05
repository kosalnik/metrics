package client

import (
	"io"
	"net/http"
)

type HttpClientWrapper struct {
	origin   *http.Client
	mutators []Mutator
}

var _ HttpSender = &HttpClientWrapper{}

type HttpClientWrapperOpt = func(c *HttpClientWrapper)

func NewHttpClient(opts ...HttpClientWrapperOpt) *HttpClientWrapper {
	c := &HttpClientWrapper{
		origin: &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (h HttpClientWrapper) Do(req *http.Request) (*http.Response, error) {
	req, err := h.applyMutators(req)
	if err != nil {
		return nil, err
	}
	return h.origin.Do(req)
}

func (h HttpClientWrapper) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return h.Do(req)
}

func (h HttpClientWrapper) applyMutators(req *http.Request) (*http.Request, error) {
	var err error
	for _, m := range h.mutators {
		req, err = m(req)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

type Mutator func(req *http.Request) (*http.Request, error)

func WithMutators(m ...Mutator) HttpClientWrapperOpt {
	return func(c *HttpClientWrapper) {
		c.mutators = m
	}
}
