package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"io"
)

type RSAEncoder struct {
	random io.Reader
	key    *rsa.PublicKey
}

func NewEncoder(key *rsa.PublicKey, random io.Reader) *RSAEncoder {
	if random == nil {
		random = rand.Reader
	}
	return &RSAEncoder{random: random, key: key}
}

func (r *RSAEncoder) Encode(b []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(r.random, r.key, b)
}

type RSADecoder struct {
	random io.Reader
	key    *rsa.PrivateKey
}

func NewDecoder(key *rsa.PrivateKey, random io.Reader) *RSADecoder {
	if random == nil {
		random = rand.Reader
	}
	return &RSADecoder{key: key, random: random}
}

func (r *RSADecoder) Decode(b []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(r.random, r.key, b)
}
