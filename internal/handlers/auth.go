package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"webui-skeleton/internal/auth"
	"webui-skeleton/internal/config"
	"webui-skeleton/internal/database"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	config  *config.Config
	db      *database.DB
	authSvc *auth.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(config *config.Config, db *database.DB, authSvc *auth.Service) *AuthHandler {
	return &AuthHandler{
		config:  config,
		db:      db,
		authSvc: authSvc,
	}
}

// LoginPage displays the login page
func (h *AuthHandler) LoginPage(c *gin.Context) {
	// Check if user is already authenticated
	if userID, _, _, exists := auth.GetUserFromContext(c); exists && userID > 0 {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - WebUI Skeleton",
	})
}

// GoogleLogin initiates Google OAuth login
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authSvc.GetGoogleAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Exchange code for user info
	userInfo, err := h.authSvc.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Create or update user
	user, err := h.authSvc.CreateOrUpdateUser(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user"})
		return
	}

	// Generate JWT token
	token, err := h.authSvc.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set auth cookie for web UI
	c.SetCookie("auth_token", token, int(h.config.Auth.JWTExpiresIn.Seconds()), "/", "", false, true)

	// Redirect to dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the auth cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

// GetCurrentUser returns the current authenticated user info
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, email, name, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    userID,
		"email": email,
		"name":  name,
	})
}

// GetProfile returns the full user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _, _, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.authSvc.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
	})
}
