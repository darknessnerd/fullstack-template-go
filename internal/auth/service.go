package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	db           *sql.DB
	jwtSecret    []byte
	jwtExpiresIn time.Duration
	jwtIssuer    string
	googleConfig oauth2.Config
}

// NewService creates a new authentication service
func NewService(db *sql.DB, jwtSecret string, jwtExpiresIn time.Duration, jwtIssuer string,
	googleClientID, googleClientSecret, googleRedirectURL string) *Service {

	googleConfig := oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Service{
		db:           db,
		jwtSecret:    []byte(jwtSecret),
		jwtExpiresIn: jwtExpiresIn,
		jwtIssuer:    jwtIssuer,
		googleConfig: googleConfig,
	}
}

// GenerateJWT generates a JWT token for a user
func (s *Service) GenerateJWT(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(s.jwtExpiresIn).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.jwtIssuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT validates a JWT token and returns the claims
func (s *Service) ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid user_id in token")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid email in token")
		}

		name, ok := claims["name"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid name in token")
		}

		return &JWTClaims{
			UserID: int(userID),
			Email:  email,
			Name:   name,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (s *Service) GetGoogleAuthURL() string {
	state := s.generateRandomState()
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken exchanges an authorization code for user info
func (s *Service) ExchangeCodeForToken(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	client := s.googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// CreateOrUpdateUser creates or updates a user from Google OAuth info
func (s *Service) CreateOrUpdateUser(userInfo *GoogleUserInfo) (*User, error) {
	// Check if user exists
	var user User
	err := s.db.QueryRow(`
		SELECT id, google_id, email, name, picture, created_at, updated_at 
		FROM users WHERE google_id = ? OR email = ?`,
		userInfo.ID, userInfo.Email).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		// Create new user
		result, err := s.db.Exec(`
			INSERT INTO users (google_id, email, name, picture) 
			VALUES (?, ?, ?, ?)`,
			userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		userID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get user ID: %w", err)
		}

		user.ID = int(userID)
		user.GoogleID = userInfo.ID
		user.Email = userInfo.Email
		user.Name = userInfo.Name
		user.Picture = userInfo.Picture
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	} else {
		// Update existing user
		_, err := s.db.Exec(`
			UPDATE users SET name = ?, picture = ?, updated_at = CURRENT_TIMESTAMP 
			WHERE id = ?`,
			userInfo.Name, userInfo.Picture, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		user.Name = userInfo.Name
		user.Picture = userInfo.Picture
		user.UpdatedAt = time.Now()
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(userID int) (*User, error) {
	var user User
	err := s.db.QueryRow(`
		SELECT id, google_id, email, name, picture, created_at, updated_at 
		FROM users WHERE id = ?`, userID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Middleware returns a Gin middleware that requires authentication
func (s *Service) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie or Authorization header
		token := s.getTokenFromRequest(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := s.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("authenticated", true)

		c.Next()
	}
}

// OptionalMiddleware returns a Gin middleware that optionally checks for authentication
// If authenticated, sets user context. If not, continues without blocking.
func (s *Service) OptionalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie or Authorization header
		token := s.getTokenFromRequest(c)
		if token == "" {
			// No token found, continue without authentication
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// Validate the token
		claims, err := s.ValidateJWT(token)
		if err != nil {
			// Invalid token, continue without authentication
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// Valid token, set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("authenticated", true)

		c.Next()
	}
}

// getTokenFromRequest extracts JWT token from cookie or Authorization header
func (s *Service) getTokenFromRequest(c *gin.Context) string {
	// First, try to get token from cookie (for web UI)
	if token, err := c.Cookie("auth_token"); err == nil && token != "" {
		return token
	}

	// Second, try to get token from Authorization header (for API)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}

func (s *Service) generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
