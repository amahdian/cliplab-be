package model

import (
	"time"

	"github.com/google/uuid"
)

type SocialPlatform string

const (
	PlatformInstagram SocialPlatform = "instagram"
	PlatformTwitter   SocialPlatform = "twitter"
	PlatformYouTube   SocialPlatform = "youtube"
	PlatformTikTok    SocialPlatform = "tiktok"
	PlatformUnknown   SocialPlatform = "unknown"
)

type Channel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	FullName  string         `json:"fullName"`
	Handler   string         `json:"handler"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Platform  SocialPlatform `json:"platform"`

	Histories   []*ChannelHistory `json:"history"`
	LastHistory *ChannelHistory   `json:"lastHistory" gorm:"-"`
}

func (*Channel) TableName() string {
	return "channels"
}
