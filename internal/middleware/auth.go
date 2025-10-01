package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type OIDC struct {
	Config   *oauth2.Config
	Verifier *oidc.IDTokenVerifier
	Issuer   string
}

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Issuer       string
	Scopes       []string
}



// NewOIDC initializes OIDC middleware with configuration from environment
func NewOIDC(ctx context.Context) (*OIDC, error) {
	_ = godotenv.Load()

	config := OIDCConfig{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_REDIRECT_URL"),
		Issuer:       os.Getenv("AUTH0_PROVIDER_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return NewOIDCWithConfig(ctx, config)
}

// NewOIDCWithConfig initializes OIDC middleware with provided configuration
func NewOIDCWithConfig(ctx context.Context, config OIDCConfig) (*OIDC, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC provider at %s: %w", config.Issuer, err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return &OIDC{
		Config:   oauth2Config,
		Verifier: verifier,
		Issuer:   config.Issuer,
	}, nil
}

// validateConfig validates OIDC configuration
func validateConfig(config OIDCConfig) error {
	if config.ClientID == "" {
		return fmt.Errorf("AUTH0_CLIENT_ID is required")
	}
	if config.ClientSecret == "" {
		return fmt.Errorf("AUTH0_CLIENT_SECRET is required")
	}
	if config.RedirectURL == "" {
		return fmt.Errorf("AUTH0_REDIRECT_URL is required")
	}
	if config.Issuer == "" {
		return fmt.Errorf("AUTH0_PROVIDER_URL is required")
	}

	// Validate redirect URL format
	if _, err := url.Parse(config.RedirectURL); err != nil {
		return fmt.Errorf("invalid redirect URL: %w", err)
	}

	// Validate issuer URL format
	if _, err := url.Parse(config.Issuer); err != nil {
		return fmt.Errorf("invalid issuer URL: %w", err)
	}

	return nil
}

// GenerateState creates a cryptographically secure random state parameter
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateNonce creates a cryptographically secure random nonce
func GenerateNonce() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
