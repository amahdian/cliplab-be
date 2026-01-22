package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	RequestStatusPending    RequestStatus = "pending"
	RequestStatusProcessing RequestStatus = "processing"
	RequestStatusCompleted  RequestStatus = "completed"
	RequestStatusFailed     RequestStatus = "failed"
)

type AnalyzeRequest struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    *uuid.UUID `json:"userId" gorm:"type:uuid"`
	UserIP    string     `json:"userIp" gorm:"type:inet"`
	Link      string     `json:"link"`
	PostId    *string    `json:"postId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`

	Status     RequestStatus `json:"status"`
	FailReason *string       `json:"failReason"`

	// for debug
	LlmRequest  string `json:"llmRequest"`
	LlmResponse string `json:"llmResponse"`

	Post *Post `json:"post,omitempty" gorm:"foreignKey:PostId;references:ID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
}

func (*AnalyzeRequest) TableName() string {
	return "analyze_requests"
}
