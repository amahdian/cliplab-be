package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type userInfoCtx struct{}

type UserInfo struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (u *UserInfo) User() model.User {
	return model.User{
		ID:    u.Id,
		Email: u.Email,
		Name:  lo.ToPtr(u.Name),
	}
}

type Authenticator interface {
	Verify(request *http.Request) (context.Context, error)
}

type authenticator struct {
	JwtSecret string
	Stg       storage.PgStorage
}

func NewAuthenticator(envs *env.Envs, stg storage.PgStorage) Authenticator {
	return &authenticator{
		JwtSecret: envs.Server.JwtSecret,
		Stg:       stg,
	}
}

func (a *authenticator) Verify(request *http.Request) (context.Context, error) {
	ctx := request.Context()
	tokenStr := request.Header.Get("Authorization")

	if tokenStr == "" {
		tokenStr = request.URL.Query().Get("authToken")
	}

	if tokenStr == "" {
		return ctx, errors.New("authorization header is empty")
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			return ctx, errs.Newf(errs.Unauthenticated, err, "Invalid user ID.")
		}

		user, err := a.Stg.User(ctx).FindById(userID)
		if err != nil {
			return ctx, errs.Newf(errs.Unauthenticated, err, "User not found.")
		}

		if user.VerifiedAt == nil && time.Since(user.CreatedAt) > 14*24*time.Hour {
			return ctx, errs.Newf(errs.Unauthenticated, nil, "Email not verified.")
		}

		userInfo := UserInfo{
			Id:    userID,
			Email: claims["email"].(string),
		}
		if claims["name"] != nil {
			userInfo.Name = claims["name"].(string)
		}

		ctx = context.WithValue(ctx, userInfoCtx{}, userInfo)

		return ctx, nil
	} else {
		return ctx, errs.Newf(errs.Unauthenticated, err, "Auth failed.")
	}
}

func UserInfoFromCtx(ctx context.Context) UserInfo {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	u, ok := ctx.Value(userInfoCtx{}).(UserInfo)
	if !ok {
		return UserInfo{}
	}
	return u
}
