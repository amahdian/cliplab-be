package router

import (
	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
)

func (r *Router) register(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.Register{}
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	authResp, err := dSvc.Register(request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, authResp)
}

func (r *Router) login(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.Login{}
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	authResp, err := dSvc.Login(request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, authResp)
}

func (r *Router) loginOauth(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.OauthLogin{}
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	authResp, err := dSvc.LoginOauth(request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, authResp)
}

func (r *Router) verify(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.Verify{}
	err := ctx.BindJSON(&request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	token, user, err := dSvc.Verify(request.Email, request.Otp)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	res := &resp.VerifyResponse{
		Token:           token,
		Name:            user.Name,
		Email:           user.Email,
		DefaultLanguage: user.DefaultLanguage,
	}

	resp.Ok(ctx, res)
}

func (r *Router) updateUser(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	// 3. Get update data from request body
	request := &req.UserUpdate{}
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	// 4. Call the service to perform the update
	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	err := dSvc.Update(user.ID, request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, true)
}

func (r *Router) me(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	userData, err := dSvc.Me(reqCtx.UserInfo)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, userData)
}
