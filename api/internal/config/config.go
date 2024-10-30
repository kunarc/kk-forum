package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc    zrpc.RpcClientConf
	ArticleRpc zrpc.RpcClientConf
	BizRedis   redis.RedisConf
	Oss        struct {
		Endpoint         string
		AccessKeyId      string
		AccessKeySecret  string
		BucketName       string
		ConnectTimeout   int64 `json:optional`
		ReadWriteTimeout int64 `json:optional`
	}
}
