package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ChannelID uuid.UUID `json:"channelId" gorm:"type:uuid"`
	CreatedAt time.Time `json:"createdAt"`

	FollowersCount int64 `json:"followersCount"`
	FollowingCount int64 `json:"followingCount"`
	MediaCount     int64 `json:"mediaCount"`

	AverageLikes      int64 `json:"averageLikes"`
	AverageComments   int64 `json:"averageComments"`
	AverageVideoViews int64 `json:"averageVideoViews"`
	AverageVideoPlays int64 `json:"averageVideoPlays"`
}

func (*ChannelHistory) TableName() string {
	return "channel_histories"
}
