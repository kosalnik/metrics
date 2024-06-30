package backup

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/models"
	"github.com/kosalnik/metrics/internal/storage/mock"
)

func TestDump_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f, err := os.CreateTemp(os.TempDir(), "test")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	require.NoError(t, os.Remove(f.Name()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s := mock.NewMockStorage(ctrl)
	s.EXPECT().GetAll(ctx).Return([]models.Metrics{{ID: "a", MType: models.MGauge, Value: 1.1}}, nil)

	d := NewDump(s, f.Name())
	err = d.Store(ctx)
	require.NoError(t, err)

	require.FileExists(t, f.Name())
	got, err := os.ReadFile(f.Name())
	require.NoError(t, err)
	require.JSONEq(t, `{"Data":[{"id":"a","type":"gauge","value":1.1}]}`, string(got))
}
