package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

// InitStore initializes the session store
func InitStore(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
	store.MaxAge(86400 * 7) // 7 days
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,      // Secure cookie
		Secure:   false,     // Set to true in production with HTTPS
		// Don't set SameSite - let browser use default behavior
		// This is more compatible with OAuth redirects on localhost
	}
}

// Get retrieves a session for the given request
func Get(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "kinde_session")
}

