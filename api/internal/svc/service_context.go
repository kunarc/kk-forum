package svc

import (
	"pkg/interceptors"

	"api/internal/config"
	"api/internal/grpc_client/user"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRpc  user.UserClient
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	userRpc := zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrInterceptor()))
	return &ServiceContext{
		Config:   c,
		UserRpc:  user.NewUserClient(userRpc.Conn()),
		BizRedis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
