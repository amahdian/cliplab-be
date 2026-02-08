package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID               uuid.UUID `json:"userId" gorm:"index"`
	PaddleSubscriptionID string    `json:"paddleSubscriptionId" gorm:"uniqueIndex"`
	PriceID              string    `json:"priceId"`
	Status               string    `json:"status"`
	CurrentPeriodStart   time.Time `json:"currentPeriodStart"`
	CurrentPeriodEnd     time.Time `json:"currentPeriodEnd"`
	CancelAtPeriodEnd    bool      `json:"cancelAtPeriodEnd"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

func (*Subscription) TableName() string {
	return "subscriptions"
}
