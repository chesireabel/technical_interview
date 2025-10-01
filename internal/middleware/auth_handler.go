package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginHandler initiates the OAuth2 authentication flow
func (o *OIDC) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate state & nonce
		state, err := GenerateState()
		if err != nil {
			log.Printf("Error generating state: %v", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		nonce, err := GenerateNonce()
		if err != nil {
			log.Printf("Error generating nonce: %v", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		// Save in session
		session := sessions.Default(c)
		session.Set("state", state)
		session.Set("nonce", nonce)
		if err := session.Save(); err != nil {
			log.Printf("Error saving session: %v", err)
			c.JSON(500, gin.H{"error": "session error"})
			return
		}

		// Redirect to Auth0 login
		authURL := o.Config.AuthCodeURL(state, oidc.Nonce(nonce))
		c.Redirect(302, authURL)
	}
}

// CallbackHandler handles the OAuth2 callback
func (o *OIDC) CallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// Validate state
		stateFromQuery := c.Query("state")
		stateFromSession, _ := session.Get("state").(string)
		if stateFromQuery != stateFromSession {
			c.JSON(400, gin.H{"error": "invalid state"})
			return
		}

		// Handle auth error
		if errMsg := c.Query("error"); errMsg != "" {
			c.JSON(401, gin.H{"error": errMsg, "description": c.Query("error_description")})
			return
		}

		// Exchange code for token
		code := c.Query("code")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		token, err := o.Config.Exchange(ctx, code)
		if err != nil {
			log.Printf("Token exchange failed: %v", err)
			c.JSON(500, gin.H{"error": "token exchange failed"})
			return
		}

		rawIDToken, _ := token.Extra("id_token").(string)
		idToken, err := o.Verifier.Verify(ctx, rawIDToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid id_token"})
			return
		}

		// Verify nonce
		nonceFromSession, _ := session.Get("nonce").(string)
		if idToken.Nonce != nonceFromSession {
			c.JSON(400, gin.H{"error": "invalid nonce"})
			return
		}

		// Extract claims
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			c.JSON(500, gin.H{"error": "failed to parse claims"})
			return
		}

		// Save user info in session
		session.Set("authenticated", true)
		session.Set("access_token", token.AccessToken)
		session.Set("id_token", rawIDToken)
		session.Set("user_sub", claims["sub"])
		if email, ok := claims["email"].(string); ok {
			session.Set("user_email", email)
		}
		if name, ok := claims["name"].(string); ok {
			session.Set("user_name", name)
		}
		if picture, ok := claims["picture"].(string); ok {
			session.Set("user_picture", picture)
		}

		session.Delete("state")
		session.Delete("nonce")
		session.Save()

		// Redirect to dashboard
		c.Redirect(303, "/")
	}
}

// LogoutHandler
func (o *OIDC) LogoutHandler(returnToURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()

		logoutURL := fmt.Sprintf("%s/v2/logout?client_id=%s&returnTo=%s",
			o.Issuer, o.Config.ClientID, returnToURL)

		c.Redirect(302, logoutURL)
	}
}
