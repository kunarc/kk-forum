package interceptors

import (
	"context"

	"pkg/xerror"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ClientErrInterceptor rpc客户端拦截器
//
//	@return grpc.UnaryClientInterceptor
//	@author kunarc
//	@update 2024-10-20 10:22:11
func ClientErrInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			// get status from err and get custom status
			grpcStatus, _ := status.FromError(err)
			ex := xerror.GrpcStatus2XError(grpcStatus)
			err = errors.WithMessage(ex, ex.Error())
		}
		return
	}
}
