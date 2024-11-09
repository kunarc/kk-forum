package logic

import (
	"context"

	"article/internal/code"
	"article/internal/svc"
	"article/internal/types"
	"article/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(in *pb.ArticleDeleteRequest) (*pb.ArticleDeleteResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvaild
	}
	if in.ArticleId <= 0 {
		return nil, code.ArticleIdInvaild
	}
	err := l.svcCtx.ArticleModel.Delete(l.ctx, uint64(in.ArticleId))
	if err != nil {
		l.Logger.Errorf("delete article by id error: err is %v, id is %v", err.Error(), in.ArticleId)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, articlesKey(in.UserId, types.PublishTimeSortType), in.ArticleId)
	if err != nil {
		l.Logger.Errorf("ZremCtx req: %v error: %v", in, err)
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, articlesKey(in.UserId, types.LikeSortType), in.ArticleId)
	if err != nil {
		l.Logger.Errorf("ZremCtx req: %v error: %v", in, err)
	}
	return &pb.ArticleDeleteResponse{}, nil
}
