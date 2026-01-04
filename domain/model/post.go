package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostFormat string

const (
	PostFormatImage PostFormat = "image"
	PostFormatVideo PostFormat = "video"
	PostFormatText  PostFormat = "text"
	PostFormatSound PostFormat = "sound"
)

type PostStatus string

const (
	PostStatusPending    PostStatus = "pending"
	PostStatusProcessing PostStatus = "processing"
	PostStatusCompleted  PostStatus = "completed"
	PostStatusFailed     PostStatus = "failed"
)

type Post struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	UserId    *uuid.UUID     `json:"userId" gorm:"type:uuid"`
	UserIP    string         `json:"userIp" gorm:"type:inet"`
	ChannelId *uuid.UUID     `json:"channelId" gorm:"type:uuid"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
	PostDate  time.Time      `json:"postDate"`

	Link     string     `json:"link"`
	ImageURL *string    `json:"imageUrl"`
	VideoURL *string    `json:"videoUrl,omitempty"` // optional field
	Format   PostFormat `json:"format"`

	Status     PostStatus `json:"status"`
	FailReason *string    `json:"failReason"`

	UserName         string `json:"userName"`
	UserAnchor       string `json:"userAnchor"`
	UserProfileLink  string `json:"userProfileLink"`
	UserProfileImage string `json:"userProfileImage"`

	LikeCount      int64 `json:"likeCount"`
	CommentCount   int64 `json:"commentCount"`
	VideoViewCount int64 `json:"videoViewCount"`
	VideoPlayCount int64 `json:"videoPlayCount"`

	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Channel *Channel `json:"channel,omitempty" gorm:"foreignKey:ChannelId;references:ID"`
}

func (*Post) TableName() string {
	return "posts"
}
