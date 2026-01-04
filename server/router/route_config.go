package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type routeConfig struct {
	RequireUserRoles bool
	Middlewares      []gin.HandlerFunc
}

func newRouteConfig() *routeConfig {
	return &routeConfig{
		RequireUserRoles: false,
		Middlewares:      []gin.HandlerFunc{},
	}
}

func (rc *routeConfig) withUserRoles(flag bool) *routeConfig {
	clone := rc.clone()
	clone.RequireUserRoles = flag
	return clone
}

func (rc *routeConfig) withMiddlewares(middlewares ...gin.HandlerFunc) *routeConfig {
	clone := rc.clone()
	clone.Middlewares = append(rc.Middlewares, middlewares...)
	return clone
}

func (rc *routeConfig) withCompression() *routeConfig {
	clone := rc.clone()
	clone.Middlewares = append(rc.Middlewares, gzip.Gzip(gzip.DefaultCompression))
	return clone
}

func (rc *routeConfig) clone() *routeConfig {
	middlewares := make([]gin.HandlerFunc, len(rc.Middlewares))
	copy(middlewares, rc.Middlewares)

	return &routeConfig{
		RequireUserRoles: rc.RequireUserRoles,
		Middlewares:      middlewares,
	}
}
