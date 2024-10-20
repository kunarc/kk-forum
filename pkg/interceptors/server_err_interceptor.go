package interceptors

import (
	"context"

	"pkg/xerror"

	"google.golang.org/grpc"
)

// ServerErrInterceptor rpc服务端错误拦截器
//
//	@return grpc.UnaryServerInterceptor
//	@author kunarc
//	@update 2024-10-20 06:33:42
func ServerErrInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			// gen custom status by err
			grpcStatus := xerror.GrpcStatusFromError(err)
			// gen err by custom status
			err = grpcStatus.Err()
		}
		return
	}
}
