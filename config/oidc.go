package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/chesireabel/Technical-Interview/internal/middleware"
)

// OIDCInitConfig holds optional configuration for OIDC initialization
type OIDCInitConfig struct {
	// InitTimeout specifies how long to wait for OIDC provider discovery
	InitTimeout time.Duration
	// CustomScopes allows overriding default scopes
	CustomScopes []string
}

// DefaultOIDCInitConfig returns sensible defaults
func DefaultOIDCInitConfig() OIDCInitConfig {
	return OIDCInitConfig{
		InitTimeout:  30 * time.Second,
		CustomScopes: []string{"openid", "profile", "email"},
	}
}

// InitOIDC initializes and returns an OIDC instance with improved error handling and validation
func InitOIDC(ctx context.Context, cfg OIDCInitConfig) (*middleware.OIDC, error) {
	// Validate required environment variables
	envVars := map[string]string{
		"AUTH0_CLIENT_ID":     os.Getenv("AUTH0_CLIENT_ID"),
		"AUTH0_CLIENT_SECRET": os.Getenv("AUTH0_CLIENT_SECRET"),
		"AUTH0_REDIRECT_URL":  os.Getenv("AUTH0_REDIRECT_URL"),
		"AUTH0_PROVIDER_URL":  os.Getenv("AUTH0_PROVIDER_URL"),
	}

	// Check for missing variables
	var missing []string
	for key, value := range envVars {
		if value == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}



	// Create OIDC configuration
	oidcConfig := middleware.OIDCConfig{
		ClientID:     envVars["AUTH0_CLIENT_ID"],
		ClientSecret: envVars["AUTH0_CLIENT_SECRET"],
		RedirectURL:  envVars["AUTH0_REDIRECT_URL"],
		Issuer:       envVars["AUTH0_PROVIDER_URL"],
		Scopes:       cfg.CustomScopes,
	}

	// Create OIDC instance with timeout handling
	type result struct {
		oidc *middleware.OIDC
		err  error
	}
	
	resultCh := make(chan result, 1)
	
	// Run OIDC initialization in goroutine
	go func() {
		// Call NewOIDC with just the config
        oidc, err := middleware.NewOIDCWithConfig(ctx, oidcConfig)
		resultCh <- result{oidc: oidc, err: err}
	}()

	// Wait for initialization or timeout
	initCtx, cancel := context.WithTimeout(ctx, cfg.InitTimeout)
	defer cancel()

	select {
	case res := <-resultCh:
		if res.err != nil {
			// Don't expose sensitive details in error
			log.Printf("❌ OIDC initialization failed: %v", sanitizeError(res.err))
			return nil, fmt.Errorf("failed to initialize OIDC: %w", res.err)
		}
		// Log success without exposing sensitive data
		log.Printf("✅ OIDC authentication initialized (issuer: %s)", maskURL(envVars["AUTH0_PROVIDER_URL"]))
		return res.oidc, nil
	case <-initCtx.Done():
		log.Printf("❌ OIDC initialization timed out after %v", cfg.InitTimeout)
		return nil, fmt.Errorf("OIDC initialization timed out after %v", cfg.InitTimeout)
	}
}

// InitOIDCWithDefaults is a convenience wrapper using default configuration
func InitOIDCWithDefaults(ctx context.Context) (*middleware.OIDC, error) {
	return InitOIDC(ctx, DefaultOIDCInitConfig())
}

// GetReturnURL returns the application URL for post-logout redirect
func GetReturnURL() string {
	returnToURL := os.Getenv("APP_URL")
	if returnToURL == "" {
		log.Println("⚠️ APP_URL not set")
		return returnToURL
	}

	

	return returnToURL
}



// maskURL masks sensitive parts of a URL for logging
func maskURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "[invalid-url]"
	}
	
	// Show scheme and host only
	return fmt.Sprintf("%s://%s/...", parsed.Scheme, parsed.Host)
}

// sanitizeError removes potentially sensitive information from errors
func sanitizeError(err error) error {
	if err == nil {
		return nil
	}
	// In production, you might want more sophisticated sanitization
	// For now, just return a generic wrapper
	return fmt.Errorf("authentication configuration error")
}

// ValidateOIDCConfig checks if all required OIDC environment variables are set
// This can be called at startup before attempting initialization
func ValidateOIDCConfig() error {
	required := []string{
		"AUTH0_CLIENT_ID",
		"AUTH0_CLIENT_SECRET",
		"AUTH0_REDIRECT_URL",
		"AUTH0_PROVIDER_URL",
	}

	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required OIDC environment variables: %v", missing)
	}

	return nil
}