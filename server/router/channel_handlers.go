package router

import (
	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
)

func (r *Router) getChannelEngagementRate(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	var request req.EngagementRateRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	channelSvc := r.svc.NewChannelSvc(reqCtx.Ctx)
	res, err := channelSvc.GetChannelEngagement(request.URL, request.Platform, reqCtx.UserInfo)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, res)
}
