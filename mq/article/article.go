package main

import (
	"context"
	"flag"

	"article-mq/internal/config"
	"article-mq/internal/logic"
	"article-mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range Consumers(ctx, svcCtx) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, logic.NewArticleLogic(ctx, svcCtx)),
		kq.MustNewQueue(svcCtx.Config.ArticleKqConsumerConf, logic.NewArticleCacheLogic(ctx, svcCtx)),
	}
}
