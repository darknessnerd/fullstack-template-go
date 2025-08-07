package server

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"webui-skeleton/internal/auth"
	"webui-skeleton/internal/config"
	"webui-skeleton/internal/database"
	"webui-skeleton/internal/handlers"
	"webui-skeleton/internal/logger"
)

type Server struct {
	config      *config.Config
	templateFS  embed.FS
	db          *database.DB
	authService *auth.Service
	handlers    *handlers.Handlers
	engine      *gin.Engine
	httpServer  *http.Server
}

// New creates a new server instance
func New(config *config.Config, templateFS embed.FS, db *database.DB) *Server {
	return &Server{
		config:     config,
		templateFS: templateFS,
		db:         db,
	}
}

// SetupEngine configures the Gin engine with routes and middleware
func (s *Server) SetupEngine() {
	// Set Gin mode based on debug flag
	if !s.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	s.engine = gin.New()

	// Add basic middleware
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())

	// Load HTML templates
	s.engine.LoadHTMLGlob("cmd/webui-be/web/templates/*")

	// Setup authentication service
	s.authService = auth.NewService(
		s.db.DB,
		s.config.Auth.JWTSecret,
		s.config.Auth.JWTExpiresIn,
		s.config.Auth.JWTIssuer,
		s.config.Auth.GoogleClientID,
		s.config.Auth.GoogleClientSecret,
		s.config.Auth.GoogleRedirectURL,
	)

	// Initialize handlers
	s.handlers = handlers.NewHandlers(s.config, s.db, s.authService)

	// Setup routes
	s.setupRoutes()

	logger.Log.Info().Msg("✅ Server engine configured")
}

// CreateHTTPServer creates the HTTP server instance
func (s *Server) CreateHTTPServer() {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Log.Info().
		Str("address", addr).
		Msg("HTTP server configured")
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	// Start server in a goroutine
	go func() {
		logger.Log.Info().
			Str("address", s.httpServer.Addr).
			Msg("🚀 Starting HTTP server")

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("❌ HTTP server failed to start")
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	logger.Log.Info().Msg("🛑 Shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("❌ HTTP server forced to shutdown")
		return err
	}

	logger.Log.Info().Msg("✅ HTTP server stopped gracefully")
	return nil
}
