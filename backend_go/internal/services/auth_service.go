package services

import (
	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
)

type AuthService struct {
	tokenSvc *TokenService
	userSvc  *UserService
	db       *db.PostgresDB
}

func NewAuthService(tokenSvc *TokenService, userSvc *UserService, db *db.PostgresDB) *AuthService {
	return &AuthService{
		tokenSvc: tokenSvc,
		userSvc:  userSvc,
		db:       db,
	}
}

func (s *AuthService) Login(c *gin.Context) {
	// Implementation stub - rate limit will be applied
	// ... login logic using s.userSvc.Login
	c.JSON(200, gin.H{"message": "login stub"})
}

func (s *AuthService) Refresh(c *gin.Context) {
	c.JSON(200, gin.H{"message": "refresh stub"})
}

func (s *AuthService) Logout(c *gin.Context) {
	c.JSON(200, gin.H{"message": "logout stub"})
}

func (s *AuthService) LogoutAll(c *gin.Context) {
	c.JSON(200, gin.H{"message": "logout all stub"})
}
