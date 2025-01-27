// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: internal/sort/estimate/estimate.proto

package estimate

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EstimatorClient is the client API for Estimator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EstimatorClient interface {
	EstimatorResp(ctx context.Context, in *EstimateRequest, opts ...grpc.CallOption) (*EstimateResponse, error)
}

type estimatorClient struct {
	cc grpc.ClientConnInterface
}

func NewEstimatorClient(cc grpc.ClientConnInterface) EstimatorClient {
	return &estimatorClient{cc}
}

func (c *estimatorClient) EstimatorResp(ctx context.Context, in *EstimateRequest, opts ...grpc.CallOption) (*EstimateResponse, error) {
	out := new(EstimateResponse)
	err := c.cc.Invoke(ctx, "/estimate.server.Estimator/EstimatorResp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EstimatorServer is the server API for Estimator service.
// All implementations must embed UnimplementedEstimatorServer
// for forward compatibility
type EstimatorServer interface {
	EstimatorResp(context.Context, *EstimateRequest) (*EstimateResponse, error)
	mustEmbedUnimplementedEstimatorServer()
}

// UnimplementedEstimatorServer must be embedded to have forward compatible implementations.
type UnimplementedEstimatorServer struct {
}

func (UnimplementedEstimatorServer) EstimatorResp(context.Context, *EstimateRequest) (*EstimateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EstimatorResp not implemented")
}
func (UnimplementedEstimatorServer) mustEmbedUnimplementedEstimatorServer() {}

// UnsafeEstimatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EstimatorServer will
// result in compilation errors.
type UnsafeEstimatorServer interface {
	mustEmbedUnimplementedEstimatorServer()
}

func RegisterEstimatorServer(s grpc.ServiceRegistrar, srv EstimatorServer) {
	s.RegisterService(&Estimator_ServiceDesc, srv)
}

func _Estimator_EstimatorResp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EstimateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EstimatorServer).EstimatorResp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/estimate.server.Estimator/EstimatorResp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EstimatorServer).EstimatorResp(ctx, req.(*EstimateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Estimator_ServiceDesc is the grpc.ServiceDesc for Estimator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Estimator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "estimate.server.Estimator",
	HandlerType: (*EstimatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EstimatorResp",
			Handler:    _Estimator_EstimatorResp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/sort/estimate/estimate.proto",
}
