package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/services"
)

type AuthHandler struct {
	authSvc *services.AuthService
}

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.authSvc.Login(c)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token from cookie
// @Tags Auth
// @Security CookieAuth
// @Success 200 {object} TokenResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	h.authSvc.Refresh(c)
}

// Logout godoc
// @Summary Logout user
// @Description Revoke tokens and clear cookies
// @Tags Auth
// @Security CookieAuth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	h.authSvc.Logout(c)
}

func (h *AuthHandler) Register(c *gin.Context) {
	h.authSvc.Register(c)
}
