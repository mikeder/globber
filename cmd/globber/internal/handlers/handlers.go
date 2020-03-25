package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/mikeder/globber/internal/auth"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/minecraft"
)

// Config contains contextual information for use within handlers.
type Config struct {
	SiteName string
}

// New returns an http.Handler with routes to support
// the API for this application.
func New(authMan *auth.Manager, bs *blog.Store, cfg *Config, mc *minecraft.Server) http.Handler {
	adminAPI := adminAPI{authMan}
	authAPI := authAPI{authMan}

	site := site{
		blogStore: bs,
		config:    cfg,
		mc:        mc,
	}
	site.loadTemplates()

	router := chi.NewRouter()

	// add global middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Seek, verify and validate JWT tokens for all requests
	router.Use(jwtauth.Verifier(authMan.TokenAuth))

	// Public routes
	router.Group(func(r chi.Router) {
		// auth handlers
		router.Post("/auth/login", authAPI.Login)
		router.Post("/auth/logout", authAPI.Logout)
		router.Post("/auth/refresh", authAPI.Refresh)

		router.Get("/*", site.root)
		router.Get("/blog", site.blogPage)
		router.Get("/blog/archive", site.blogArchive)
		router.Get("/blog/entry/{slug}", site.blogEntry)

		router.Get("/favicon.ico", faviconHandler)

		router.Get("/minecraft", site.minecraft)
		router.Get("/minecraft/ping", site.minecraftPing)
		router.Get("/minecraft/status", site.minecraftStatus)

		fs := http.FileServer(http.Dir("./static/"))
		router.Get("/static/*", http.HandlerFunc(http.StripPrefix("/static/", fs).ServeHTTP))
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		// Handle valid / invalid tokens
		r.Use(jwtauth.Authenticator)

		r.Post("/admin/user/add", adminAPI.AddUser)

		r.Get("/blog/compose", site.blogCompose)
		r.Post("/blog/compose", site.blogCompose)

	})

	return router
}

// adminAPI contains admin level handlers
type adminAPI struct {
	manager *auth.Manager
}

// auth contains authentication/authorization handlers
type authAPI struct {
	manager *auth.Manager
}

// site contains handlers for rendering page templates
type site struct {
	blogStore *blog.Store
	config    *Config
	templates *template.Template
	mc        *minecraft.Server
}
