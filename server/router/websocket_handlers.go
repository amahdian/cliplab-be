package router

import (
	"net/http"

	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// webSocketHandler handles the WebSocket connection for real-time updates.
//
//	@Summary		Establish WebSocket connection
//	@Description	Upgrades the HTTP connection to a WebSocket connection for real-time, bidirectional communication.
//	@Description	This endpoint is used for sending server-side events, such as chat title updates.
//	@Tags			WebSocket
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		101	"Switching Protocols"
//	@Failure		401	"Unauthorized"
//	@Failure		500	"Internal Server Error"
//	@Router			/api/v1/ws [get]
func (r *Router) webSocketHandler(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Errorf("failed to upgrade websocket: %v", err)
		return
	}

	wsSvc := r.svc.NewWebSocketSvc(reqCtx.Ctx)
	// Register this user connection
	wsSvc.Register(user.ID, conn)

	defer func() {
		wsSvc.Unregister(user.ID)
		_ = conn.Close()
	}()

	// Listen for client messages if needed
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
