package logic

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"

	"article/internal/code"
	"article/internal/model"
	"article/internal/svc"
	"article/internal/types"
	"article/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

const articlesExpire = 3600 * 24 * 2

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *pb.ArticlesRequest) (*pb.ArticlesResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvaild
	}
	if in.SortType < types.PublishTimeSortType || in.SortType > types.LikeSortType {
		return nil, code.ArticleSortTypeError
	}
	if in.PageSize > types.DefaultArticleLimit {
		return nil, code.ArticlePageSizeBig
	} else if in.PageSize <= 0 {
		return nil, code.ArticlePageSizeError
	}
	// 先读缓存
	var (
		isCache        = false
		isEnd          = false
		err            error
		articles       []*model.Article
		curPage        []*pb.ArticleItem
		lastId, cursor int64
	)
	articleIds, _ := l.cacheArticles(in.UserId, in.Cursor, int(in.PageSize), in.SortType)
	count := len(articleIds)
	// NOTE：isEnd 暂时不考虑
	// if count == 0 {
	// 	var key string
	// 	if in.SortType == types.PublishTimeSortType {
	// 		key = fmt.Sprintf("%s%d", types.CacheArticlePublishTimePrefix, in.UserId)
	// 	} else {
	// 		key = fmt.Sprintf("%s%d", types.CacheArticleLikePrefix, in.UserId)
	// 	}
	// 	_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, key, 0, strconv.Itoa(-1))
	// 	if err != nil {
	// 		l.Logger.Errorf("add -1 EOF to zset error: err is %v", err.Error())
	// 	}
	// }
	if count > 0 {
		isCache = true
		if articleIds[count-1] == -1 {
			isEnd = true
		}
		articles, err = l.articleByIds(l.ctx, articleIds)
		if err != nil {
			l.Logger.Errorf("[mapreduce] articleByIds error: err is %v", err.Error())
			return nil, err
		}
		// 通过sortFiled对articles进行排序
		var cmpFunc func(a, b *model.Article) int
		if in.SortType == types.PublishTimeSortType {
			cmpFunc = func(a, b *model.Article) int {
				return cmp.Compare(b.LikeNum, a.LikeNum)
			}
		} else {
			cmpFunc = func(a, b *model.Article) int {
				return cmp.Compare(b.PublishTime.Unix(), a.PublishTime.Unix())
			}
		}
		slices.SortFunc(articles, cmpFunc)
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("ArticlesByUserId:%d:%d", in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.ArticleModel.FindIdsByCursor(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
		})
		if err != nil {
			logx.Errorf("FindIdsByCursor userId: %d  error: %v", in.UserId, err)
			return nil, err
		}
		if v == nil {
			return &pb.ArticlesResponse{}, nil
		}
		articles = v.([]*model.Article)
	}
	// TODO: 时间时区问题， redis的时间戳与前端传入的不一致
	for _, a := range articles {
		curPage = append(
			curPage,
			&pb.ArticleItem{
				Id:           int64(a.Id),
				Title:        a.Title,
				Content:      a.Content,
				Description:  a.Description,
				Cover:        a.Cover,
				CommentCount: a.CommentNum,
				LikeCount:    a.LikeNum,
				PublishTime:  a.PublishTime.Unix(),
				AuthorId:     int64(a.AuthorId),
			},
		)
	}
	if len(curPage) > 0 {
		lastId = curPage[len(curPage)-1].Id
		if in.SortType == types.PublishTimeSortType {
			cursor = curPage[len(curPage)-1].PublishTime
		} else if in.SortType == types.LikeSortType {
			cursor = curPage[len(curPage)-1].LikeCount
		}
		if cursor < 0 {
			cursor = 0
		}
		// 可能存在重复数据（多个 article publish_time/likeNum 相同）
		for k, a := range curPage {
			if in.SortType == types.PublishTimeSortType {
				if a.Id == in.ArticleId && a.PublishTime == in.Cursor {
					curPage = curPage[k:]
					break
				}
			} else if in.SortType == types.LikeSortType {
				if a.Id == in.ArticleId && a.LikeCount == in.Cursor {
					curPage = curPage[k:]
					break
				}
			}
		}
	}

	// 缓存一手
	if !isCache {
		threading.GoSafe(func() {
			for _, cur := range curPage {
				err = l.addCacheArticles(cur, in.UserId, cur.PublishTime, cur.LikeCount, in.SortType)
				if err != nil {
					l.Logger.Errorf("cache articleId error: err is %v, curId is %v", err.Error(), cur.Id)
					break
				}
			}
		})
	}
	return &pb.ArticlesResponse{
		Articles:  curPage,
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
	}, nil
}

func (l *ArticlesLogic) cacheArticles(uId, cursor int64, ps int, sortType int32) ([]int64, error) {
	var key string
	if sortType == types.PublishTimeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticlePublishTimePrefix, uId)
	} else {
		key = fmt.Sprintf("%s%d", types.CacheArticleLikePrefix, uId)
	}
	b, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if err != nil {
		l.Logger.Errorf("ExistsCtx key: %s error: %v", key, err)
		return nil, err
	}
	if b {
		err = l.svcCtx.BizRedis.ExpireCtx(l.ctx, key, articlesExpire)
		if err != nil {
			l.Logger.Errorf("ExpireCtx key: %s error: %v", key, err)
			return nil, err
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(l.ctx, key, 0, cursor, 0, ps)
	if err != nil {
		l.Logger.Errorf("get cache articleId by sortType error: err is %s, sortType is %v, cursor is %v", err.Error(), sortType, cursor)
		return nil, err
	}
	var articleIds []int64
	for _, pair := range pairs {
		id, _ := strconv.ParseInt(pair.Key, 10, 64)
		articleIds = append(articleIds, id)
	}
	return articleIds, nil
}

func (l *ArticlesLogic) articleByIds(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](func(source chan<- int64) {
		for _, aid := range articleIds {
			source <- aid
		}
	}, func(id int64, writer mr.Writer[*model.Article], cancel func(error)) {
		p, err := l.svcCtx.ArticleModel.FindOne(ctx, uint64(id))
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) {
		var articles []*model.Article
		for article := range pipe {
			articles = append(articles, article)
		}
		writer.Write(articles)
	})
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (l *ArticlesLogic) addCacheArticles(cur *pb.ArticleItem, uId, publishTime, likeCount int64, sortType int32) error {
	var (
		key   string
		score int64
	)
	if sortType == types.PublishTimeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticlePublishTimePrefix, uId)
		score = publishTime
	} else if sortType == types.LikeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticleLikePrefix, uId)
		score = likeCount
	}

	_, err := l.svcCtx.BizRedis.Zadd(key, score, strconv.Itoa(int(cur.Id)))
	if err != nil {
		return err
	}
	return l.svcCtx.BizRedis.ExpireCtx(context.Background(), key, articlesExpire)
}

func articlesKey(uid int64, sortType int32) (key string) {
	if sortType == types.PublishTimeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticlePublishTimePrefix, uid)
	} else if sortType == types.LikeSortType {
		key = fmt.Sprintf("%s%d", types.CacheArticleLikePrefix, uid)
	}
	return
}
