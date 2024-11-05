package svc

import (
	"article/internal/config"
	"article/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(sqlx.NewMysql(c.Datasource), c.CacheRedis),
		BizRedis:     redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
