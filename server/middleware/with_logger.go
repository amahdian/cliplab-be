package middleware

import (
	"fmt"
	"time"

	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/gin-gonic/gin"
)

func WithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			return
		}

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		logger.Infof(defaultLogFormatter(param))
	}
}

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[REQUEST] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
