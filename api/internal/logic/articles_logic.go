package logic

import (
	"context"

	"api/internal/grpc_client/article"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticlesLogic) Articles(req *types.ArticlesRequest) (resp *types.ArticlesResponse, err error) {
	in := &article.ArticlesRequest{
		UserId:    req.AuthorId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		SortType:  req.SortType,
		ArticleId: req.ArticleId,
	}
	res, err := l.svcCtx.ArticleRpc.Articles(l.ctx, in)
	if err != nil {
		l.Logger.Errorf("get articles rpc error: err is %s", err.Error())
		return nil, err
	}
	if res.Articles == nil || len(res.Articles) == 0 {
		return &types.ArticlesResponse{}, nil
	}
	articleList := []types.ArticleInfo{}
	for i := 0; i < len(res.Articles); i++ {
		articleList = append(
			articleList,
			types.ArticleInfo{
				ArticleId:   res.Articles[i].Id,
				Title:       res.Articles[i].Title,
				Content:     res.Articles[i].Content,
				Description: res.Articles[i].Description,
				Cover:       res.Articles[i].Cover,
			},
		)
	}
	return &types.ArticlesResponse{Articles: articleList}, nil
}
