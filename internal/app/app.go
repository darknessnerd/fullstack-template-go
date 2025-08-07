package app

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"

	"webui-skeleton/internal/config"
	"webui-skeleton/internal/database"
	"webui-skeleton/internal/logger"
	"webui-skeleton/internal/server"
)

// Application represents the main application
type Application struct {
	config     *config.Config
	server     *server.Server
	db         *database.DB
	templateFS embed.FS
}

// New creates a new application instance
func New(templateFS embed.FS) *Application {
	return &Application{
		templateFS: templateFS,
	}
}

// Initialize sets up all application components
func (app *Application) Initialize() error {
	// Load configuration and setup logger
	var err error
	app.config, err = config.LoadConfiguration()
	if err != nil {
		return err
	}

	// Initialize logger
	logger.Initialize(app.config.Debug, app.config.LogLevel)

	// Display banner
	displayBanner()

	// Initialize database
	app.db = database.New(&app.config.Database)
	if err := app.db.Connect(); err != nil {
		return err
	}

	// Run database migrations
	if err := app.db.Migrate(); err != nil {
		return err
	}

	// Setup server
	app.server = server.New(app.config, app.templateFS, app.db)
	app.server.SetupEngine()
	app.server.CreateHTTPServer()

	logger.Log.Info().Msg("✅ Application initialized successfully")
	return nil
}

// Run starts the application
func (app *Application) Run() error {
	// Setup graceful shutdown
	ctx, cancel := setupGracefulShutdown()
	defer cancel()

	// Start the HTTP server (this blocks until shutdown)
	if err := app.server.Start(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("❌ HTTP server failed")
		return err
	}

	logger.Log.Info().Msg("👋 Application stopped gracefully")
	return nil
}

// Cleanup performs cleanup operations
func (app *Application) Cleanup() {
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("❌ Failed to close database connection")
		}
	}
}

func displayBanner() {
	banner := `
╔═══════════════════════════════════════╗
║            WebUI Skeleton             ║
║                                       ║
║    A Go web application skeleton      ║
║    with authentication & web UI       ║
╚═══════════════════════════════════════╝
`
	logger.Log.Info().Msg(banner)
}

func setupGracefulShutdown() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	// Listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Log.Info().Msg("🛑 Received shutdown signal")
		cancel()
	}()

	return ctx, cancel
}
