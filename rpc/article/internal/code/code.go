package code

import "pkg/xerror"

var (
	UserIdInvaild        = xerror.New(300000, "用户id不合法")
	ArticleTitleEmpty    = xerror.New(300001, "文章标题不能为空")
	ArticleContentSmall  = xerror.New(300002, "文章字数不足100")
	ArticleSortTypeError = xerror.New(400001, "文章排序方式错误")
	ArticlePageSizeBig   = xerror.New(400002, "页大小超过200")
	ArticlePageSizeError = xerror.New(400003, "页大小不合法")
	ArticleIdInvaild     = xerror.New(400004, "文章id不合法")
)
