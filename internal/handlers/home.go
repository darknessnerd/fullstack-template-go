package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"webui-skeleton/internal/auth"
	"webui-skeleton/internal/config"
	"webui-skeleton/internal/database"
)

// HomeHandler handles home and dashboard pages
type HomeHandler struct {
	config  *config.Config
	db      *database.DB
	authSvc *auth.Service
}

// NewHomeHandler creates a new home handler
func NewHomeHandler(config *config.Config, db *database.DB, authSvc *auth.Service) *HomeHandler {
	return &HomeHandler{
		config:  config,
		db:      db,
		authSvc: authSvc,
	}
}

// HomePage handles the main home page
func (h *HomeHandler) HomePage(c *gin.Context) {
	// Get user info if authenticated
	userID, email, name, isAuthenticated := auth.GetUserFromContext(c)

	// Example data to display
	exampleData := gin.H{
		"items": []gin.H{
			{"id": 1, "name": "Example Item 1", "description": "This is a sample item to demonstrate data binding"},
			{"id": 2, "name": "Example Item 2", "description": "Another sample item showing template rendering"},
			{"id": 3, "name": "Example Item 3", "description": "A third item to showcase the grid layout"},
		},
		"stats": gin.H{
			"total_items":  3,
			"last_updated": "2025-08-07",
		},
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":           "WebUI Skeleton - Home",
		"isAuthenticated": isAuthenticated,
		"userID":          userID,
		"email":           email,
		"name":            name,
		"data":            exampleData,
	})
}

// DashboardPage handles the dashboard page (protected route)
func (h *HomeHandler) DashboardPage(c *gin.Context) {
	userID, email, name, _ := auth.GetUserFromContext(c)

	// Example dashboard data
	dashboardData := gin.H{
		"user_stats": gin.H{
			"login_count":   42,
			"last_activity": "2025-08-07 14:30:00",
			"member_since":  "2025-01-15",
		},
		"recent_activities": []gin.H{
			{"action": "Logged in", "timestamp": "2025-08-07 14:30:00"},
			{"action": "Viewed profile", "timestamp": "2025-08-07 12:15:00"},
			{"action": "API call made", "timestamp": "2025-08-07 10:45:00"},
		},
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":  "Dashboard",
		"userID": userID,
		"email":  email,
		"name":   name,
		"data":   dashboardData,
	})
}
