package model

import (
	"time"

	"github.com/google/uuid"
)

type UsageStatus string

const (
	UsageStatusProcessing UsageStatus = "processing"
	UsageStatusCompleted  UsageStatus = "completed"
	UsageStatusFailed     UsageStatus = "failed"
)

type Tool string

const (
	ToolVideoAnalysis     Tool = "video_analysis"
	ToolChannelEngagement Tool = "channel_engagement"
)

func (t Tool) DisplayName() string {
	switch t {
	case ToolVideoAnalysis:
		return "Video Analysis"
	case ToolChannelEngagement:
		return "Channel Engagement"
	default:
		return string(t)
	}
}

type UsageHistory struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId           uuid.UUID      `json:"userId" gorm:"type:uuid"`
	ToolName         Tool           `json:"toolName"` // "video_analysis", "channel_engagement", etc.
	CreditsUsed      int            `json:"creditsUsed"`
	RemainingCredits int            `json:"remainingCredits"`
	Status           UsageStatus    `json:"status"`                     // "completed", "failed", "processing"
	Details          string         `json:"details"`                    // e.g., the URL or handle
	Platform         SocialPlatform `json:"platform"`                   // "instagram", "tiktok", "youtube", "general"
	Type             string         `json:"type"`                       // "analysis", "engagement"
	RequestID        *uuid.UUID     `json:"requestId" gorm:"type:uuid"` // Link to AnalyzeRequest ID if applicable
	ResponseLog      string         `json:"responseLog"`                // JSON log of the response
	CreatedAt        time.Time      `json:"createdAt"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserId"`
}

func (*UsageHistory) TableName() string {
	return "usage_histories"
}
