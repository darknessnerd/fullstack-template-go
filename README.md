# WebUI Skeleton

A Go web application skeleton with built-in authentication and web UI components, extracted from the card-swap project.

## Features

- ✅ Google OAuth2 Authentication
- ✅ JWT Token Management
- ✅ Database Support (SQLite & PostgreSQL)
- ✅ RESTful API Endpoints
- ✅ Web UI Templates
- ✅ Middleware Support
- ✅ Configuration Management
- ✅ Structured Logging with Zerolog
- ✅ Graceful Shutdown
- ✅ Health Check Endpoints
- ✅ Live Reloading with Air

## Starting a New Project from the Skeleton

**Method 1: Manual Copy**
- Copy all files from webui-skeleton to your new project directory.
- Rename the directory and update project names as needed.
- Install dependencies:
  ```bash
  go mod tidy
  ```
- Copy .env.example to .env and configure your environment variables.
- Start development as described in the Quick Start section.

**Method 2: Using Git**
- Clone the skeleton repository:
  ```bash
  git clone git@github.com:darknessnerd/fullstack-template-go.git my-new-project
  cd my-new-project
  ```
- Remove the existing git history:
  ```bash
  rm -rf .git
  ```
- Initialize a new git repository:
  ```bash
  git init
  git add .
  git commit -m "Initial commit from skeleton"
  ```
- Continue with dependency installation and configuration as above.

Both methods give you a clean starting point for your own project.

## Quick Start

1. **Clone and setup**:
   ```bash
   cd webui-skeleton
   cp .env.example .env
   ```

2. **Install Air for live reloading** (recommended for development):
   ```bash
   go install github.com/air-verse/air@latest
   ```

3. **Configure Google OAuth** (optional, but recommended):
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing
   - Enable Google+ API
   - Create OAuth 2.0 credentials
   - Update `.env` with your credentials

4. **Install dependencies**:
   ```bash
   go mod tidy
   ```

5. **Run the application**:
   
   **With live reloading (recommended for development):**
   ```bash
   air
   ```
   
   **Or run normally:**
   ```bash
   go run cmd/webui-be/main.go
   ```

6. **Access the application**:
   - Open http://localhost:8080
   - Health check: http://localhost:8080/health
   - API status: http://localhost:8080/api/v1/status

## Development with Air

Air provides live reloading during development. When you run `air`, it will:

- 🔄 **Auto-reload** when you modify Go files, templates, or configuration
- 🚀 **Fast rebuilds** with intelligent file watching
- 🎨 **Colored output** for easy debugging
- 📁 **Smart ignoring** of vendor, tmp, and test files

The Air configuration is in `.air.toml` and watches:
- `.go` files (source code)
- `.html` files (templates)
- `.tpl`, `.tmpl` files (additional templates)

Build artifacts are stored in `./tmp/` directory.

## Configuration

The application supports configuration via environment variables or `.env` file:

### Server
- `SERVER_HOST`: Server bind address (default: 0.0.0.0)
- `SERVER_PORT`: Server port (default: 8080)

### Database
- `DB_TYPE`: Database type (sqlite/postgresql, default: sqlite)
- `DB_DATABASE`: Database name/path (default: app.db)
- `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`: PostgreSQL settings
- `DB_SSL_MODE`: PostgreSQL SSL mode (default: disable)

### Authentication
- `JWT_SECRET`: Secret key for JWT tokens (required for production)
- `JWT_EXPIRES_IN`: Token expiration time (default: 24h)
- `JWT_ISSUER`: Token issuer (default: webui-skeleton)
- `REQUIRE_AUTH`: Require authentication for all routes (default: false)
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`: Google OAuth credentials
- `GOOGLE_REDIRECT_URL`: OAuth redirect URL
- `SESSION_SECRET`: Session secret key

### Logging
- `DEBUG`: Enable debug mode (default: false)
- `LOG_LEVEL`: Log level (trace/debug/info/warn/error/fatal/panic, default: info)

## API Endpoints

### Public Endpoints
- `GET /health` - Health check
- `GET /api/v1/status` - API status
- `GET /` - Home page
- `GET /auth/login` - Login page
- `GET /auth/google` - Google OAuth login
- `GET /auth/google/callback` - OAuth callback
- `POST /auth/logout` - Logout

### Protected Endpoints (require authentication when REQUIRE_AUTH=true)
- `GET /dashboard` - User dashboard
- `GET /auth/user` - Current user info
- `GET /api/v1/profile` - User profile

## Project Structure

```
webui-skeleton/
├── cmd/webui-be/           # Main application
│   ├── main.go
│   └── web/templates/      # HTML templates
├── internal/               # Internal packages
│   ├── app/               # Application setup
│   ├── auth/              # Authentication service
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   ├── logger/            # Logging setup
│   └── server/            # HTTP server and routes
├── .env.example           # Environment configuration example
├── go.mod                 # Go module file
└── README.md             # This file
```

## Development

### Adding New Routes

1. Add route handlers in `internal/server/handlers.go`
2. Register routes in `internal/server/server.go` (`setupRoutes` method)
3. Add templates in `cmd/webui-be/web/templates/` if needed

### Adding New Database Tables

1. Add migration SQL in `internal/database/database.go` (`Migrate` method)
2. Create repository interfaces and implementations
3. Update the application setup in `internal/app/app.go`

### Authentication

The skeleton includes two middleware options:
- `AuthMiddleware()`: Requires valid JWT token
- `OptionalAuthMiddleware()`: Sets user context if token is present, but doesn't require it

Use `auth.GetUserFromContext(c)` to retrieve user information in handlers.

## Building for Production

1. **Set production environment variables**:
   ```bash
   export DEBUG=false
   export JWT_SECRET=your-production-secret
   export GOOGLE_CLIENT_ID=your-production-client-id
   # ... other production settings
   ```

2. **Build the application**:
   ```bash
   go build -o webui-skeleton cmd/webui-be/main.go
   ```

3. **Run**:
   ```bash
   ./webui-skeleton
   ```

## Docker Support

Create a `Dockerfile`:
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o webui-skeleton cmd/webui-be/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/webui-skeleton .
CMD ["./webui-skeleton"]
```

## License

This skeleton is provided as-is for development purposes. Modify as needed for your projects.
