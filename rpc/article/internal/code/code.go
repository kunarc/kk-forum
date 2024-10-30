package code

import "pkg/xerror"

var (
	UserIdInvaild       = xerror.New(30000, "用户id不合法")
	ArticleTitleEmpty   = xerror.New(300001, "文章标题不能为空")
	ArticleContentSmall = xerror.New(300002, "文章字数不足100")
)
