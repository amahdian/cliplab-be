package svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	var email, providerID, name, picture string

	if data.Provider == model.ProviderGoogle {
		googleInfo, err := s.verifyGoogleToken(data.Token)
		if err != nil {
			return nil, err
		}
		email = googleInfo.Email
		providerID = googleInfo.Sub
		name = googleInfo.Name
		picture = googleInfo.Picture
	} else {
		return nil, errors.New("provider not supported yet")
	}

	user, err := s.stg.User(s.ctx).FindByProvider(data.Provider, providerID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Try to find by email if provider ID not found
		user, err = s.stg.User(s.ctx).FindByEmail(email)
		if err != nil {
			return nil, err
		}

		if user != nil {
			// Link provider to existing account
			user.ProviderID = &providerID
			user.Provider = data.Provider
			if user.Name == nil {
				user.Name = &name
			}
			if user.ProfileImage == nil && picture != "" {
				user.ProfileImage = &picture
			}
			now := time.Now()
			if user.VerifiedAt == nil {
				user.VerifiedAt = &now
			}
			err = s.stg.User(s.ctx).UpdateOne(user, false)
			if err != nil {
				return nil, err
			}
		} else {
			// Create new user
			now := time.Now()
			user = &model.User{
				Email:        email,
				Name:         &name,
				Provider:     data.Provider,
				ProviderID:   &providerID,
				ProfileImage: &picture,
				VerifiedAt:   &now,
			}
			err = s.stg.User(s.ctx).CreateOne(user)
			if err != nil {
				return nil, err
			}
		}
	}

	return s.generateAuthResponse(user)
}

type googleTokenInfo struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	Error         string `json:"error"`
}

func (s *userSvc) verifyGoogleToken(token string) (*googleTokenInfo, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// We use the userinfo endpoint which accepts both access tokens (via header or query) and provides user details
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token verification failed with status: %d", resp.StatusCode)
	}

	var info googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	if info.Error != "" {
		return nil, fmt.Errorf("google token verification failed: %s", info.Error)
	}

	return &info, nil
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
