package server

import (
	"context"

	"github.com/kosalnik/metrics/internal/models"
	"github.com/kosalnik/metrics/internal/storage"
	pb "github.com/kosalnik/metrics/pkg/metrics"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServer
	storage storage.Storager
}

func (g *GRPCServer) AddGauge(ctx context.Context, in *pb.MetricsItem) (*pb.SimpleResponse, error) {
	_, err := g.storage.SetGauge(ctx, in.Id, *in.Value)
	if err != nil {
		return &pb.SimpleResponse{Error: "fail set gauge"}, nil
	}
	return &pb.SimpleResponse{Success: true}, nil
}

func (g *GRPCServer) AddCounter(ctx context.Context, in *pb.MetricsItem) (*pb.SimpleResponse, error) {
	_, err := g.storage.IncCounter(ctx, in.Id, *in.Delta)
	if err != nil {
		return &pb.SimpleResponse{Error: "fail set counter"}, nil
	}
	return &pb.SimpleResponse{Success: true}, nil
}
func (g *GRPCServer) AddBatch(ctx context.Context, in *pb.MetricsList) (*pb.SimpleResponse, error) {
	var list []models.Metrics
	for _, v := range in.Items {
		var m models.Metrics
		switch v.Type {
		case pb.MType_GAUGE:
			m = models.Metrics{ID: v.Id, MType: models.MGauge, Value: *v.Value}
		case pb.MType_COUNTER:
			m = models.Metrics{ID: v.Id, MType: models.MCounter, Delta: *v.Delta}
		}
		if m.MType == "" {
			continue
		}
		list = append(list, m)
	}
	if len(list) == 0 {
		return &pb.SimpleResponse{Success: true}, nil
	}
	if err := g.storage.UpsertAll(ctx, list); err != nil {
		return &pb.SimpleResponse{Error: "fail batch update"}, nil
	}
	return &pb.SimpleResponse{Success: true}, nil
}
