package req

import (
	"context"
	"net"

	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/gin-gonic/gin"
)

type RequestContext struct {
	Ctx      context.Context
	UserInfo *auth.UserInfo
	Ip       net.IP
}

func GetRequestContext(c *gin.Context) RequestContext {
	ctx := c.Request.Context()
	userInfo := auth.UserInfoFromCtx(ctx)

	return RequestContext{
		Ctx:      ctx,
		UserInfo: &userInfo,
		Ip:       net.ParseIP(c.ClientIP()),
	}
}
