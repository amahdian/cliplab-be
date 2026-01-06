package router

import (
	"net/http"

	"github.com/amahdian/cliplab-be/server/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoutes() {
	r.publicGroup = r.Group("/api/v1")
	r.authGroup = r.Group(
		"/api/v1",
		middleware.VerifyAuth(r.authenticator),
	)

	r.registerPublicRoutes()
	r.registerUserRoutes()
	r.registerPostRoutes()
	r.registerWebSocketRoutes()
}

func (r *Router) registerPublicRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.publicGroup, http.MethodGet, "/health", r.healthCheck, config)
	r.registerRoute(r.publicGroup, http.MethodGet, "/server-time", r.getServerTime)
	r.registerRoute(r.publicGroup, http.MethodGet, "/swagger/*any", r.swaggerHandler, config)
}

func (r *Router) registerUserRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.publicGroup, http.MethodPost, "/users/login", r.login, config)
	r.registerRoute(r.publicGroup, http.MethodPost, "/users/register", r.register, config)
	r.registerRoute(r.publicGroup, http.MethodPost, "/users/verify", r.verify, config)
	r.registerRoute(r.authGroup, http.MethodPut, "/users/update", r.updateUser, config)
	r.registerRoute(r.authGroup, http.MethodGet, "/users/me", r.me, config)
}

func (r *Router) registerPostRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.publicGroup, http.MethodPost, "/posts/analyze", r.addPostToAnalyzeQueue, config.withMiddlewares(middleware.VerifyRecaptcha(r.configs.Recaptcha.Secret)))
	r.registerRoute(r.publicGroup, http.MethodGet, "/posts/:id", r.getPostData, config)
}

func (r *Router) registerWebSocketRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.authGroup, http.MethodGet, "/ws", r.webSocketHandler, config)
}

func (r *Router) registerRoute(routerGroup *gin.RouterGroup, method, path string, handler gin.HandlerFunc, configs ...*routeConfig) {
	config := newRouteConfig()
	if len(configs) > 0 {
		config = configs[0]
	}

	handlers := make([]gin.HandlerFunc, 0)

	//if config.RequireUserRoles {
	//	handlers = append(handlers, middleware.WithUserRoles(r.authorizer))
	//}

	if len(config.Middlewares) > 0 {
		handlers = append(handlers, config.Middlewares...)
	}

	handlers = append(handlers, handler)
	routerGroup.Handle(method, path, handlers...)
}
