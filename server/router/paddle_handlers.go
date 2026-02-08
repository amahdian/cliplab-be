package router

import (
	"context"
	"io"
	"net/http"

	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
)

func (r *Router) paddleWebhook(ctx *gin.Context) {
	signature := ctx.GetHeader("Paddle-Signature")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	// Use background context for webhook processing to avoid "context canceled" errors
	// as billing and credit updates are critical operations that should finish
	// even if the client (Paddle) closes the connection prematurely.
	dSvc := r.svc.NewBillingSvc(context.Background())
	if err := dSvc.HandleWebhook(body, signature); err != nil {
		// Log the error locally since the client might be gone
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
