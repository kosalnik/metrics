package client

import (
	"context"

	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/models"
	pb "github.com/kosalnik/metrics/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSender struct {
	client pb.MetricsClient
}

func NewGRPCSender(ctx context.Context, addr string) *GRPCSender {
	var conn *grpc.ClientConn
	var err error
	log.Info().Str("addr", addr).Msg("Dial grpc client")
	conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("fail grpc connect")
	}
	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()
	return &GRPCSender{client: pb.NewMetricsClient(conn)}
}

func (g *GRPCSender) SendGauge(k string, v float64) {
	res, err := g.client.AddGauge(context.Background(), &pb.MetricsItem{Value: &v, Type: pb.MType_GAUGE, Id: k})
	if err != nil {
		log.Error().Str("key", k).Float64("val", v).Err(err).Msg("send gauge. error")
	}
	if res.Error != "" {
		log.Error().Str("key", k).Float64("val", v).Str("err", res.Error).Msg("send gauge. response")
	}
}

func (g *GRPCSender) SendCounter(k string, v int64) {
	res, err := g.client.AddCounter(context.Background(), &pb.MetricsItem{Delta: &v, Type: pb.MType_GAUGE, Id: k})
	if err != nil {
		log.Error().Str("key", k).Int64("val", v).Err(err).Msg("send counter. error")
	}
	if res.Error != "" {
		log.Error().Str("key", k).Int64("val", v).Str("err", res.Error).Msg("send counter. response")
	}
}

func (g *GRPCSender) SendBatch(ctx context.Context, list []models.Metrics) error {
	req := pb.MetricsList{}
	for i := range list {
		v := list[i]
		switch v.MType {
		case models.MGauge:
			req.Items = append(req.Items, &pb.MetricsItem{Value: &v.Value, Id: v.ID, Type: pb.MType_GAUGE})
		case models.MCounter:
			req.Items = append(req.Items, &pb.MetricsItem{Delta: &v.Delta, Id: v.ID, Type: pb.MType_COUNTER})
		}
	}
	if len(req.Items) == 0 {
		return nil
	}
	res, err := g.client.AddBatch(ctx, &req)
	if err != nil {
		log.Error().Any("data", list).Err(err).Msg("send batch. error")
		return err
	}
	if res.Error != "" {
		log.Error().Any("data", list).Str("err", res.Error).Msg("send batch. response")
	}
	return nil
}

var _ Sender = &GRPCSender{}
