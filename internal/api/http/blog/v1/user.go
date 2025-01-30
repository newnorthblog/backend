package v1

import (
	"errors"
	"net/http"

	"github.com/newnorthblog/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	users.POST("/register", h.userRegister)
	users.POST("/login", h.userAuth)
	users.POST("ping", h.userIdentityMiddleware, h.ping)
}

type userRegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

// @Summary Регистрация
// @Tags Client
// @Description Регистрация
// @ModuleID Client
// @Accept  json
// @Produce  json
// @Param input body userRegisterRequest true "Регистрация"
// @Success 201
// @Failure 400 {object} ErrorStruct
// @Router /users/register [post]
func (h *Handler) userRegister(c *gin.Context) {
	var req userRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	if err := h.services.Users.Register(c.Request.Context(), &service.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			errorResponse(c, UserAlreadyExistsCode)
			return
		}
		h.logger.Error("failed to register user",
			"error", err,
		)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusCreated)
}

type userLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userLoginResponse struct {
	AccessToken string `json:"access_token"`
}

// @Summary Авторизация
// @Tags Client
// @Description Авторизация
// @ModuleID Client
// @Accept  json
// @Produce  json
// @Param input body userLoginRequest true "Авторизация"
// @Success 200 {object} userLoginResponse
// @Failure 400 {object} ErrorStruct
// @Router /users/login [post]
func (h *Handler) userAuth(c *gin.Context) {
	var req userLoginRequest
	if err := c.BindJSON(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	token, err := h.services.Users.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrUserInvalidCredentials) {
			errorResponse(c, UserNotFoundCode)
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, UserNotFoundCode)
			return
		}

		h.logger.Error("failed to login client",
			"error", err,
		)
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, userLoginResponse{AccessToken: token.AccessToken})
}

// @Summary Ping
// @Tags Client
// @Description Проверка доступности сервера
// @ModuleID Client
// @Accept  json
// @Produce  json
// @Success 200 {string} string "pong"
// @Failure 400 {object} ErrorStruct
// @Router /users/ping [post]
// @Security Bearer
func (h *Handler) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
