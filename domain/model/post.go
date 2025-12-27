package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialPlatform string

const (
	PlatformInstagram SocialPlatform = "instagram"
	PlatformTwitter   SocialPlatform = "twitter"
	PlatformYouTube   SocialPlatform = "youtube"
	PlatformTikTok    SocialPlatform = "tiktok"
	PlatformUnknown   SocialPlatform = "unknown"
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
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    *uuid.UUID     `json:"userId" gorm:"type:uuid"`
	UserIP    string         `json:"userIp" gorm:"type:inet"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	PostDate  time.Time      `json:"postDate"`

	Link     string         `json:"link"`
	ImageURL *string        `json:"imageUrl"`
	VideoURL *string        `json:"videoUrl,omitempty"` // optional field
	Platform SocialPlatform `json:"platform"`
	Format   PostFormat     `json:"format"`

	Status     PostStatus `json:"status"`
	FailReason *string    `json:"failReason"`

	UserName         string `json:"userName"`
	UserAnchor       string `json:"userAnchor"`
	UserProfileLink  string `json:"userProfileLink"`
	UserProfileImage string `json:"userProfileImage"`

	User     *User          `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Contents []*PostContent `json:"contents,omitempty" gorm:"foreignKey:PostID"`
}

func (*Post) TableName() string {
	return "posts"
}
