package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware for JWT authentication
func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		claims, err := s.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)

		c.Next()
	}
}

// OptionalAuthMiddleware creates a middleware that doesn't require authentication but sets user context if available
func (s *Service) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				token := tokenParts[1]

				// Validate token
				if claims, err := s.ValidateJWT(token); err == nil {
					// Set user information in context
					c.Set("user_id", claims.UserID)
					c.Set("user_email", claims.Email)
					c.Set("user_name", claims.Name)
				}
			}
		}

		c.Next()
	}
}

// GetUserFromContext retrieves user information from the Gin context
func GetUserFromContext(c *gin.Context) (userID int, email string, name string, exists bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, "", "", false
	}

	emailInterface, _ := c.Get("user_email")
	nameInterface, _ := c.Get("user_name")

	userID, ok := userIDInterface.(int)
	if !ok {
		return 0, "", "", false
	}

	email, _ = emailInterface.(string)
	name, _ = nameInterface.(string)

	return userID, email, name, true
}
