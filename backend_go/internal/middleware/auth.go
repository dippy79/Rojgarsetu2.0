package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rojgarsetu/backend/config"
)

type Claims struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	MFAVerified bool   `json:"mfa_verified"`
	jwt.RegisteredClaims
}

// var JWT_SECRET []byte

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Authorization header required",
				},
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Bearer token required",
				},
			})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Invalid token claims",
				},
			})
			c.Abort()
			return
		}

		// Validate issuer and audience
		if claims.Issuer != cfg.JWT.Issuer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Invalid token issuer",
				},
			})
			c.Abort()
			return
		}

		if len(claims.Audience) > 0 && !contains(claims.Audience, cfg.JWT.Audience) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Invalid token audience",
				},
			})
			c.Abort()
			return
		}

		// Validate exp claim explicitly - redundant since token.Valid already checks
		// claims.Valid() handled by jwt library in token.Valid

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RoleMiddleware restricts access based on user roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "User role not found in token",
				},
			})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    401,
					"message": "Invalid user role",
				},
			})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    403,
				"message": "Access denied. Insufficient permissions.",
			},
		})
		c.Abort()
	}
}

// CandidateMiddleware restricts access to candidates only
func CandidateMiddleware() gin.HandlerFunc {
	return RoleMiddleware("candidate")
}

// CompanyMiddleware restricts access to companies only
func CompanyMiddleware() gin.HandlerFunc {
	return RoleMiddleware("company")
}

// AdminMiddleware restricts access to admin and superadmin
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin", "superadmin")
}

// AdminMFAMiddleware requires a valid admin JWT with mfa_verified=true.
func AdminMFAMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			if tokenString == "" {
				if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
					tokenString = cookie
				}
			}
			if tokenString == "" {
				c.JSON(http.StatusForbidden, gin.H{"error": "MFA required"})
				c.Abort()
				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWT.Secret), nil
			})
			if err != nil || !token.Valid {
				c.JSON(http.StatusForbidden, gin.H{"error": "MFA required"})
				c.Abort()
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || !claims.MFAVerified {
				c.JSON(http.StatusForbidden, gin.H{"error": "MFA required"})
				c.Abort()
				return
			}

			c.Set("role", claims.Role)
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("mfa_verified", claims.MFAVerified)
			userRole = claims.Role
		}

		roleStr, ok := userRole.(string)
		if !ok || (roleStr != "admin" && roleStr != "superadmin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "MFA required"})
			c.Abort()
			return
		}

		if mfaVerified, ok := c.Get("mfa_verified"); !ok || mfaVerified != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "MFA required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SuperAdminMiddleware restricts access to superadmin only
func SuperAdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("superadmin")
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}
