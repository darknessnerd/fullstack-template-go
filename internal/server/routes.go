package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Health check endpoints
	s.engine.GET("/health", s.handleHealth)
	s.engine.GET("/hh", s.handleHealth)

	// Authentication routes (unprotected)
	s.setupAuthRoutes()

	// Web routes (optionally protected)
	s.setupWebRoutes()

	// API routes (protected)
	s.setupAPIRoutes()

	// Admin routes (protected)
	s.setupAdminRoutes()
}

// setupAuthRoutes configures authentication routes
func (s *Server) setupAuthRoutes() {
	authGroup := s.engine.Group("/auth")
	{
		authGroup.GET("/google/login", s.handlers.Auth.GoogleLogin)
		authGroup.GET("/google/callback", s.handlers.Auth.GoogleCallback)
		authGroup.GET("/logout", s.handlers.Auth.Logout)  // Changed from POST to GET
		authGroup.POST("/logout", s.handlers.Auth.Logout) // Keep POST for API compatibility

		// Protected auth routes
		protected := authGroup.Group("")
		protected.Use(s.authService.Middleware())
		{
			protected.GET("/profile", s.handlers.Auth.GetProfile)
		}
	}
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// setupWebRoutes configures web-related routes
func (s *Server) setupWebRoutes() {
	// Login page (no authentication required)
	s.engine.GET("/login", s.handleLoginPage)

	// Main application routes (require authentication when REQUIRE_AUTH is true)
	webGroup := s.engine.Group("/")
	if s.config.Auth.RequireAuth {
		// If authentication is required, use the mandatory middleware
		webGroup.Use(s.authService.Middleware())
	} else {
		// If authentication is optional, use the optional middleware
		webGroup.Use(s.authService.OptionalMiddleware())
	}
	{
		webGroup.GET("/", s.handleHomePage)
		webGroup.GET("/dashboard", s.handlers.Home.DashboardPage)
		// Add other web routes here
	}
}

// setupAPIRoutes configures API routes with authentication
func (s *Server) setupAPIRoutes() {
	apiGroup := s.engine.Group("/api/v1")

	// Public API routes
	apiGroup.GET("/status", s.handlers.API.Status)
	apiGroup.GET("/example", s.handlers.API.GetExampleData)

	// All other API routes require authentication
	protected := apiGroup.Group("")
	protected.Use(s.authService.Middleware())
	{
		// Example API endpoints (protected)
		protected.POST("/example", s.handlers.API.CreateExampleItem)
		protected.GET("/items", s.handlers.API.GetItems)

		// User profile endpoints
		protected.GET("/profile", s.handlers.Auth.GetProfile)
	}
}

// setupAdminRoutes configures admin routes (all protected)
func (s *Server) setupAdminRoutes() {
	adminGroup := s.engine.Group("/admin")
	adminGroup.Use(s.authService.Middleware())
	{
		// Admin dashboard
		adminGroup.GET("/", s.handleAdminDashboard)

		// User management
		adminGroup.GET("/users", s.handleAdminUsers)

		// System info
		adminGroup.GET("/system", s.handleAdminSystem)
	}
}

// handleLoginPage handles the login page
func (s *Server) handleLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - WebUI Skeleton",
	})
}

// handleHomePage handles the home page
func (s *Server) handleHomePage(c *gin.Context) {
	// Delegate to the Home handler
	s.handlers.Home.HomePage(c)
}

// handleAdminDashboard handles the admin dashboard
func (s *Server) handleAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title": "Admin Dashboard",
	})
}

// handleAdminUsers handles the admin users page
func (s *Server) handleAdminUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin users endpoint",
		"users":   []string{"admin", "user1", "user2"},
	})
}

// handleAdminSystem handles the admin system info page
func (s *Server) handleAdminSystem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "System info endpoint",
		"system": gin.H{
			"version": "1.0.0",
			"uptime":  "running",
		},
	})
}
