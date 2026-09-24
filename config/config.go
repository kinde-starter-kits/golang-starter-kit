package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	KindeDomain       string
	KindeClientID     string
	KindeClientSecret string
	RedirectURI       string
	LogoutRedirectURI string
	Port              string
	SessionSecret     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "change-this-secret-in-production"
	}

	return &Config{
		KindeDomain:       os.Getenv("KINDE_DOMAIN"),
		KindeClientID:     os.Getenv("KINDE_CLIENT_ID"),
		KindeClientSecret: os.Getenv("KINDE_CLIENT_SECRET"),
		RedirectURI:       os.Getenv("KINDE_REDIRECT_URI"),
		LogoutRedirectURI: os.Getenv("KINDE_LOGOUT_REDIRECT_URI"),
		Port:              port,
		SessionSecret:     sessionSecret,
	}, nil
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	if c.KindeDomain == "" {
		return fmt.Errorf("KINDE_DOMAIN is required")
	}
	if c.KindeClientID == "" {
		return fmt.Errorf("KINDE_CLIENT_ID is required")
	}
	if c.KindeClientSecret == "" {
		return fmt.Errorf("KINDE_CLIENT_SECRET is required")
	}
	if c.RedirectURI == "" {
		return fmt.Errorf("KINDE_REDIRECT_URI is required")
	}
	if c.LogoutRedirectURI == "" {
		return fmt.Errorf("KINDE_LOGOUT_REDIRECT_URI is required")
	}
	return nil
}


