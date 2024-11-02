package logic

import (
	"context"
	"encoding/json"

	"like/internal/svc"
	"like/internal/types"
	"like/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *pb.ThumbupRequest) (*pb.ThumbupResponse, error) {
	thumbupMsg := &types.ThumbupMqMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}
	threading.GoSafe(func() {
		msg, err := json.Marshal(thumbupMsg)
		if err != nil {
			l.Logger.Errorf("[Thumbup] marshal msg error: msg is %v, err is %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(msg))
		if err != nil {
			l.Logger.Errorf("[Thumbup] push msg error: msg is %v, err is %v", msg, err)
			return
		}
	})
	return &pb.ThumbupResponse{}, nil
}
