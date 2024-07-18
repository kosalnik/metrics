package crypt_test

import (
	"crypto/rsa"
	"math/rand"
)

func generateRSAKeyPair(seed int64) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.New(rand.NewSource(seed)), 2048)
	if err != nil {
		return
	}
	publicKey = &privateKey.PublicKey
	return
}
