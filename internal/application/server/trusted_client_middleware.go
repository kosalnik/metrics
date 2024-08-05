package server

import (
	"net"
	"net/http"

	"github.com/kosalnik/metrics/internal/log"
)

var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func TrustedClientMiddleware(trustedSubnet *net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				if trustedSubnet != nil {
					realIP := request.Header.Get(xRealIP)
					log.Info().Str("ip", realIP).Msg("request from client")
					if realIP == "" || !trustedSubnet.Contains(net.ParseIP(realIP)) {
						log.Warn().Msg("access denied")
						http.Error(writer, "Access denied", http.StatusForbidden)
						return
					}
				}
				next.ServeHTTP(writer, request)
			},
		)
	}
}
