package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type FollowCount struct {
	ID          int64 `gorm:"primary_key"`
	UserID      int64
	FollowCount int
	FansCount   int
	CreateTime  time.Time
	UpdateTime  time.Time
}

func (m *FollowCount) TableName() string {
	return "follow_count"
}

type FollowCountModel struct {
	db *gorm.DB
}

func NewFollowCountModel(db *gorm.DB) *FollowCountModel {
	return &FollowCountModel{
		db: db,
	}
}

func (m *FollowCountModel) InsertFollowCount(ctx context.Context, uID int64, followedID int64, followTime time.Time) error {
	followsCount := make([]FollowCount, 2)
	followsCount[0] = FollowCount{
		UserID:     uID,
		CreateTime: followTime,
		UpdateTime: followTime,
	}
	followsCount[1] = FollowCount{
		UserID:     followedID,
		CreateTime: followTime,
		UpdateTime: followTime,
	}
	return m.db.WithContext(ctx).Save(followsCount).Error
}

func (m *FollowCountModel) IncrFollowCount(ctx context.Context, uID int64, followTime time.Time) error {
	err := m.db.WithContext(ctx).Exec(
		"UPDATE follow_count SET follow_count = follow_count + 1, update_time = ? WHERE user_id = ?",
		followTime, uID,
	).Error
	return err
}

func (m *FollowCountModel) ReduceFollowCount(ctx context.Context, uID int64, followTime time.Time) error {
	err := m.db.WithContext(ctx).Exec(
		"UPDATE follow_count SET follow_count = follow_count - 1, update_time = ? WHERE user_id = ?",
		followTime, uID,
	).Error
	return err
}

func (m *FollowCountModel) IncrFansCount(ctx context.Context, uID int64, followTime time.Time) error {
	err := m.db.WithContext(ctx).Exec(
		"UPDATE follow_count SET  fans_count = fans_count + 1, update_time = ? WHERE user_id = ?",
		followTime, uID,
	).Error
	return err
}

func (m *FollowCountModel) ReduceFansCount(ctx context.Context, uID int64, followTime time.Time) error {
	err := m.db.WithContext(ctx).Exec(
		"UPDATE follow_count SET  fans_count = fans_count - 1, update_time = ? WHERE user_id = ?",
		followTime, uID,
	).Error
	return err
}

func (m *FollowCountModel) FindOne(ctx context.Context, uID int64) (fc *FollowCount, err error) {
	err = m.db.WithContext(ctx).Where("user_id = ?", uID).Find(&fc).Error
	return
}
