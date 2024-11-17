package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	KqConsumerConf        kq.KqConf
	ArticleKqConsumerConf kq.KqConf
	Datasource            string
	UserRpc               zrpc.RpcClientConf
	BizRedis              redis.RedisConf
	EsConfig              struct {
		Addresses []string
		Username  string
		Password  string
	}
}
