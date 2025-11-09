package middleware

import (
	"log"
	"net/http"

	"github.com/kinde-starter-kits/golang-starter-kit/session"
)

// AuthRequired is a middleware that checks if the user is authenticated
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.Get(r)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		userID := sess.Values["user_id"]
		if userID == nil {
			// Save the requested URL to redirect after login
			sess.Values["redirect_after_login"] = r.URL.Path
			sess.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
