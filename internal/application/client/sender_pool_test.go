package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/application/client"
	"github.com/kosalnik/metrics/internal/application/client/mock"
	"github.com/kosalnik/metrics/internal/models"
)

func TestSenderPool_SendBatch(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	data := []models.Metrics{{ID: "ASD", MType: models.MCounter, Delta: 3}}
	senderMock := mock.NewMockSender(ctrl)
	expectedContext := context.Background()
	senderMock.EXPECT().SendBatch(expectedContext, data).Times(5).Do(func(_ context.Context, _ []models.Metrics) {
		time.Sleep(time.Second * 2)
	})
	p := client.NewSenderPool(ctx, senderMock, 5)
	for i := 0; i < 5; i++ {
		require.NoError(t, p.SendBatch(expectedContext, data))
	}
	time.Sleep(time.Second)
}

func TestSenderPool_SendBatchContinue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	data := []models.Metrics{{ID: "ASD", MType: models.MCounter, Delta: 3}}
	senderMock := mock.NewMockSender(ctrl)
	expectedContext := context.Background()
	senderMock.EXPECT().SendBatch(expectedContext, data).Times(6).Do(func(_ context.Context, _ []models.Metrics) {
		time.Sleep(time.Millisecond * 500)
	})
	p := client.NewSenderPool(ctx, senderMock, 5)
	for i := 0; i < 5; i++ {
		require.NoError(t, p.SendBatch(expectedContext, data))
	}
	time.Sleep(time.Second)
	require.NoError(t, p.SendBatch(expectedContext, data))
	time.Sleep(time.Second)
}
