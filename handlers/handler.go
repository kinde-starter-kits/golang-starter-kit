package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kinde-starter-kits/golang-starter-kit/config"
	"github.com/kinde-starter-kits/golang-starter-kit/session"
	"golang.org/x/oauth2"
)

// Handler holds the application handlers
type Handler struct {
	config      *Config
	oauthConfig *oauth2.Config
	templates   *template.Template
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config) *Handler {
	// Load all templates
	templates := template.Must(template.ParseGlob("templates/*"))

	return &Handler{
		config: &Config{
			KindeDomain:       cfg.KindeDomain,
			KindeClientID:     cfg.KindeClientID,
			KindeClientSecret: cfg.KindeClientSecret,
			RedirectURI:       cfg.RedirectURI,
			LogoutRedirectURI: cfg.LogoutRedirectURI,
		},
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.KindeClientID,
			ClientSecret: cfg.KindeClientSecret,
			RedirectURL:  cfg.RedirectURI,
			Scopes:       []string{"openid", "profile", "email", "offline"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("%s/oauth2/auth", cfg.KindeDomain),
				TokenURL: fmt.Sprintf("%s/oauth2/token", cfg.KindeDomain),
			},
		},
		templates: templates,
	}
}

// Config is a local config structure for handlers
type Config struct {
	KindeDomain       string
	KindeClientID     string
	KindeClientSecret string
	RedirectURI       string
	LogoutRedirectURI string
}

// Home renders the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get(r)

	var userID interface{}
	var userEmail interface{}

	if err != nil {
		log.Printf("Error getting session (will show as not authenticated): %v", err)
		// Clear bad cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "kinde_session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		userID = nil
		userEmail = nil
	} else {
		userID = sess.Values["user_id"]
		userEmail = sess.Values["user_email"]
	}

	data := map[string]interface{}{
		"authenticated":     userID != nil,
		"user":              sess.Values["user_name"],
		"user_email":        userEmail,
		"user_first_name":   sess.Values["user_first_name"],
		"user_last_name":    sess.Values["user_last_name"],
		"user_initials":     sess.Values["user_initials"],
		"user_picture":      sess.Values["user_picture"],
	}

	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Login initiates the OAuth2 login flow
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== LOGIN START ===")
	log.Printf("Request Host: %s", r.Host)
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Cookies received: %v", r.Cookies())

	// Create a fresh session (don't use existing one to avoid cookie issues)
	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Session error, creating new session: %v", err)
		// Clear the bad cookie by setting MaxAge to -1
		http.SetCookie(w, &http.Cookie{
			Name:   "kinde_session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		// Get a fresh session
		sess, err = session.Get(r)
		if err != nil {
			log.Printf("Failed to create new session: %v", err)
			h.renderError(w, "Session error. Please try again or use an incognito window.")
			return
		}
	}

	// Clear any existing OAuth state from previous attempts
	delete(sess.Values, "oauth_state")
	delete(sess.Values, "code_verifier")

	// Generate state parameter for CSRF protection
	state := generateRandomString(32)
	sess.Values["oauth_state"] = state

	// Generate PKCE code verifier and challenge
	codeVerifier := generateRandomString(64)
	sess.Values["code_verifier"] = codeVerifier

	log.Printf("Login: Generated state=%s, codeVerifier=%s", state, codeVerifier[:10]+"...")
	log.Printf("Login: Session ID before save: %v", sess.ID)

	if err := sess.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		h.renderError(w, "Failed to initiate login")
		return
	}

	log.Printf("Login: Session saved successfully")
	log.Printf("Login: Session values: %+v", sess.Values)

	// Build authorization URL with PKCE
	authURL := h.oauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge(codeVerifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	log.Printf("Login: Redirecting to: %s", authURL)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Register initiates the OAuth2 registration flow
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		h.renderError(w, "Failed to initiate registration")
		return
	}

	// Generate state parameter for CSRF protection
	state := generateRandomString(32)
	sess.Values["oauth_state"] = state

	// Generate PKCE code verifier
	codeVerifier := generateRandomString(64)
	sess.Values["code_verifier"] = codeVerifier

	if err := sess.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		h.renderError(w, "Failed to initiate registration")
		return
	}

	// Build authorization URL with screen_hint=registration
	authURL := h.oauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge(codeVerifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("screen_hint", "registration"),
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Callback handles the OAuth2 callback
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== CALLBACK START ===")
	log.Printf("Request Host: %s", r.Host)
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Cookies received: %v", r.Cookies())

	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Error getting session in callback: %v", err)
		h.renderError(w, "Session error. Please try logging in again.")
		return
	}

	log.Printf("Callback: Session ID: %v", sess.ID)
	log.Printf("Callback: Session IsNew: %v", sess.IsNew)
	log.Printf("Callback: All session values: %+v", sess.Values)

	// Check if user is already logged in (callback was already processed)
	if sess.Values["user_id"] != nil {
		log.Printf("Callback: User already logged in, redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}

	// Verify state parameter
	state := r.URL.Query().Get("state")
	savedState, ok := sess.Values["oauth_state"].(string)

	// Debug logging
	log.Printf("Callback: state from URL: %s", state)
	log.Printf("Callback: state from session: %s, ok: %v", savedState, ok)

	if !ok {
		log.Printf("Callback: State not found - Session values: %+v", sess.Values)
		h.renderError(w, "State not found in session. Your session may have expired. Please try logging in again.")
		return
	}

	if state == "" {
		h.renderError(w, "No state parameter received from OAuth provider")
		return
	}

	if state != savedState {
		log.Printf("State mismatch - expected: %s, got: %s", savedState, state)
		h.renderError(w, "Invalid state parameter (CSRF check failed). Please try logging in again.")
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		h.renderError(w, "No authorization code received")
		return
	}

	// Exchange code for token with PKCE
	codeVerifier, ok := sess.Values["code_verifier"].(string)
	if !ok {
		h.renderError(w, "Code verifier not found in session")
		return
	}

	log.Printf("Callback: Exchanging code for token...")
	token, err := h.oauthConfig.Exchange(
		context.Background(),
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		log.Printf("Token exchange error: %v", err)
		h.renderError(w, fmt.Sprintf("Failed to exchange token: %v", err))
		return
	}
	log.Printf("Callback: Token exchange successful")

	// Fetch user information from Kinde's userinfo endpoint
	log.Printf("Callback: Fetching user info...")
	userInfo, err := h.fetchUserInfo(context.Background(), token.AccessToken)
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		h.renderError(w, fmt.Sprintf("Failed to fetch user information: %v", err))
		return
	}
	log.Printf("Callback: User info fetched successfully: %+v", userInfo)

	// Store only essential user information in session (not the large tokens)
	// Kinde returns 'id' not 'sub'
	if id, ok := userInfo["id"].(string); ok {
		sess.Values["user_id"] = id
	}
	if email, ok := userInfo["preferred_email"].(string); ok {
		sess.Values["user_email"] = email
	}

	// Get first and last name
	var fullName string
	firstName := ""
	lastName := ""
	if fn, ok := userInfo["given_name"].(string); ok && fn != "" {
		firstName = fn
		fullName = fn
	} else if fn, ok := userInfo["first_name"].(string); ok && fn != "" {
		firstName = fn
		fullName = fn
	}
	if ln, ok := userInfo["family_name"].(string); ok && ln != "" {
		lastName = ln
		if fullName != "" {
			fullName += " " + ln
		} else {
			fullName = ln
		}
	} else if ln, ok := userInfo["last_name"].(string); ok && ln != "" {
		lastName = ln
		if fullName != "" {
			fullName += " " + ln
		} else {
			fullName = ln
		}
	}

	sess.Values["user_name"] = fullName
	sess.Values["user_first_name"] = firstName
	sess.Values["user_last_name"] = lastName

	// Compute initials for avatar
	initials := ""
	if firstName != "" && len(firstName) > 0 {
		initials += string(firstName[0])
	}
	if lastName != "" && len(lastName) > 0 {
		initials += string(lastName[0])
	}
	sess.Values["user_initials"] = initials

	if picture, ok := userInfo["picture"].(string); ok {
		sess.Values["user_picture"] = picture
	}

	log.Printf("Callback: Stored user data in session: user_id=%v, email=%v", sess.Values["user_id"], sess.Values["user_email"])

	// Clear OAuth state and code verifier
	delete(sess.Values, "oauth_state")
	delete(sess.Values, "code_verifier")

	log.Printf("Callback: Saving session with user data...")
	if err := sess.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		h.renderError(w, "Failed to save session")
		return
	}

	log.Printf("Callback: Login successful! Redirecting to dashboard")
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

