package logic

import (
	"context"
	"strconv"
	"time"

	"follow/code"
	"follow/internal/model"
	"follow/internal/svc"
	"follow/internal/types"
	"follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消关注
func (l *UnFollowLogic) UnFollow(in *pb.UnFollowRequest) (*pb.UnFollowResponse, error) {
	if in.UserId == 0 {
		return nil, code.CancelUserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.CancelFollowedUserIdEmpty
	}
	if in.UserId == in.FollowedUserId {
		return nil, code.CannotCancelSelf
	}
	follow, err := l.svcCtx.FollowModel.FindByUserIDAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[Follow] FollowModel.FindByUserIDAndFollowedUserID err: %v req: %v", err, in)
		return nil, err
	}
	if follow == nil || follow.FollowStatus == types.FollowStatusUnfollow {
		return nil, code.CancelObjectInvaild
	}
	followTime := time.Now()
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 更新follow状态
		err = model.NewFollowModel(tx).UpdateFollowStatus(l.ctx, followTime, follow.ID, types.FollowStatusUnfollow)
		if err != nil {
			l.Logger.Errorf("[UpdateFollowStatus] error: err is %v, id is %v", err.Error(), follow.ID)
		}

		// 更新 follow_count 表
		err = model.NewFollowCountModel(tx).ReduceFollowCount(l.ctx, in.UserId, followTime)
		if err != nil {
			l.Logger.Errorf("[ReduceFollowCount] error: err is %v, userId is %v", err.Error(), in.UserId)
		}
		err = model.NewFollowCountModel(tx).ReduceFansCount(l.ctx, in.FollowedUserId, followTime)
		if err != nil {
			l.Logger.Errorf("[ReduceFansCount] error: err is %v, userId is %v", err.Error(), in.FollowedUserId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	followExist, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, cacheFollowKey(types.CacheFollowType, in.UserId))
	if followExist {
		_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, cacheFollowKey(types.CacheFollowType, in.UserId), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[ZremCtx] del cache follow key error: err is %v, key is %v", err.Error(), cacheFollowKey(types.CacheFollowType, in.UserId))
		}
	}
	fansExist, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, cacheFollowKey(types.CacheFansType, in.UserId))
	if fansExist {
		_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, cacheFollowKey(types.CacheFollowType, in.FollowedUserId), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[ZremCtx] del cache fans key error: err is %v, key is %v", err.Error(), cacheFollowKey(types.CacheFansType, in.FollowedUserId))
		}
	}
	return &pb.UnFollowResponse{}, nil
}
