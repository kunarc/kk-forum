package logic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(r *http.Request) (resp *types.UploadCoverResponse, err error) {
	// todo: add your logic here and delete this line
	err = r.ParseMultipartForm(maxFileSize)
	if err != nil {
		return nil, err
	}
	file, header, err := r.FormFile("cover")
	if err != nil {
		l.Logger.Errorf("read cover from req error: err is %s", err.Error())
	}
	defer file.Close()
	bucket, err := l.svcCtx.Oss.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		l.Logger.Errorf("get oss bucket error: err is %s", err.Error())
		return nil, err
	}
	fileName := header.Filename
	objectKey := l.genObjectKey(fileName)
	err = bucket.PutObject(objectKey, file)
	if err != nil {
		l.Logger.Errorf("put object to bucket error: err is %s", err.Error())
		return nil, err
	}
	return &types.UploadCoverResponse{CoverUrl: l.genFileURL(objectKey)}, nil
}

func (l *UploadCoverLogic) genObjectKey(fileName string) string {
	return fmt.Sprintf("%d_%s", time.Now().Unix(), fileName)
}

func (l *UploadCoverLogic) genFileURL(objectKey string) string {
	return fmt.Sprintf("https://%s.%s/%s", l.svcCtx.Config.Oss.BucketName, l.svcCtx.Config.Oss.Endpoint, objectKey)
}
