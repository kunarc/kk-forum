package svc

import (
	"api/internal/config"
	"api/internal/grpc_client/article"
	"api/internal/grpc_client/user"
	"pkg/interceptors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	ConnectTimeout   = 1
	ReadWriteTimeout = 3
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    user.UserClient
	ArticleRpc article.ArticleClient
	BizRedis   *redis.Redis
	Oss        *oss.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Oss.ConnectTimeout == 0 {
		c.Oss.ConnectTimeout = ConnectTimeout
	}
	if c.Oss.ReadWriteTimeout == 0 {
		c.Oss.ReadWriteTimeout = ReadWriteTimeout
	}
	oss, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKeyId, c.Oss.AccessKeySecret, oss.Timeout(c.Oss.ConnectTimeout, c.Oss.ReadWriteTimeout))
	if err != nil {
		panic(err)
	}
	userRpc := zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrInterceptor()))
	aritcleRpc := zrpc.MustNewClient(c.ArticleRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrInterceptor()))
	return &ServiceContext{
		Config:     c,
		UserRpc:    user.NewUserClient(userRpc.Conn()),
		ArticleRpc: article.NewArticleClient(aritcleRpc.Conn()),
		BizRedis:   redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		Oss:        oss,
	}
}
