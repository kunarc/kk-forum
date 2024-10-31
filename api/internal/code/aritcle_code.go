package code

import "pkg/xerror"

var (
	ArticleTitleEmpty   = xerror.New(300001, "文章标题不能为空")
	ArticleContentSmall = xerror.New(300002, "文章字数不足100")
	ArticleCoverBig     = xerror.New(300003, "文件封面太大")
)
