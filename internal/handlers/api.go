package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"webui-skeleton/internal/auth"
	"webui-skeleton/internal/config"
	"webui-skeleton/internal/database"
)

// APIHandler handles API endpoints
type APIHandler struct {
	config  *config.Config
	db      *database.DB
	authSvc *auth.Service
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(config *config.Config, db *database.DB, authSvc *auth.Service) *APIHandler {
	return &APIHandler{
		config:  config,
		db:      db,
		authSvc: authSvc,
	}
}

// Status returns API status information
func (h *APIHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"api":     "v1",
		"app":     "webui-skeleton",
	})
}

// GetExampleData returns example data for demonstration
func (h *APIHandler) GetExampleData(c *gin.Context) {
	// Simulate some business logic
	data := gin.H{
		"message": "Hello from the webui-skeleton API!",
		"data": []gin.H{
			{"id": 1, "value": "Sample API data 1", "active": true, "category": "demo"},
			{"id": 2, "value": "Sample API data 2", "active": false, "category": "test"},
			{"id": 3, "value": "Sample API data 3", "active": true, "category": "example"},
		},
		"meta": gin.H{
			"total":     3,
			"page":      1,
			"per_page":  10,
			"timestamp": "2025-08-07T12:00:00Z",
		},
	}

	c.JSON(http.StatusOK, data)
}

// CreateExampleItem creates a new example item
func (h *APIHandler) CreateExampleItem(c *gin.Context) {
	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Active      bool   `json:"active"`
		Category    string `json:"category"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get user info for auditing
	userID, _, _, _ := auth.GetUserFromContext(c)

	// Simulate creating an item (in a real app, this would save to database)
	response := gin.H{
		"id":          42, // fake generated ID
		"name":        request.Name,
		"description": request.Description,
		"active":      request.Active,
		"category":    request.Category,
		"created_by":  userID,
		"created_at":  "2025-08-07T12:00:00Z",
		"updated_at":  "2025-08-07T12:00:00Z",
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Example item created successfully",
		"item":    response,
	})
}

// GetItems returns a paginated list of items (example endpoint)
func (h *APIHandler) GetItems(c *gin.Context) {
	// Get query parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	category := c.Query("category")

	// Simulate pagination and filtering
	items := []gin.H{
		{"id": 1, "name": "Item 1", "category": "demo", "active": true},
		{"id": 2, "name": "Item 2", "category": "test", "active": false},
		{"id": 3, "name": "Item 3", "category": "example", "active": true},
	}

	// Filter by category if provided
	if category != "" {
		filtered := []gin.H{}
		for _, item := range items {
			if item["category"] == category {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"meta": gin.H{
			"page":     page,
			"limit":    limit,
			"total":    len(items),
			"category": category,
		},
	})
}
