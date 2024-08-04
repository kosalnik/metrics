package server

import (
	"context"

	pb "github.com/kosalnik/metrics/pkg/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServer
}

func NewGRPCServer() {

}

func (g *GRPCServer) AddGauge(context.Context, *pb.MetricsItem) (*pb.SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddGauge not implemented")
}
func (g *GRPCServer) AddCounter(context.Context, *pb.MetricsItem) (*pb.SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCounter not implemented")
}
func (g *GRPCServer) AddBatch(context.Context, *pb.MetricsList) (*pb.SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBatch not implemented")
}
