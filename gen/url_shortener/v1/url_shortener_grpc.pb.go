// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: url_shortener/v1/url_shortener.proto

package v1

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
	URLShortenerServiceV1_ShortenV1_FullMethodName        = "/url_shortener.v1.URLShortenerServiceV1/ShortenV1"
	URLShortenerServiceV1_ShortenBatchV1_FullMethodName   = "/url_shortener.v1.URLShortenerServiceV1/ShortenBatchV1"
	URLShortenerServiceV1_GetV1_FullMethodName            = "/url_shortener.v1.URLShortenerServiceV1/GetV1"
	URLShortenerServiceV1_PingV1_FullMethodName           = "/url_shortener.v1.URLShortenerServiceV1/PingV1"
	URLShortenerServiceV1_UserUrlsV1_FullMethodName       = "/url_shortener.v1.URLShortenerServiceV1/UserUrlsV1"
	URLShortenerServiceV1_UserUrlsDeleteV1_FullMethodName = "/url_shortener.v1.URLShortenerServiceV1/UserUrlsDeleteV1"
	URLShortenerServiceV1_StatsV1_FullMethodName          = "/url_shortener.v1.URLShortenerServiceV1/StatsV1"
)

// URLShortenerServiceV1Client is the client API for URLShortenerServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerServiceV1Client interface {
	ShortenV1(ctx context.Context, in *ShortenRequestV1, opts ...grpc.CallOption) (*ShortenResponseV1, error)
	ShortenBatchV1(ctx context.Context, in *ShortenBatchRequestV1, opts ...grpc.CallOption) (*ShortenBatchResponseV1, error)
	GetV1(ctx context.Context, in *GetRequestV1, opts ...grpc.CallOption) (*GetResponseV1, error)
	PingV1(ctx context.Context, in *PingRequestV1, opts ...grpc.CallOption) (*PingResponseV1, error)
	UserUrlsV1(ctx context.Context, in *UserUrlsRequestV1, opts ...grpc.CallOption) (*UserUrlsResponseV1, error)
	UserUrlsDeleteV1(ctx context.Context, in *UserUrlsDeleteRequestV1, opts ...grpc.CallOption) (*UserUrlsDeleteResponseV1, error)
	StatsV1(ctx context.Context, in *StatsRequestV1, opts ...grpc.CallOption) (*StatsResponseV1, error)
}

type uRLShortenerServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerServiceV1Client(cc grpc.ClientConnInterface) URLShortenerServiceV1Client {
	return &uRLShortenerServiceV1Client{cc}
}

