package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_standardAnalyzers(t *testing.T) {
	require.NotEmpty(t, standardAnalyzers())
}

func Test_staticCheckAnalyzers(t *testing.T) {
	require.NotEmpty(t, staticCheckAnalyzers())
}
