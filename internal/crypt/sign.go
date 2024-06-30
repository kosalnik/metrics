// Package crypt contains working with signs.
package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

const hashHeader = "HashSHA256"

func ExtractSign(r *http.Request) string {
	return r.Header.Get(hashHeader)
}

func ToSignRequest(r *http.Request, value string) {
	r.Header.Set(hashHeader, value)
}

func GetSign(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}

func VerifySign(data []byte, sign string, key []byte) bool {
	return GetSign(data, key) == sign
}
