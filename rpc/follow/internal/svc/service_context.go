package svc

import (
	"pkg/orm"

	"follow/internal/config"
	"follow/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	BizRedis         redis.Redis
	FollowModel      *model.FollowModle
	FollowCountModel *model.FollowCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.BizRedis)
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})
	return &ServiceContext{
		Config:           c,
		DB:               db,
		BizRedis:         *rds,
		FollowModel:      model.NewFollowModel(db.DB),
		FollowCountModel: model.NewFollowCountModel(db.DB),
	}
}
