package svc

import (
	"context"
	"errors"
	"time"

	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/amahdian/cliplab-be/svc/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserSvc interface {
	Register(data *req.Register) (*resp.AuthResponse, error)
	Login(data *req.Login) (*resp.AuthResponse, error)
	LoginOauth(data *req.OauthLogin) (*resp.AuthResponse, error)
	Verify(email, otp string) (string, *model.User, error)
	Update(userID uuid.UUID, updateData *req.UserUpdate) error
	Me(userInfo *auth.UserInfo) (*model.User, error)
}

type userSvc struct {
	ctx  context.Context
	stg  storage.PgStorage
	envs *env.Envs
}

func newUserSvc(ctx context.Context, stg storage.PgStorage, envs *env.Envs) UserSvc {
	return &userSvc{
		ctx:  ctx,
		stg:  stg,
		envs: envs,
	}
}

func (s *userSvc) Register(data *req.Register) (*resp.AuthResponse, error) {
	existingUser, err := s.stg.User(s.ctx).FindByEmail(data.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:           data.Email,
		Password:        &hashedPassword,
		Name:            data.Name,
		Provider:        model.ProviderLocal,
		DefaultLanguage: model.LanguageEnglish,
	}

	err = s.stg.User(s.ctx).CreateOne(user)
	if err != nil {
		return nil, err
	}

	return s.generateAuthResponse(user)
}

func (s *userSvc) Login(data *req.Login) (*resp.AuthResponse, error) {
	user, err := s.stg.User(s.ctx).FindByEmail(data.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Password == nil {
		return nil, errors.New("invalid credentials")
	}

	if user.VerifiedAt == nil && time.Since(user.CreatedAt) > 14*24*time.Hour {
		return nil, errors.New("email not verified")
	}

	if !utils.CheckPasswordHash(data.Password, *user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

func (s *userSvc) LoginOauth(data *req.OauthLogin) (*resp.AuthResponse, error) {
	// This is a placeholder. In a real application, you would validate
	// the token with the OAuth provider and get the user's info.
	// For now, we'll just simulate it.

	// Example: In a real scenario you would call a function like:
	// userInfo, err := oauth.VerifyToken(data.Provider, data.Token)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// email := userInfo.Email
	// providerID := userInfo.ID
	// name := userInfo.Name

	// Simulated user info for demonstration
	email := "user@example.com"
	providerID := "123456789"
	name := "OAuth User"

	user, err := s.stg.User(s.ctx).FindByProvider(data.Provider, providerID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// User doesn't exist, create a new one
		now := time.Now()
		user = &model.User{
			Email:      email,
			Name:       &name,
			Provider:   data.Provider,
			ProviderID: &providerID,
			VerifiedAt: &now, // OAuth users are considered verified
		}
		err = s.stg.User(s.ctx).CreateOne(user)
		if err != nil {
			return nil, err
		}
	}

	return s.generateAuthResponse(user)
}

func (s *userSvc) Verify(email, otp string) (string, *model.User, error) {
	user, err := s.stg.User(s.ctx).FindByEmail(email)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, errors.New("user not found")
	}

	// Logic to verify OTP would go here.
	// For now, we'll assume it's always successful.

	now := time.Now()
	user.VerifiedAt = &now
	err = s.stg.User(s.ctx).UpdateOne(user, false)
	if err != nil {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().AddDate(1, 0, 0).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(s.envs.Server.JwtSecret))

	return tokenStr, user, nil
}

func (s *userSvc) Update(userID uuid.UUID, updateData *req.UserUpdate) error {
	user, err := s.stg.User(s.ctx).FindById(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}
	if user.ID != userID {
		return errors.New("unauthorized")
	}

	// Update fields if provided
	if updateData.Name != nil {
		user.Name = updateData.Name
	}
	if updateData.DefaultLanguage != nil {
		user.DefaultLanguage = *updateData.DefaultLanguage
	}

	return s.stg.User(s.ctx).UpdateOne(user, false)
}

func (s *userSvc) Me(userInfo *auth.UserInfo) (*model.User, error) {
	user, err := s.stg.User(s.ctx).FindById(userInfo.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *userSvc) generateAuthResponse(user *model.User) (*resp.AuthResponse, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().AddDate(1, 0, 0).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.envs.Server.JwtSecret))
	if err != nil {
		return nil, err
	}

	return &resp.AuthResponse{
		Token: tokenStr,
		User:  user,
	}, nil
}
