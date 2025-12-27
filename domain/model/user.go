package model

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderLocal    Provider = "local"
	ProviderGoogle   Provider = "google"
	ProviderFacebook Provider = "facebook"
)

type User struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Email           string    `json:"email"`
	Password        *string   `json:"-"`
	Name            *string   `json:"name"`
	Provider        Provider  `json:"-"`
	ProviderID      *string   `json:"-"` // User's ID from the provider
	ProfileImage    *string   `json:"profileImage"`
	DefaultLanguage Language  `json:"defaultLanguage"`

	VerifiedAt *time.Time `json:"verifiedAt,omitempty"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
}

func (*User) TableName() string {
	return "users"
}
