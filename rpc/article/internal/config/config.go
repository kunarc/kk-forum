package config

import (
	"pkg/consul"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Datasource string
	Consul     consul.Conf
	CacheRedis cache.CacheConf
	BizRedis   redis.RedisConf
}
