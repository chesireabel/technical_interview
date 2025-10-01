package middleware

import (
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	SessionName   = "oidc-session"
	SessionMaxAge = 86400 * 7 // 7 days
)

// InitSessionStore initializes the session store middleware for Gin
func InitSessionStore() gin.HandlerFunc {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		// fallback secret for local dev (should be set in .env in prod)
		secret = "super-secret-fallback"
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   SessionMaxAge,
		HttpOnly: true,   // prevent JS access
		Secure:   false,  // set true in production (requires HTTPS)
		SameSite: 2,      // Lax mode
	})

	return sessions.Sessions(SessionName, store)
}

// SessionInfo represents current session information
type SessionInfo struct {
	Authenticated bool      `json:"authenticated"`
	UserSub       string    `json:"user_sub,omitempty"`
	UserEmail     string    `json:"user_email,omitempty"`
	UserName      string    `json:"user_name,omitempty"`
	UserPicture   string    `json:"user_picture,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
}
