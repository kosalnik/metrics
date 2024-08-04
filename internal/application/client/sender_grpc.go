package client

import (
	"context"

	"github.com/kosalnik/metrics/internal/models"
	pb "github.com/kosalnik/metrics/pkg/metrics"
	"google.golang.org/grpc"
)

type GRPCSender struct {
	cl pb.MetricsClient
}

func NewGRPCSender(conn grpc.ClientConnInterface) *GRPCSender {
	// устанавливаем соединение с сервером
	//conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatal().Err(err).Msg("fail grpc connect")
	//}
	//defer conn.Close()
	return &GRPCSender{cl: pb.NewMetricsClient(conn)}
}

func (G GRPCSender) SendGauge(k string, v float64) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCSender) SendCounter(k string, v int64) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCSender) SendBatch(ctx context.Context, list []models.Metrics) error {
	//TODO implement me
	panic("implement me")
}

var _ Sender = &GRPCSender{}