func (c *uRLShortenerServiceV1Client) ShortenV1(ctx context.Context, in *ShortenRequestV1, opts ...grpc.CallOption) (*ShortenResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_ShortenV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceV1Client) ShortenBatchV1(ctx context.Context, in *ShortenBatchRequestV1, opts ...grpc.CallOption) (*ShortenBatchResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenBatchResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_ShortenBatchV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceV1Client) GetV1(ctx context.Context, in *GetRequestV1, opts ...grpc.CallOption) (*GetResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_GetV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceV1Client) PingV1(ctx context.Context, in *PingRequestV1, opts ...grpc.CallOption) (*PingResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_PingV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceV1Client) UserUrlsV1(ctx context.Context, in *UserUrlsRequestV1, opts ...grpc.CallOption) (*UserUrlsResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserUrlsResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_UserUrlsV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceV1Client) UserUrlsDeleteV1(ctx context.Context, in *UserUrlsDeleteRequestV1, opts ...grpc.CallOption) (*UserUrlsDeleteResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserUrlsDeleteResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_UserUrlsDeleteV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceV1Client) StatsV1(ctx context.Context, in *StatsRequestV1, opts ...grpc.CallOption) (*StatsResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatsResponseV1)
	err := c.cc.Invoke(ctx, URLShortenerServiceV1_StatsV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServiceV1Server is the server API for URLShortenerServiceV1 service.
// All implementations must embed UnimplementedURLShortenerServiceV1Server
// for forward compatibility.
type URLShortenerServiceV1Server interface {
	ShortenV1(context.Context, *ShortenRequestV1) (*ShortenResponseV1, error)
	ShortenBatchV1(context.Context, *ShortenBatchRequestV1) (*ShortenBatchResponseV1, error)
	GetV1(context.Context, *GetRequestV1) (*GetResponseV1, error)
	PingV1(context.Context, *PingRequestV1) (*PingResponseV1, error)
	UserUrlsV1(context.Context, *UserUrlsRequestV1) (*UserUrlsResponseV1, error)
	UserUrlsDeleteV1(context.Context, *UserUrlsDeleteRequestV1) (*UserUrlsDeleteResponseV1, error)
	StatsV1(context.Context, *StatsRequestV1) (*StatsResponseV1, error)
	mustEmbedUnimplementedURLShortenerServiceV1Server()
}

// UnimplementedURLShortenerServiceV1Server must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLShortenerServiceV1Server struct{}

func (UnimplementedURLShortenerServiceV1Server) ShortenV1(context.Context, *ShortenRequestV1) (*ShortenResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) ShortenBatchV1(context.Context, *ShortenBatchRequestV1) (*ShortenBatchResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenBatchV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) GetV1(context.Context, *GetRequestV1) (*GetResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) PingV1(context.Context, *PingRequestV1) (*PingResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) UserUrlsV1(context.Context, *UserUrlsRequestV1) (*UserUrlsResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserUrlsV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) UserUrlsDeleteV1(context.Context, *UserUrlsDeleteRequestV1) (*UserUrlsDeleteResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserUrlsDeleteV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) StatsV1(context.Context, *StatsRequestV1) (*StatsResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StatsV1 not implemented")
}
func (UnimplementedURLShortenerServiceV1Server) mustEmbedUnimplementedURLShortenerServiceV1Server() {}
func (UnimplementedURLShortenerServiceV1Server) testEmbeddedByValue()                               {}

// UnsafeURLShortenerServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServiceV1Server will
// result in compilation errors.
type UnsafeURLShortenerServiceV1Server interface {
	mustEmbedUnimplementedURLShortenerServiceV1Server()
}

func RegisterURLShortenerServiceV1Server(s grpc.ServiceRegistrar, srv URLShortenerServiceV1Server) {
	// If the following call pancis, it indicates UnimplementedURLShortenerServiceV1Server was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLShortenerServiceV1_ServiceDesc, srv)
}

func _URLShortenerServiceV1_ShortenV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).ShortenV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_ShortenV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).ShortenV1(ctx, req.(*ShortenRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerServiceV1_ShortenBatchV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenBatchRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).ShortenBatchV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_ShortenBatchV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).ShortenBatchV1(ctx, req.(*ShortenBatchRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerServiceV1_GetV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).GetV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_GetV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).GetV1(ctx, req.(*GetRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerServiceV1_PingV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).PingV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_PingV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).PingV1(ctx, req.(*PingRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerServiceV1_UserUrlsV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUrlsRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).UserUrlsV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_UserUrlsV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).UserUrlsV1(ctx, req.(*UserUrlsRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerServiceV1_UserUrlsDeleteV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUrlsDeleteRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).UserUrlsDeleteV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_UserUrlsDeleteV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).UserUrlsDeleteV1(ctx, req.(*UserUrlsDeleteRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerServiceV1_StatsV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceV1Server).StatsV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerServiceV1_StatsV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceV1Server).StatsV1(ctx, req.(*StatsRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortenerServiceV1_ServiceDesc is the grpc.ServiceDesc for URLShortenerServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortenerServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "url_shortener.v1.URLShortenerServiceV1",
	HandlerType: (*URLShortenerServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenV1",
			Handler:    _URLShortenerServiceV1_ShortenV1_Handler,
		},
		{
			MethodName: "ShortenBatchV1",
			Handler:    _URLShortenerServiceV1_ShortenBatchV1_Handler,
		},
		{
			MethodName: "GetV1",
			Handler:    _URLShortenerServiceV1_GetV1_Handler,
		},
		{
			MethodName: "PingV1",
			Handler:    _URLShortenerServiceV1_PingV1_Handler,
		},
		{
			MethodName: "UserUrlsV1",
			Handler:    _URLShortenerServiceV1_UserUrlsV1_Handler,
		},
		{
			MethodName: "UserUrlsDeleteV1",
			Handler:    _URLShortenerServiceV1_UserUrlsDeleteV1_Handler,
		},
		{
			MethodName: "StatsV1",
			Handler:    _URLShortenerServiceV1_StatsV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "url_shortener/v1/url_shortener.proto",
}
