package logic

import (
	"context"
	"errors"

	"article/internal/svc"
	"article/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *pb.ArticleDetailRequest) (*pb.ArticleDetailResponse, error) {
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, uint64(in.ArticleId))
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &pb.ArticleDetailResponse{}, nil
		}
		return nil, err
	}
	return &pb.ArticleDetailResponse{
		Article: &pb.ArticleItem{
			Id:          int64(article.Id),
			Title:       article.Title,
			Content:     article.Content,
			Description: article.Description,
			Cover:       article.Cover,
			AuthorId:    int64(article.AuthorId),
			LikeCount:   article.LikeNum,
			PublishTime: article.PublishTime.Unix(),
		},
	}, nil
}
