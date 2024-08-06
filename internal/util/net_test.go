package util_test

import (
	"testing"

	"github.com/kosalnik/metrics/internal/util"
	"github.com/stretchr/testify/require"
)

func TestGetMyHostIP(t *testing.T) {
	ip, err := util.GetMyHostIP()
	require.NoError(t, err)
	require.NotNil(t, ip)
}
