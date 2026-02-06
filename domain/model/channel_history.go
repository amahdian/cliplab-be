package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelHistory struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ChannelID uuid.UUID `json:"channelId" gorm:"type:uuid"`
	CreatedAt time.Time `json:"createdAt"`

	FollowersCount    int64  `json:"followersCount"`
	FollowingCount    int64  `json:"followingCount"`
	MediaCount        int64  `json:"mediaCount"`
	ProfileImage      string `json:"profileImage"`
	ProfileDescriptor string `json:"profileDescriptor"`

	AverageLikes          float64   `json:"averageLikes"`
	AverageComments       float64   `json:"averageComments"`
	AverageVideoViews     float64   `json:"averageVideoViews"`
	AverageVideoPlays     float64   `json:"averageVideoPlays"`
	AverageEngagementRate float64   `json:"averageEngagementRate"`
	LatestPostPublishDate time.Time `json:"latestPostPublishDate"`
}

func (*ChannelHistory) TableName() string {
	return "channel_histories"
}
