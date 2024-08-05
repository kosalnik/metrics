package client

import (
	"net/http"

	"github.com/kosalnik/metrics/internal/util"
)

var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func SetRealIPMutator() func(r *http.Request) (*http.Request, error) {
	return func(request *http.Request) (*http.Request, error) {
		myIP, err := util.GetMyHostIP()
		if err != nil {
			return nil, err
		}
		request.Header.Set(xRealIP, myIP.String())
		return request, nil
	}
}
