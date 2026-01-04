package router

import (
	"net/http"
	"time"

	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/amahdian/cliplab-be/version"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ginSwaggerHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)

func (r *Router) healthCheck(ctx *gin.Context) {
	response := resp.HealthResponseDto{
		AppName:    version.AppName,
		AppVersion: version.AppVersion,
	}

	resp.Ok(ctx, response)
}

func (r *Router) swaggerHandler(ctx *gin.Context) {
	// the ginSwaggerHandler by default recognizes the "/swagger/index.html"  but not"/swagger" or "/swagger/".
	// therefore we add support for these endpoints by redirecting to "/swagger/index.html"
	if ctx.Request.RequestURI == "/api/swagger" || ctx.Request.RequestURI == "/api/swagger/" {
		ctx.Redirect(http.StatusFound, "/api/swagger/index.html")
	}
	ginSwaggerHandler(ctx)
}

func (r *Router) getServerTime(ctx *gin.Context) {
	resp.Ok(ctx, time.Now())
}
