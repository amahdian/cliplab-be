package router

import (
	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
)

func (r *Router) createCheckout(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	var request req.CreateCheckout
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewBillingSvc(reqCtx.Ctx)
	checkoutURL, transactionID, err := dSvc.CreateCheckout(user.ID, request.PriceID)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, resp.CheckoutResponse{
		URL:           checkoutURL,
		TransactionID: transactionID,
	})
}

func (r *Router) getCustomerPortalUrl(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	dSvc := r.svc.NewBillingSvc(reqCtx.Ctx)
	url, err := dSvc.GetCustomerPortalUrl(user.ID)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, gin.H{
		"url": url,
	})
}
