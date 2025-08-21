# Authentication

WebUI Skeleton provides built-in authentication using Google OAuth2 and JWT tokens.

## Google OAuth2
- Configure your Google credentials in the `.env` file:
  - `GOOGLE_CLIENT_ID`
  - `GOOGLE_CLIENT_SECRET`
  - `GOOGLE_REDIRECT_URL`
- Users can log in via Google, and their profile is stored in the database.

## JWT Token Management
- After login, a JWT token is issued and stored in a secure cookie.
- The token is used to authenticate API requests and access protected routes.
- Configure the secret and expiration in `.env`:
  - `JWT_SECRET`
  - `JWT_EXPIRES_IN`

## Middleware
- `AuthMiddleware()`: Requires a valid JWT token for access.
- `OptionalAuthMiddleware()`: Allows access but sets user context if token is present.

## Example Usage
- Protect routes by adding the authentication middleware in `internal/server/routes.go`.
- Retrieve user info in handlers using `auth.GetUserFromContext(c)`.

See the source code in `internal/auth/` for implementation details.
