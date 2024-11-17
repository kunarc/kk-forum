package svc

import (
	"pkg/es"

	"article-mq/internal/config"
	"article-mq/internal/grpc_client/user"
	"article-mq/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis
	UserRpc      user.UserClient
	Es           *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	esConfig := &es.Config{
		Addresses: c.EsConfig.Addresses,
		Username:  c.EsConfig.Username,
		Password:  c.EsConfig.Password,
	}
	userRpc := zrpc.MustNewClient(c.UserRpc)
	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(sqlx.NewMysql(c.Datasource)),
		BizRedis:     redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		UserRpc:      user.NewUserClient(userRpc.Conn()),
		Es:           es.MustNewEs(esConfig),
	}
}
