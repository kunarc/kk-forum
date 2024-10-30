// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.2

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	articleFieldNames          = builder.RawFieldNames(&Article{})
	articleRows                = strings.Join(articleFieldNames, ",")
	articleRowsExpectAutoSet   = strings.Join(stringx.Remove(articleFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	articleRowsWithPlaceHolder = strings.Join(stringx.Remove(articleFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheKkArticleArticleIdPrefix = "cache:kkArticle:article:id:"
)

type (
	articleModel interface {
		Insert(ctx context.Context, data *Article) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Article, error)
		Update(ctx context.Context, data *Article) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultArticleModel struct {
		sqlc.CachedConn
		table string
	}

	Article struct {
		Id          uint64    `db:"id"`           // 主键ID
		Title       string    `db:"title"`        // 标题
		Content     string    `db:"content"`      // 内容
		Cover       string    `db:"cover"`        // 封面
		Description string    `db:"description"`  // 描述
		AuthorId    uint64    `db:"author_id"`    // 作者ID
		Status      int64     `db:"status"`       // 状态 0:待审核 1:审核不通过 2:可见 3:用户删除
		CommentNum  int64     `db:"comment_num"`  // 评论数
		LikeNum     int64     `db:"like_num"`     // 点赞数
		CollectNum  int64     `db:"collect_num"`  // 收藏数
		ViewNum     int64     `db:"view_num"`     // 浏览数
		ShareNum    int64     `db:"share_num"`    // 分享数
		TagIds      string    `db:"tag_ids"`      // 标签ID
		PublishTime time.Time `db:"publish_time"` // 发布时间
		CreateTime  time.Time `db:"create_time"`  // 创建时间
		UpdateTime  time.Time `db:"update_time"`  // 最后修改时间
	}
)

func newArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultArticleModel {
	return &defaultArticleModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`article`",
	}
}

func (m *defaultArticleModel) Delete(ctx context.Context, id uint64) error {
	kkArticleArticleIdKey := fmt.Sprintf("%s%v", cacheKkArticleArticleIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, kkArticleArticleIdKey)
	return err
}

func (m *defaultArticleModel) FindOne(ctx context.Context, id uint64) (*Article, error) {
	kkArticleArticleIdKey := fmt.Sprintf("%s%v", cacheKkArticleArticleIdPrefix, id)
	var resp Article
	err := m.QueryRowCtx(ctx, &resp, kkArticleArticleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", articleRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultArticleModel) Insert(ctx context.Context, data *Article) (sql.Result, error) {
	kkArticleArticleIdKey := fmt.Sprintf("%s%v", cacheKkArticleArticleIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, articleRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Title, data.Content, data.Cover, data.Description, data.AuthorId, data.Status, data.CommentNum, data.LikeNum, data.CollectNum, data.ViewNum, data.ShareNum, data.TagIds, data.PublishTime)
	}, kkArticleArticleIdKey)
	return ret, err
}

func (m *defaultArticleModel) Update(ctx context.Context, data *Article) error {
	kkArticleArticleIdKey := fmt.Sprintf("%s%v", cacheKkArticleArticleIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, articleRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Title, data.Content, data.Cover, data.Description, data.AuthorId, data.Status, data.CommentNum, data.LikeNum, data.CollectNum, data.ViewNum, data.ShareNum, data.TagIds, data.PublishTime, data.Id)
	}, kkArticleArticleIdKey)
	return err
}

func (m *defaultArticleModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheKkArticleArticleIdPrefix, primary)
}

func (m *defaultArticleModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", articleRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultArticleModel) tableName() string {
	return m.table
}
