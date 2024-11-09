package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"article-mq/internal/svc"
	"article-mq/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleCacheLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleCacheLogic {
	return &ArticleCacheLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleCacheLogic) Consume(ctx context.Context, key, val string) error {
	l.Logger.Debugf("[article-cache-mq] comsume article, notice cache article: key is %v, val is %v", key, val)
	var msg *types.CanalArticleMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		l.Logger.Errorf("Unmarshal msg error: err is %v, msg is %v", err.Error(), val)
		return err
	}
	l.articleOperate(msg)

	return nil
}

func (l *ArticleCacheLogic) articleOperate(msg *types.CanalArticleMsg) (err error) {
	handleType := msg.Type
	for _, data := range msg.Data {
		authorId, _ := strconv.ParseInt(data.AuthorId, 10, 64)
		publishTime, _ := time.Parse("2006-01-02 15:04:05", data.PublishTime)
		likeCount, _ := strconv.ParseInt(data.LikeNum, 10, 64)
		publishTimeKey := articlesKey(authorId, types.PublishTimeSortType)
		likeCountKey := articlesKey(authorId, types.LikeSortType)
		if handleType == "INSERT" || handleType == "UPDATE" {
			_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, publishTime.Unix(), data.ID)
			if err != nil {
				l.Logger.Errorf("[ZddCtx] cache key: publishTimeKey, val: %v error: err is %v", data.ID, err.Error())
			}
			_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeCountKey, likeCount, data.ID)
			if err != nil {
				l.Logger.Errorf("[ZddCtx] cache key: likeCountKey, val: %v error: err is %v", data.ID, err.Error())
			}
		} else {
			_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, publishTimeKey, data.ID)
			if err != nil {
				l.Logger.Errorf("[ZremCtx] del key: publishTimeKey, val: %v error: err is %v", data.ID, err.Error())
			}
			_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, likeCountKey, data.ID)
			if err != nil {
				l.Logger.Errorf("[ZremCtx] del key: likeCountKey, val: %v error: err is %v", data.ID, err.Error())
			}
		}
	}
	MsgtToEs()
	return
}

func MsgtToEs() {
}

func articlesKey(uid int64, sortType int32) (key string) {
	if sortType == types.PublishTimeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticlePublishTimePrefix, uid)
	} else if sortType == types.LikeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticleLikePrefix, uid)
	}
	return
}
