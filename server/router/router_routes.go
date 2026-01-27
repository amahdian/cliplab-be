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
	r.registerAnalyzeRoutes()
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
	recaptchaMiddleware := middleware.VerifyRecaptcha(r.configs.Recaptcha.Secret)

	r.registerRoute(r.publicGroup, http.MethodPost, "/users/login", r.login, config.withMiddlewares(recaptchaMiddleware))
	r.registerRoute(r.publicGroup, http.MethodPost, "/users/login/oauth", r.loginOauth, config.withMiddlewares(recaptchaMiddleware))
	r.registerRoute(r.publicGroup, http.MethodPost, "/users/register", r.register, config.withMiddlewares(recaptchaMiddleware))
	r.registerRoute(r.publicGroup, http.MethodPost, "/users/verify", r.verify, config)
	r.registerRoute(r.authGroup, http.MethodPut, "/users/update", r.updateUser, config)
	r.registerRoute(r.authGroup, http.MethodGet, "/users/me", r.me, config)
}

func (r *Router) registerAnalyzeRoutes() {
	config := newRouteConfig()
	//r.registerRoute(r.publicGroup, http.MethodPost, "/analyze", r.addRequestToAnalyzeQueue, config)
	r.registerRoute(r.publicGroup, http.MethodPost, "/analyze", r.addRequestToAnalyzeQueue, config.withMiddlewares(middleware.VerifyRecaptcha(r.configs.Recaptcha.Secret)))
	r.registerRoute(r.publicGroup, http.MethodGet, "/analyze/:id", r.getAnalyzeResult, config)
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
