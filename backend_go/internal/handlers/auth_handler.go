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

func (h *AuthHandler) Refresh(c *gin.Context) {
	h.authSvc.Refresh(c)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.authSvc.Logout(c)
}

func (h *AuthHandler) Register(c *gin.Context) {
	h.authSvc.Register(c)
}
