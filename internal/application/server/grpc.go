package server

import (
	"context"
	"net"

	"github.com/kosalnik/metrics/internal/log"
	pb "github.com/kosalnik/metrics/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServer
}

func NewGRPCServer(ctx context.Context, addr string) {
	log.Info().Str("addr", addr).Msg("Listen grpc")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("new grpc server fails")
	}
	s := grpc.NewServer()
	pb.RegisterMetricsServer(s, &GRPCServer{})
	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()
	if err := s.Serve(listen); err != nil {
		log.Error().Err(err).Msg("Listen grpc fails")
	}
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
