package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

var hashHeader string = "HashSHA256"

func ExtractSign(r *http.Request) string {
	return r.Header.Get("HashSHA256")
}

func ToSignRequest(r *http.Request, value string) {
	r.Header.Set("HashSHA256", value)
}

func GetSign(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}

func VerifySign(data []byte, sign string, key []byte) bool {
	return GetSign(data, key) == sign
}
