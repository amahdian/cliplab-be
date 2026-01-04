package req

import "github.com/amahdian/cliplab-be/domain/model"

type Register struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Name     *string `json:"name" binding:"omitempty,min=2,max=100"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OauthLogin struct {
	Provider model.Provider `json:"provider" binding:"required,oneof=google facebook"`
	Token    string         `json:"token" binding:"required"`
}

type Verify struct {
	Email string `json:"email" binding:"required,mobile"`
	Otp   string `json:"otp" binding:"required"`
}

type UserUpdate struct {
	Name            *string         `json:"name" binding:"omitempty,min=2,max=100"`
	DefaultLanguage *model.Language `json:"defaultLanguage" binding:"omitempty,oneof=en fa"`
}
