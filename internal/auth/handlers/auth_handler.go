package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	model "github.com/thangchung/go-coffeeshop/internal/auth/domain"
	"github.com/thangchung/go-coffeeshop/internal/auth/service"
)

// AuthHandler wires HTTP routes to the AuthService.
type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// RegisterRoutes attaches all auth routes to the given router group.
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
	rg.GET("/health", h.Health)
}

// Register godoc
// POST /auth/register
// Body: { "username": "...", "email": "...", "password": "..." }
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Register(&req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		case errors.Is(err, service.ErrUsernameTaken):
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		default:
			logrus.WithError(err).Error("register failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// POST /auth/login
// Body: { "email": "...", "password": "..." }
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCreds):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		default:
			logrus.WithError(err).Error("login failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Health is a simple liveness probe endpoint.
// GET /auth/health
func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "auth-service"})
}
