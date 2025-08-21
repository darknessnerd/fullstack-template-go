# Configuration

WebUI Skeleton uses environment variables for configuration. You can set these in a `.env` file or via your shell environment.

## Main Variables

### Server
- `SERVER_HOST`: Bind address (default: 0.0.0.0)
- `SERVER_PORT`: Port (default: 8080)

### Database
- `DB_TYPE`: `sqlite` or `postgresql` (default: sqlite)
- `DB_DATABASE`: Database name/path (default: app.db)
- `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`: PostgreSQL settings
- `DB_SSL_MODE`: PostgreSQL SSL mode (default: disable)

### Authentication
- `JWT_SECRET`: Secret for JWT tokens
- `JWT_EXPIRES_IN`: Token expiration (default: 24h)
- `JWT_ISSUER`: Token issuer (default: webui-skeleton)
- `REQUIRE_AUTH`: Require authentication for all routes (default: false)
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`: Google OAuth credentials
- `SESSION_SECRET`: Session secret key

### Logging
- `DEBUG`: Enable debug mode (default: false)
- `LOG_LEVEL`: Log level (trace/debug/info/warn/error/fatal/panic, default: info)

## Example `.env` File
```
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_TYPE=sqlite
DB_DATABASE=app.db
JWT_SECRET=your-secret
JWT_EXPIRES_IN=24h
REQUIRE_AUTH=false
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
SESSION_SECRET=your-session-secret
DEBUG=true
LOG_LEVEL=info
```

See `internal/config/config.go` for implementation details.
