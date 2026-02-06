package middleware

import (
	"github.com/amahdian/cliplab-be/svc/auth"
	"github.com/gin-gonic/gin"
)

func WithAuth(authenticator auth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		ctx, err := authenticator.Verify(r)
		if err == nil {
			c.Request = r.WithContext(ctx)
		}
		c.Next()
	}
}
