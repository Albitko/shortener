// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/shortener.proto

package proto

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

const (
	Shortener_ShortenURL_FullMethodName      = "/shortener.Shortener/ShortenURL"
	Shortener_ShortenURLBatch_FullMethodName = "/shortener.Shortener/ShortenURLBatch"
	Shortener_DeleteURLBatch_FullMethodName  = "/shortener.Shortener/DeleteURLBatch"
	Shortener_GetURL_FullMethodName          = "/shortener.Shortener/GetURL"
	Shortener_GetURLsByUserID_FullMethodName = "/shortener.Shortener/GetURLsByUserID"
	Shortener_Ping_FullMethodName            = "/shortener.Shortener/Ping"
	Shortener_GetStats_FullMethodName        = "/shortener.Shortener/GetStats"
)

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	ShortenURL(ctx context.Context, in *PostURLRequest, opts ...grpc.CallOption) (*PostURLResponse, error)
	ShortenURLBatch(ctx context.Context, in *PostURLBatchRequest, opts ...grpc.CallOption) (*PostURLBatchResponse, error)
	DeleteURLBatch(ctx context.Context, in *DeleteURLBatchRequest, opts ...grpc.CallOption) (*DeleteURLBatchResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	GetURLsByUserID(ctx context.Context, in *GetURLsByUserIDRequest, opts ...grpc.CallOption) (*GetURLsByUserIDResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) ShortenURL(ctx context.Context, in *PostURLRequest, opts ...grpc.CallOption) (*PostURLResponse, error) {
	out := new(PostURLResponse)
	err := c.cc.Invoke(ctx, Shortener_ShortenURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) ShortenURLBatch(ctx context.Context, in *PostURLBatchRequest, opts ...grpc.CallOption) (*PostURLBatchResponse, error) {
	out := new(PostURLBatchResponse)
	err := c.cc.Invoke(ctx, Shortener_ShortenURLBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteURLBatch(ctx context.Context, in *DeleteURLBatchRequest, opts ...grpc.CallOption) (*DeleteURLBatchResponse, error) {
	out := new(DeleteURLBatchResponse)
	err := c.cc.Invoke(ctx, Shortener_DeleteURLBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, Shortener_GetURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetURLsByUserID(ctx context.Context, in *GetURLsByUserIDRequest, opts ...grpc.CallOption) (*GetURLsByUserIDResponse, error) {
	out := new(GetURLsByUserIDResponse)
	err := c.cc.Invoke(ctx, Shortener_GetURLsByUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, Shortener_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, Shortener_GetStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	ShortenURL(context.Context, *PostURLRequest) (*PostURLResponse, error)
	ShortenURLBatch(context.Context, *PostURLBatchRequest) (*PostURLBatchResponse, error)
	DeleteURLBatch(context.Context, *DeleteURLBatchRequest) (*DeleteURLBatchResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	GetURLsByUserID(context.Context, *GetURLsByUserIDRequest) (*GetURLsByUserIDResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) ShortenURL(context.Context, *PostURLRequest) (*PostURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURL not implemented")
}
func (UnimplementedShortenerServer) ShortenURLBatch(context.Context, *PostURLBatchRequest) (*PostURLBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURLBatch not implemented")
}
func (UnimplementedShortenerServer) DeleteURLBatch(context.Context, *DeleteURLBatchRequest) (*DeleteURLBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLBatch not implemented")
}
func (UnimplementedShortenerServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedShortenerServer) GetURLsByUserID(context.Context, *GetURLsByUserIDRequest) (*GetURLsByUserIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURLsByUserID not implemented")
}
func (UnimplementedShortenerServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortenerServer) GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_ShortenURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).ShortenURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_ShortenURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).ShortenURL(ctx, req.(*PostURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_ShortenURLBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostURLBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).ShortenURLBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_ShortenURLBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).ShortenURLBatch(ctx, req.(*PostURLBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteURLBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteURLBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_DeleteURLBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteURLBatch(ctx, req.(*DeleteURLBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetURLsByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLsByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetURLsByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetURLsByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetURLsByUserID(ctx, req.(*GetURLsByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetStats(ctx, req.(*GetStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenURL",
			Handler:    _Shortener_ShortenURL_Handler,
		},
		{
			MethodName: "ShortenURLBatch",
			Handler:    _Shortener_ShortenURLBatch_Handler,
		},
		{
			MethodName: "DeleteURLBatch",
			Handler:    _Shortener_DeleteURLBatch_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _Shortener_GetURL_Handler,
		},
		{
			MethodName: "GetURLsByUserID",
			Handler:    _Shortener_GetURLsByUserID_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Shortener_Ping_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _Shortener_GetStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}
