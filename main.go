package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kinde-starter-kits/golang-starter-kit/config"
	"github.com/kinde-starter-kits/golang-starter-kit/handlers"
	"github.com/kinde-starter-kits/golang-starter-kit/middleware"
	"github.com/kinde-starter-kits/golang-starter-kit/session"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration - if invalid, show setup page instead
	configValid := cfg.Validate() == nil

	// Initialize session store
	session.InitStore(cfg.SessionSecret)

	// Create ServeMux (router)
	mux := http.NewServeMux()

	// Initialize handlers
	h := handlers.NewHandler(cfg)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	if !configValid {
		// If config is not valid, show setup instructions page
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			h.SetupPage(w, r, cfg)
		})
	} else {
		// Public routes
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			h.Home(w, r)
		})
		mux.HandleFunc("/login", h.Login)
		mux.HandleFunc("/register", h.Register)
		mux.HandleFunc("/callback", h.Callback)
		mux.HandleFunc("/logout", h.Logout)

		// Protected routes
		mux.Handle("/dashboard", middleware.AuthRequired(http.HandlerFunc(h.Dashboard)))
		mux.Handle("/profile", middleware.AuthRequired(http.HandlerFunc(h.Profile)))

		// API routes (for demonstration)
		mux.Handle("/api/user", middleware.AuthRequired(http.HandlerFunc(h.GetUser)))
	}

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Start server
	port := cfg.Port
	addr := ":" + port

	log.Printf("🚀 Kinde Golang Starter Kit running on http://localhost:%s", port)
	if !configValid {
		log.Printf("⚠️  Configuration incomplete - visit http://localhost:%s for setup instructions", port)
	} else {
		log.Printf("📝 Make sure you've configured your Kinde app with redirect URI: %s", cfg.RedirectURI)
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

