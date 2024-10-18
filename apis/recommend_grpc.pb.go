// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: apis/recommend.proto

package apis

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

// RecommenderClient is the client API for Recommender service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecommenderClient interface {
	RecommendServer(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*NoteResponse, error)
}

type recommenderClient struct {
	cc grpc.ClientConnInterface
}

func NewRecommenderClient(cc grpc.ClientConnInterface) RecommenderClient {
	return &recommenderClient{cc}
}

func (c *recommenderClient) RecommendServer(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*NoteResponse, error) {
	out := new(NoteResponse)
	err := c.cc.Invoke(ctx, "/recommend.server.Recommender/RecommendServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecommenderServer is the server API for Recommender service.
// All implementations must embed UnimplementedRecommenderServer
// for forward compatibility
type RecommenderServer interface {
	RecommendServer(context.Context, *UserRequest) (*NoteResponse, error)
	mustEmbedUnimplementedRecommenderServer()
}

// UnimplementedRecommenderServer must be embedded to have forward compatible implementations.
type UnimplementedRecommenderServer struct {
}

func (UnimplementedRecommenderServer) RecommendServer(context.Context, *UserRequest) (*NoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecommendServer not implemented")
}
func (UnimplementedRecommenderServer) mustEmbedUnimplementedRecommenderServer() {}

// UnsafeRecommenderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecommenderServer will
// result in compilation errors.
type UnsafeRecommenderServer interface {
	mustEmbedUnimplementedRecommenderServer()
}

func RegisterRecommenderServer(s grpc.ServiceRegistrar, srv RecommenderServer) {
	s.RegisterService(&Recommender_ServiceDesc, srv)
}

func _Recommender_RecommendServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommenderServer).RecommendServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recommend.server.Recommender/RecommendServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommenderServer).RecommendServer(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Recommender_ServiceDesc is the grpc.ServiceDesc for Recommender service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Recommender_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "recommend.server.Recommender",
	HandlerType: (*RecommenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RecommendServer",
			Handler:    _Recommender_RecommendServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/recommend.proto",
}