package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"article-mq/internal/svc"
	"article-mq/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLogic) Consume(ctx context.Context, key, val string) error {
	l.Logger.Debugf("[article-mq] comsume like num, notice article: key is %v, val is %v", key, val)
	var msg *types.CanalLikeMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		l.Logger.Errorf("Unmarshal msg error: err is %v, msg is %v", err.Error(), val)
		return err
	}
	return l.updateArticleLikeNum(msg)
}

func (l *ArticleLogic) updateArticleLikeNum(msg *types.CanalLikeMsg) error {
	manageType := msg.Type
	if manageType == "DELETE" {
		return nil
	}
	for _, d := range msg.Data {
		if d.BizID != types.ArticleBizID {
			continue
		}
		id, err := strconv.ParseInt(d.ObjID, 10, 64)
		if err != nil {
			l.Logger.Errorf("strconv.ParseInt id: %s error: %v", d.ID, err)
			continue
		}
		likeNum, err := strconv.ParseInt(d.LikeNum, 10, 64)
		if err != nil {
			l.Logger.Errorf("strconv.ParseInt likeNum: %s error: %v", d.LikeNum, err)
			continue
		}
		err = l.svcCtx.ArticleModel.UpdateLikeNum(l.ctx, id, likeNum)
		if err != nil {
			l.Logger.Errorf("UpdateLikeNum id: %d like: %d", id, likeNum)
		}
	}
	return nil
}
