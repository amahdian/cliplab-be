package resp

import (
	"github.com/amahdian/cliplab-be/domain/model"
)

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type VerifyResponse struct {
	Email           string         `json:"email"`
	Name            *string        `json:"name"`
	DefaultLanguage model.Language `json:"defaultLanguage"`
	Token           string         `json:"token"`
}
