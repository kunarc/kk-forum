package model

import (
	"context"
	"fmt"
	"time"

	"article/internal/types"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		FindIdsByCursor(ctx context.Context, uId, cursor, ps int64, sortType int32) ([]*Article, error)
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn, c, opts...),
	}
}

func (m *customArticleModel) FindIdsByCursor(ctx context.Context, uId, cursor, ps int64, sortType int32) ([]*Article, error) {
	var (
		sql         string
		articles    []*Article
		publishTime = time.Unix(cursor, 0).Format("2006-01-02 15:04:05")
		likeNum     = cursor
	)
	if sortType == types.PublishTimeSortType {
		sql = fmt.Sprintf("select "+articleRows+" from"+m.table+" where author_id = %v and publish_time <= '%v' order by publish_time desc limit %v", uId, publishTime, ps)
	} else {
		sql = fmt.Sprintf("select "+articleRows+" from"+m.table+" where author_id = %v and like_num <= %v order by like_num desc limit %v", uId, likeNum, ps)
	}
	err := m.QueryRowsNoCacheCtx(ctx, &articles, sql)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
