package model

import (
	"context"
	"time"

	"follow/internal/types"
	"follow/pb"

	"gorm.io/gorm"
)

type Follow struct {
	ID             int64 `gorm:"primary_key"`
	UserID         int64
	FollowedUserID int64
	FollowStatus   int
	CreateTime     time.Time
	UpdateTime     time.Time
}

func (m *Follow) TableName() string {
	return "follow"
}

type FollowModle struct {
	db *gorm.DB
}

func NewFollowModel(db *gorm.DB) *FollowModle {
	return &FollowModle{db: db}
}

func (m *FollowModle) FindByUserIDAndFollowedUserID(ctx context.Context, uID, fuId int64) (*Follow, error) {
	var follow Follow
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND followed_user_id = ?", uID, fuId).
		First(&follow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &follow, err
}

func (m *FollowModle) InsertFollowRecord(ctx context.Context, follow *Follow) error {
	err := m.db.WithContext(ctx).Create(follow).Error
	return err
}

func (m *FollowModle) UpdateFollowStatus(ctx context.Context, followTime time.Time, id int64, status int) error {
	err := m.db.WithContext(ctx).Model(&Follow{}).
		Where("id = ?", id).Updates(map[string]any{
		"follow_status": status,
		"update_time":   followTime,
	}).Error
	return err
}

func (m *FollowModle) GetFollowItemList(ctx context.Context, uID, cursor, ps int64) (res []*pb.FollowItem, err error) {
	data := []*struct {
		pb.FollowItem
		CreateTime time.Time
	}{}
	err = m.db.WithContext(ctx).Model(&Follow{}).Where("follow.user_id = ? AND follow.create_time < ? AND follow_status = ?", uID, cursor, types.FollowStatusFollow).
		Joins("JOIN follow_count ON follow.followed_user_id = follow_count.user_id").
		Select("follow.id", "follow.followed_user_id", "fans_count", "follow.create_time").Limit(int(ps)).Order("create_time desc").Find(&data).Error
	for i := 0; i < len(data); i++ {
		res = append(res, &pb.FollowItem{
			Id:             data[i].Id,
			FollowedUserId: data[i].FollowedUserId,
			FansCount:      data[i].FansCount,
			CreateTime:     data[i].CreateTime.Unix(),
		})
	}
	return
}

func (m *FollowModle) FindOne(ctx context.Context, uid, fid int64) (res *Follow, err error) {
	err = m.db.WithContext(ctx).Model(&Follow{}).Where("user_id = ? AND followed_user_id = ?", uid, fid).Find(&res).Error
	return
}
