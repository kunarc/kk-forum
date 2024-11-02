package logic

import (
	"context"
	"time"

	"article/internal/code"
	"article/internal/model"
	"article/internal/svc"
	"article/internal/types"
	"article/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvaild
	}
	if len(in.Title) == 0 {
		return nil, code.ArticleTitleEmpty
	}
	if len(in.Content) < types.MinContentCount {
		return nil, code.ArticleContentSmall
	}
	res, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		AuthorId:    uint64(in.UserId),
		Title:       in.Title,
		Content:     in.Content,
		Description: in.Description,
		Cover:       in.Cover,
		Status:      types.ArticleStatusVisible, // 正常逻辑不会这样写，这里为了演示方便
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("publish insert req %v error %s", in, err.Error())
		return nil, err
	}
	articleId, err := res.LastInsertId()
	if err != nil {
		l.Logger.Errorf("get articleId error %s", err.Error())
		return nil, err
	}
	return &pb.PublishResponse{ArticleId: articleId}, nil
}
