package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"article-mq/internal/grpc_client/user"
	"article-mq/internal/svc"
	"article-mq/internal/types"

	"github.com/elastic/go-elasticsearch/v8/esutil"
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
	return l.articleOperate(msg)
}

func (l *ArticleCacheLogic) articleOperate(msg *types.CanalArticleMsg) (err error) {
	handleType := msg.Type
	var esData []*types.ArticleEsMsg
	for _, data := range msg.Data {
		status, _ := strconv.Atoi(data.Status)
		likNum, _ := strconv.ParseInt(data.LikeNum, 10, 64)
		articleId, _ := strconv.ParseInt(data.ID, 10, 64)
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
		u, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
			UserId: authorId,
		})
		if err != nil {
			l.Logger.Errorf("FindById userId: %d error: %v", authorId, err)
			return err
		}
		esData = append(esData, &types.ArticleEsMsg{
			ArticleId:   articleId,
			AuthorId:    authorId,
			AuthorName:  u.Username,
			Title:       data.Title,
			Content:     data.Content,
			Description: data.Description,
			Status:      status,
			LikeNum:     likNum,
		})
	}
	err = l.batchArticleToEs(esData)
	if err != nil {
		l.Logger.Errorf("batchArticleToEs data: %v error: %v", esData, err)
	}

	return err
}

func (l *ArticleCacheLogic) batchArticleToEs(data []*types.ArticleEsMsg) error {
	if len(data) == 0 {
		return nil
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "article-index",
	})
	if err != nil {
		return err
	}
	for _, d := range data {
		v, err := json.Marshal(d)
		if err != nil {
			return err
		}
		payload := fmt.Sprintf(`{"doc":%s,"doc_as_upsert":true}`, string(v))
		err = bi.Add(l.ctx, esutil.BulkIndexerItem{
			Action:     "update",
			DocumentID: fmt.Sprintf("%d", d.ArticleId),
			Body:       strings.NewReader(payload),
			OnSuccess:  func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {},
			OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
			},
		})
		if err != nil {
			return err
		}
	}
	return bi.Close(l.ctx)
}

func articlesKey(uid int64, sortType int32) (key string) {
	if sortType == types.PublishTimeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticlePublishTimePrefix, uid)
	} else if sortType == types.LikeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticleLikePrefix, uid)
	}
	return
}
