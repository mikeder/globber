package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/mikeder/globber/internal/blog"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

// Config contains contextual information for use within handlers.
type Config struct {
	SiteName string
}

// New returns an http.Handler with routes to support
// the API for this application.
func New(bs *blog.Store, cfg *Config) http.Handler {
	site := site{
		blogStore: bs,
		config:    cfg,
	}
	site.loadTemplates()

	router := chi.NewRouter()

	// add global middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Public routes
	router.Group(func(r chi.Router) {
		// auth handlers
		// router.Post("/auth/login", api.auth.login)
		// router.Post("/auth/logout", api.auth.logout)
	
		router.Get("/*", site.root)
		router.Get("/blog", site.blogPage)
		router.Get("/blog/archive", site.blogArchive)
		router.Get("/blog/entry/{slug}", site.blogEntry)

		router.Get("/favicon.ico", faviconHandler)

		fs := http.FileServer(http.Dir("./static/"))
		router.Get("/static/*", http.HandlerFunc(http.StripPrefix("/static/", fs).ServeHTTP))
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})

	return router
}

// api contains json/rest handlers
type api struct {
	auth
}

// auth contains authentication/authorization handlers
type auth struct{}

// site contains handlers for rendering page templates
type site struct {
	blogStore *blog.Store
	config    *Config
	templates *template.Template
}
