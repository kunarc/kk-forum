package logic

import (
	"context"
	"encoding/json"

	"api/internal/code"
	"api/internal/grpc_client/article"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const minContentCount = 100

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (resp *types.PublishResponse, err error) {
	if len(req.Title) == 0 {
		return nil, code.ArticleTitleEmpty
	}
	if len(req.Content) < minContentCount {
		return nil, code.ArticleContentSmall
	}
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		l.Logger.Errorf("get userId error err is %s", err.Error())
		return nil, err
	}
	res, err := l.svcCtx.ArticleRpc.Publish(l.ctx, &article.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Cover:       req.Cover,
	})
	if err != nil {
		l.Logger.Errorf("article public rpc error err is %s", err.Error())
		return nil, err
	}
	return &types.PublishResponse{
		ArticleId: res.ArticleId,
	}, nil
}
