package req

import (
	"context"
	"net"

	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestContext struct {
	Ctx      context.Context
	UserInfo *auth.UserInfo
	Ip       net.IP
}

func GetRequestContext(c *gin.Context) RequestContext {
	ctx := c.Request.Context()
	userInfo := auth.UserInfoFromCtx(ctx)
	var userInfoPtr *auth.UserInfo
	if userInfo.Id != uuid.Nil {
		userInfoPtr = &userInfo
	}

	return RequestContext{
		Ctx:      ctx,
		UserInfo: userInfoPtr,
		Ip:       net.ParseIP(c.ClientIP()),
	}
}
