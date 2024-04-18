package crypt

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

var hashHeader string = "HashSHA256"

func SetHashHeaderName(name string) {
	hashHeader = name
}

func ExtractSign(r *http.Request) string {
	return r.Header.Get("HashSHA256")
}

func ToSignRequest(r *http.Request, value string) {
	r.Header.Set("HashSHA256", value)
}

func GetSign(data []byte) string {
	h := sha256.Sum256(data)

	return hex.EncodeToString(h[:])
}

func VerifySign(data []byte, sign string) bool {
	return GetSign(data) == sign
}
