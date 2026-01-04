package router

import (
	"fmt"
	"net/url"
	"time"

	"github.com/amahdian/cliplab-be/docs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/amahdian/cliplab-be/svc"
	"github.com/amahdian/cliplab-be/svc/auth"

	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/server/binding"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	validator "github.com/gin-gonic/gin/binding"

	"github.com/amahdian/cliplab-be/server/middleware"
)

type Router struct {
	*gin.Engine

	storage   storage.PgStorage
	svc       svc.Svc
	validator validator.StructValidator

	configs *env.Envs

	authenticator auth.Authenticator

	publicGroup *gin.RouterGroup
	authGroup   *gin.RouterGroup
}

func NewRouter(
	svc svc.Svc,
	configs *env.Envs,
	authenticator auth.Authenticator) *Router {
	gin.SetMode(configs.Server.GinMode)
	router := &Router{
		Engine:        gin.New(),
		svc:           svc,
		validator:     validator.Validator,
		configs:       configs,
		authenticator: authenticator,
	}
	router.Use(
		middleware.WithLogger(),
		middleware.WithRecovery(),
	)
	pprof.Register(router.Engine)
	router.setupBindings()
	router.setupCors()
	//router.setupSwagger()
	router.setupRoutes()

	return router
}

func (r *Router) setupBindings() {
	binding.Init()
}

func (r *Router) setupCors() {
	origins := []string{"https://cliplab.dev"}
	if r.configs.Server.LocalCors {
		origins = append(origins, "http://localhost:3000")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
}

func (r *Router) setupSwagger() {
	scheme := "http"
	host := fmt.Sprintf("localhost:%s", r.configs.Server.HttpPort)
	basePath := "/"

	if r.configs.Server.SwaggerHostAddr != "" {
		uri, err := url.ParseRequestURI(r.configs.Server.SwaggerHostAddr)
		if err != nil {
			logger.Errorf("Failed to parse swagger host address: %v", err)
		} else {
			scheme = uri.Scheme
			host = uri.Host
			basePath = uri.Path
		}
	}

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = []string{scheme}
}
