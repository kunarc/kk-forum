package logic

import (
	"context"
	"strconv"

	"api/internal/grpc_client/article"
	"api/internal/grpc_client/user"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AtricleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAtricleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AtricleDetailLogic {
	return &AtricleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AtricleDetailLogic) AtricleDetail(req *types.ArticleDetailRequest) (resp *types.ArticleDetailRespones, err error) {
	articleInfo, err := l.svcCtx.ArticleRpc.ArticleDetail(l.ctx, &article.ArticleDetailRequest{ArticleId: req.AtricleId})
	if err != nil {
		l.Logger.Errorf("get article detail rpc error: err is %v, articleId is %v", err.Error(), req.AtricleId)
		return nil, err
	}
	if articleInfo == nil || articleInfo.Article == nil {
		return nil, nil
	}
	articleItem := articleInfo.Article
	userInfo, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{UserId: articleItem.AuthorId})
	if err != nil {
		l.Logger.Errorf("get user info rpc error: err is %v, userId is %v", err.Error(), articleItem.AuthorId)
	}
	return &types.ArticleDetailRespones{
		Title:       articleItem.Title,
		Content:     articleItem.Content,
		Description: articleItem.Description,
		Cover:       articleItem.Cover,
		AuthorId:    strconv.FormatInt(articleInfo.Article.AuthorId, 10),
		AuthorName:  userInfo.Username,
	}, nil
}
