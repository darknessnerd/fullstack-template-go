package database

import (
	"database/sql"
	"fmt"

	"webui-skeleton/internal/config"
	"webui-skeleton/internal/logger"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	DB     *sql.DB
	config *config.DatabaseConfig
}

// New creates a new database instance
func New(config *config.DatabaseConfig) *DB {
	return &DB{
		config: config,
	}
}

// Connect establishes a connection to the database
func (db *DB) Connect() error {
	var err error

	logger.Log.Info().
		Str("type", string(db.config.Type)).
		Str("database", db.config.Database).
		Msg("Connecting to database")

	db.DB, err = sql.Open(db.config.GetDriverName(), db.config.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.DB.SetMaxOpenConns(db.config.MaxOpenConns)
	db.DB.SetMaxIdleConns(db.config.MaxIdleConns)
	db.DB.SetConnMaxLifetime(db.config.ConnMaxLifetime)

	// Test the connection
	if err := db.DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log.Info().Msg("✅ Database connection established")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		logger.Log.Info().Msg("Closing database connection")
		return db.DB.Close()
	}
	return nil
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	logger.Log.Info().Msg("Running database migrations")

	// Create users table
	usersSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			google_id VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			picture VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`

	if db.config.Type == config.PostgreSQL {
		usersSQL = `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				google_id VARCHAR(255) UNIQUE NOT NULL,
				email VARCHAR(255) UNIQUE NOT NULL,
				name VARCHAR(255) NOT NULL,
				picture VARCHAR(500),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	}

	if _, err := db.DB.Exec(usersSQL); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create sessions table for session management
	sessionsSQL := `
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token VARCHAR(500) NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`

	if db.config.Type == config.PostgreSQL {
		sessionsSQL = `
			CREATE TABLE IF NOT EXISTS sessions (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL,
				token VARCHAR(500) NOT NULL,
				expires_at TIMESTAMP NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			)`
	}

	if _, err := db.DB.Exec(sessionsSQL); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	logger.Log.Info().Msg("✅ Database migrations completed")
	return nil
}
