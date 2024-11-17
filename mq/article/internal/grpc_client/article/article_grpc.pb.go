// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: article.proto

package article

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

// ArticleClient is the client API for Article service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArticleClient interface {
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	Articles(ctx context.Context, in *ArticlesRequest, opts ...grpc.CallOption) (*ArticlesResponse, error)
	ArticleDelete(ctx context.Context, in *ArticleDeleteRequest, opts ...grpc.CallOption) (*ArticleDeleteResponse, error)
	ArticleDetail(ctx context.Context, in *ArticleDetailRequest, opts ...grpc.CallOption) (*ArticleDetailResponse, error)
}

type articleClient struct {
	cc grpc.ClientConnInterface
}

func NewArticleClient(cc grpc.ClientConnInterface) ArticleClient {
	return &articleClient{cc}
}

func (c *articleClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, "/pb.Article/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleClient) Articles(ctx context.Context, in *ArticlesRequest, opts ...grpc.CallOption) (*ArticlesResponse, error) {
	out := new(ArticlesResponse)
	err := c.cc.Invoke(ctx, "/pb.Article/Articles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleClient) ArticleDelete(ctx context.Context, in *ArticleDeleteRequest, opts ...grpc.CallOption) (*ArticleDeleteResponse, error) {
	out := new(ArticleDeleteResponse)
	err := c.cc.Invoke(ctx, "/pb.Article/ArticleDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleClient) ArticleDetail(ctx context.Context, in *ArticleDetailRequest, opts ...grpc.CallOption) (*ArticleDetailResponse, error) {
	out := new(ArticleDetailResponse)
	err := c.cc.Invoke(ctx, "/pb.Article/ArticleDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArticleServer is the server API for Article service.
// All implementations must embed UnimplementedArticleServer
// for forward compatibility
type ArticleServer interface {
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	Articles(context.Context, *ArticlesRequest) (*ArticlesResponse, error)
	ArticleDelete(context.Context, *ArticleDeleteRequest) (*ArticleDeleteResponse, error)
	ArticleDetail(context.Context, *ArticleDetailRequest) (*ArticleDetailResponse, error)
	mustEmbedUnimplementedArticleServer()
}

// UnimplementedArticleServer must be embedded to have forward compatible implementations.
type UnimplementedArticleServer struct {
}

func (UnimplementedArticleServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedArticleServer) Articles(context.Context, *ArticlesRequest) (*ArticlesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Articles not implemented")
}
func (UnimplementedArticleServer) ArticleDelete(context.Context, *ArticleDeleteRequest) (*ArticleDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArticleDelete not implemented")
}
func (UnimplementedArticleServer) ArticleDetail(context.Context, *ArticleDetailRequest) (*ArticleDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArticleDetail not implemented")
}
func (UnimplementedArticleServer) mustEmbedUnimplementedArticleServer() {}

// UnsafeArticleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArticleServer will
// result in compilation errors.
type UnsafeArticleServer interface {
	mustEmbedUnimplementedArticleServer()
}

func RegisterArticleServer(s grpc.ServiceRegistrar, srv ArticleServer) {
	s.RegisterService(&Article_ServiceDesc, srv)
}

func _Article_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Article/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Article_Articles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticlesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServer).Articles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Article/Articles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServer).Articles(ctx, req.(*ArticlesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Article_ArticleDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServer).ArticleDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Article/ArticleDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServer).ArticleDelete(ctx, req.(*ArticleDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Article_ArticleDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServer).ArticleDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Article/ArticleDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServer).ArticleDetail(ctx, req.(*ArticleDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Article_ServiceDesc is the grpc.ServiceDesc for Article service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Article_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Article",
	HandlerType: (*ArticleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Article_Publish_Handler,
		},
		{
			MethodName: "Articles",
			Handler:    _Article_Articles_Handler,
		},
		{
			MethodName: "ArticleDelete",
			Handler:    _Article_ArticleDelete_Handler,
		},
		{
			MethodName: "ArticleDetail",
			Handler:    _Article_ArticleDetail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "article.proto",
}