// Logout handles user logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	// Clear session
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		log.Printf("Error clearing session: %v", err)
	}

	// Redirect to Kinde logout endpoint
	logoutURL := fmt.Sprintf("%s/logout?redirect=%s", h.config.KindeDomain, h.config.LogoutRedirectURI)
	http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
}

// Dashboard renders the dashboard page (protected)
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	data := map[string]interface{}{
		"user_id":         sess.Values["user_id"],
		"user_name":       sess.Values["user_name"],
		"user_first_name": sess.Values["user_first_name"],
		"user_last_name":  sess.Values["user_last_name"],
		"user_initials":   sess.Values["user_initials"],
		"user_email":      sess.Values["user_email"],
		"user_picture":    sess.Values["user_picture"],
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Profile renders the profile page (protected)
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	data := map[string]interface{}{
		"user_id":      sess.Values["user_id"],
		"user_name":    sess.Values["user_name"],
		"user_email":   sess.Values["user_email"],
		"user_picture": sess.Values["user_picture"],
	}

	if err := h.templates.ExecuteTemplate(w, "profile.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetUser returns user information as JSON (protected)
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get(r)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      sess.Values["user_id"],
		"name":    sess.Values["user_name"],
		"email":   sess.Values["user_email"],
		"picture": sess.Values["user_picture"],
	})
}

// SetupPage renders the setup instructions page
func (h *Handler) SetupPage(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	data := map[string]interface{}{
		"env_file":       ".env",
		"current_domain": cfg.KindeDomain,
	}

	if err := h.templates.ExecuteTemplate(w, "setup.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Helper functions

func (h *Handler) renderError(w http.ResponseWriter, errMsg string) {
	data := map[string]interface{}{
		"error": errMsg,
	}
	w.WriteHeader(http.StatusInternalServerError)
	if err := h.templates.ExecuteTemplate(w, "error.html", data); err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
	}
}

// fetchUserInfo fetches user information from Kinde's userinfo endpoint
func (h *Handler) fetchUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/oauth2/user_profile", h.config.KindeDomain), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status: %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}

// codeChallenge generates a PKCE code challenge from a code verifier
func codeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
