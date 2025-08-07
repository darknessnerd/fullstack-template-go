package handlers

import (
	"github.com/gin-gonic/gin"
	"webui-skeleton/internal/auth"
	"webui-skeleton/internal/config"
	"webui-skeleton/internal/database"
)

// Handlers contains all handler instances
type Handlers struct {
	Home        *HomeHandler
	Auth        *AuthHandler
	API         *APIHandler
	Health      *HealthHandler
	config      *config.Config
	db          *database.DB
	authService *auth.Service
}

// NewHandlers creates a new handlers container with all handler instances
func NewHandlers(config *config.Config, db *database.DB, authService *auth.Service) *Handlers {
	return &Handlers{
		Home:        NewHomeHandler(config, db, authService),
		Auth:        NewAuthHandler(config, db, authService),
		API:         NewAPIHandler(config, db, authService),
		Health:      NewHealthHandler(config),
		config:      config,
		db:          db,
		authService: authService,
	}
}

// SetupRoutes configures all application routes
func (h *Handlers) SetupRoutes(engine *gin.Engine) {
	// Health check
	engine.GET("/health", h.Health.HealthCheck)

	// Home page (public)
	engine.GET("/", h.Home.HomePage)

	// Authentication routes
	authGroup := engine.Group("/auth")
	{
		authGroup.GET("/login", h.Auth.LoginPage)
		authGroup.GET("/google", h.Auth.GoogleLogin)
		authGroup.GET("/google/callback", h.Auth.GoogleCallback)
		authGroup.POST("/logout", h.Auth.Logout)
		authGroup.GET("/user", h.authService.AuthMiddleware(), h.Auth.GetCurrentUser)
	}

	// Protected web routes
	webGroup := engine.Group("")
	if h.config.Auth.RequireAuth {
		webGroup.Use(h.authService.AuthMiddleware())
	} else {
		webGroup.Use(h.authService.OptionalAuthMiddleware())
	}
	{
		webGroup.GET("/dashboard", h.Home.DashboardPage)
	}

	// API routes
	apiGroup := engine.Group("/api/v1")
	{
		// Public API routes
		apiGroup.GET("/status", h.API.Status)
		apiGroup.GET("/example", h.API.GetExampleData)

		// Protected API routes (require authentication)
		protectedAPI := apiGroup.Group("")
		if h.config.Auth.RequireAuth {
			protectedAPI.Use(h.authService.AuthMiddleware())
		} else {
			protectedAPI.Use(h.authService.OptionalAuthMiddleware())
		}
		{
			protectedAPI.GET("/profile", h.Auth.GetProfile)
			protectedAPI.POST("/example", h.API.CreateExampleItem)
		}
	}

	// Static files
	engine.Static("/static", "cmd/webui-be/web/static")
}
