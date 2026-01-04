package router

import (
	"net/url"

	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/gin-gonic/gin"
)

func (r *Router) addPostToAnalyzeQueue(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo

	request := &req.AddPost{}
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	link, err := url.Parse(request.Url)
	if err != nil {
		resp.AbortWithError(ctx, errs.Newf(errs.InvalidArgument, nil, "invalid link: %v", err))
		return
	}

	dSvc := r.svc.NewPostSvc(reqCtx.Ctx)
	id, err := dSvc.AddPostToAnalyzeQueue(*link, user, reqCtx.Ip)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, id)
}

func (r *Router) getPostData(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.IdUri{}
	if err := ctx.BindUri(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewPostSvc(reqCtx.Ctx)
	post, err := dSvc.GetPostById(request.Id)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, post)
}
