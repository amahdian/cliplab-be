package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ChannelID uuid.UUID `json:"channelId" gorm:"type:uuid"`
	CreatedAt time.Time `json:"createdAt"`

	FollowersCount int `json:"followersCount"`
	FollowingCount int `json:"followingCount"`
	MediaCount     int `json:"mediaCount"`

	AverageLikes      int `json:"averageLikes"`
	AverageComments   int `json:"averageComments"`
	AverageVideoViews int `json:"averageVideoViews"`
	AverageVideoPlays int `json:"averageVideoPlays"`
}

func (*ChannelHistory) TableName() string {
	return "channel_histories"
}
