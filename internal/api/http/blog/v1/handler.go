package v1

import (
	"log/slog"

	"github.com/newnorthblog/backend/internal/pkg/tokenmanager"
	"github.com/newnorthblog/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// @title New-North Backend API
// @version 1.0
// @description Backend API for New-North Blog
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

type Handler struct {
	services     *service.Services
	logger       *slog.Logger
	tokenManager *tokenmanager.Manager
}

func NewHandler(
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("v1")
	h.initUserRoutes(v1)
}
