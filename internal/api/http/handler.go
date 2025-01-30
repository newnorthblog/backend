package apiHttp

import (
	"log/slog"
	"net/http"

	sloggin "github.com/samber/slog-gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/newnorthblog/backend/pkg/limiter"
	"github.com/newnorthblog/backend/pkg/validator"

	_ "github.com/newnorthblog/backend/docs"

	blogV1 "github.com/newnorthblog/backend/internal/api/http/blog/v1"
	"github.com/newnorthblog/backend/internal/pkg/tokenmanager"

	"github.com/newnorthblog/backend/internal/config"
	"github.com/newnorthblog/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	logger       *slog.Logger
	tokenManager *tokenmanager.Manager
}

func NewHandlers(
	services *service.Services,
	logger *slog.Logger,
	tokenManager *tokenmanager.Manager,
) *Handler {
	return &Handler{
		services:     services,
		logger:       logger,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	validator.RegisterGinValidator()

	router.Use(
		sloggin.NewWithConfig(h.logger, sloggin.Config{
			WithSpanID:  true,
			WithTraceID: true,
		}),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL, h.logger),
		corsMiddleware,
	)

	router.Use(gin.Recovery())

	if cfg.HTTPServer.SwaggerEnabled {
		h.logger.Info("swagger enabled")

		router.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		})

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("blogV1")))
	}

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	appHandlersV1 := blogV1.NewHandler(h.services, h.logger, h.tokenManager)
	api := router.Group("/api")
	{
		appHandlersV1.Init(api)
	}
}
