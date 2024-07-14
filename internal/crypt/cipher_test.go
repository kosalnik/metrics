package crypt_test

import (
	"math/rand"
	"testing"

	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/stretchr/testify/require"
)

func TestEncoderDecoder(t *testing.T) {
	var seed int64 = 1
	privateKey, publicKey, err := generateRSAKeyPair(seed)
	require.NoError(t, err)
	encoder := crypt.NewEncoder(publicKey, rand.New(rand.NewSource(seed)))
	decoder := crypt.NewDecoder(privateKey, rand.New(rand.NewSource(seed)))

	want := []byte("Hello, World!")
	encoded, err := encoder.Encode(want)
	require.NoError(t, err)
	require.NotEqual(t, want, encoded)
	got, err := decoder.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
