package router

import (
	"net/http"

	"github.com/amahdian/cliplab-be/server/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoutes() {
	r.publicGroup = r.Group("/v1")
	r.authGroup = r.Group(
		"/v1",
		middleware.VerifyAuth(r.authenticator),
	)
	r.optionalAuthGroup = r.Group(
		"/v1",
		middleware.WithAuth(r.authenticator),
	)

	r.registerPublicRoutes()
	r.registerUserRoutes()
	r.registerAnalyzeRoutes()
	r.registerChannelRoutes()
	r.registerWebSocketRoutes()
	r.registerWebhookRoutes()
	r.registerBillingRoutes()
}

func (r *Router) registerBillingRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.authGroup, http.MethodPost, "/billing/checkout", r.createCheckout, config)
	r.registerRoute(r.authGroup, http.MethodGet, "/billing/portal", r.getCustomerPortalUrl, config)
}

func (r *Router) registerWebhookRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.publicGroup, http.MethodPost, "/webhook/paddle", r.paddleWebhook, config)
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
	r.registerRoute(r.authGroup, http.MethodGet, "/users/history", r.getUserHistory, config)
}

func (r *Router) registerAnalyzeRoutes() {
	config := newRouteConfig()
	//r.registerRoute(r.publicGroup, http.MethodPost, "/analyze", r.addRequestToAnalyzeQueue, config)
	r.registerRoute(r.optionalAuthGroup, http.MethodPost, "/analyze", r.addRequestToAnalyzeQueue, config.withMiddlewares(middleware.VerifyRecaptcha(r.configs.Recaptcha.Secret)))
	r.registerRoute(r.publicGroup, http.MethodGet, "/analyze/:id", r.getAnalyzeResult, config)
	r.registerRoute(r.publicGroup, http.MethodGet, "/instagram/image", r.getProxyImage, config)
}

func (r *Router) registerChannelRoutes() {
	config := newRouteConfig()
	recaptchaMiddleware := middleware.VerifyRecaptcha(r.configs.Recaptcha.Secret)

	r.registerRoute(r.optionalAuthGroup, http.MethodGet, "/channels/engagement-rate", r.getChannelEngagementRate, config.withMiddlewares(recaptchaMiddleware))
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
