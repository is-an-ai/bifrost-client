package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"bifrost-client/internal/platform"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrTokenExpired     = errors.New("token expired")
)

// TokenPayload represents the JWT payload
type TokenPayload struct {
	Exp int64 `json:"exp"`
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	// GetAuthURL returns the GitHub OAuth URL
	GetAuthURL(ctx context.Context) (string, error)

	// HandleCallback processes the token from the callback URL
	HandleCallback(ctx context.Context, callbackURL string) error

	// IsAuthenticated checks if the user is authenticated
	IsAuthenticated(ctx context.Context) (bool, error)

	// Logout performs the logout operation
	Logout(ctx context.Context) error

	// StartLogin initiates the login process
	StartLogin(ctx context.Context) error

	// CheckAndStartLogin checks authentication and starts login if needed
	CheckAndStartLogin(ctx context.Context) error
}

// Config represents the authentication configuration
type Config struct {
	// API Server configuration
	APIServerURL string
}

// Service implements the AuthService interface
type Service struct {
	config  Config
	storage Storage
}

// NewService creates a new authentication service
func NewService(config Config, storage Storage) *Service {
	return &Service{
		config:  config,
		storage: storage,
	}
}

// GetAuthURL implements AuthService interface
func (s *Service) GetAuthURL(ctx context.Context) (string, error) {
	baseURL := fmt.Sprintf("%s/v1/user/auth/github", s.config.APIServerURL)
	params := url.Values{}
	params.Add("client_type", "bifrost-client")

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// HandleCallback implements AuthService interface
func (s *Service) HandleCallback(ctx context.Context, callbackURL string) error {
	// Parse the callback URL
	parsedURL, err := url.Parse(callbackURL)
	if err != nil {
		return fmt.Errorf("failed to parse callback URL: %w", err)
	}

	// Extract token from query parameters
	token := parsedURL.Query().Get("token")
	if token == "" {
		return fmt.Errorf("no token found in callback URL")
	}

	// Save the token
	if err := s.storage.SaveToken(ctx, token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// IsAuthenticated implements AuthService interface
func (s *Service) IsAuthenticated(ctx context.Context) (bool, error) {
	token, err := s.storage.GetToken(ctx)
	if err != nil {
		return false, err
	}
	if token == "" {
		return false, nil
	}

	// Split the token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}

	// Decode the payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Parse the payload
	var tokenPayload TokenPayload
	if err := json.Unmarshal(payload, &tokenPayload); err != nil {
		return false, fmt.Errorf("failed to parse token payload: %w", err)
	}

	// Check if the token is expired
	if tokenPayload.Exp < time.Now().Unix() {
		// Token is expired, delete it
		if err := s.storage.DeleteToken(ctx); err != nil {
			return false, fmt.Errorf("failed to delete expired token: %w", err)
		}
		return false, ErrTokenExpired
	}

	return true, nil
}

// Logout implements AuthService interface
func (s *Service) Logout(ctx context.Context) error {
	return s.storage.DeleteToken(ctx)
}

// StartLogin implements AuthService interface
func (s *Service) StartLogin(ctx context.Context) error {
	// Check if already authenticated
	authenticated, err := s.IsAuthenticated(ctx)
	if err != nil {
		return fmt.Errorf("failed to check authentication status: %w", err)
	}
	if authenticated {
		return nil // Already authenticated
	}

	// Get the GitHub OAuth URL
	authURL, err := s.GetAuthURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth URL: %w", err)
	}

	// Open the URL in the default browser
	if err := platform.OpenBrowser(authURL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}

// CheckAndStartLogin implements AuthService interface
func (s *Service) CheckAndStartLogin(ctx context.Context) error {
	// Check if already authenticated
	authenticated, err := s.IsAuthenticated(ctx)
	if err != nil {
		return fmt.Errorf("failed to check authentication status: %w", err)
	}

	// If not authenticated, start login process
	if !authenticated {
		return s.StartLogin(ctx)
	}

	return nil
}

func (s *Service) GetAuthToken(ctx context.Context) (string, error) {
	token, err := s.storage.GetToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return token, nil
}
