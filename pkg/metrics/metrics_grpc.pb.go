// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: pkg/metrics/metrics.proto

package metrics

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Metrics_AddGauge_FullMethodName   = "/metrics.Metrics/AddGauge"
	Metrics_AddCounter_FullMethodName = "/metrics.Metrics/AddCounter"
	Metrics_AddBatch_FullMethodName   = "/metrics.Metrics/AddBatch"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	AddGauge(ctx context.Context, in *MetricsItem, opts ...grpc.CallOption) (*SimpleResponse, error)
	AddCounter(ctx context.Context, in *MetricsItem, opts ...grpc.CallOption) (*SimpleResponse, error)
	AddBatch(ctx context.Context, in *MetricsList, opts ...grpc.CallOption) (*SimpleResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) AddGauge(ctx context.Context, in *MetricsItem, opts ...grpc.CallOption) (*SimpleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, Metrics_AddGauge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) AddCounter(ctx context.Context, in *MetricsItem, opts ...grpc.CallOption) (*SimpleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, Metrics_AddCounter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) AddBatch(ctx context.Context, in *MetricsList, opts ...grpc.CallOption) (*SimpleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, Metrics_AddBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility.
type MetricsServer interface {
	AddGauge(context.Context, *MetricsItem) (*SimpleResponse, error)
	AddCounter(context.Context, *MetricsItem) (*SimpleResponse, error)
	AddBatch(context.Context, *MetricsList) (*SimpleResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsServer struct{}

func (UnimplementedMetricsServer) AddGauge(context.Context, *MetricsItem) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddGauge not implemented")
}
func (UnimplementedMetricsServer) AddCounter(context.Context, *MetricsItem) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCounter not implemented")
}
func (UnimplementedMetricsServer) AddBatch(context.Context, *MetricsList) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBatch not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}
func (UnimplementedMetricsServer) testEmbeddedByValue()                 {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	// If the following call pancis, it indicates UnimplementedMetricsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_AddGauge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsItem)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddGauge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddGauge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddGauge(ctx, req.(*MetricsItem))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_AddCounter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsItem)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddCounter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddCounter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddCounter(ctx, req.(*MetricsItem))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_AddBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddBatch(ctx, req.(*MetricsList))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddGauge",
			Handler:    _Metrics_AddGauge_Handler,
		},
		{
			MethodName: "AddCounter",
			Handler:    _Metrics_AddCounter_Handler,
		},
		{
			MethodName: "AddBatch",
			Handler:    _Metrics_AddBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/metrics/metrics.proto",
}
