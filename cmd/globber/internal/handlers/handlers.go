package handlers

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/web"
	"github.com/pkg/errors"
)

// Config contains contextual information for use within handlers.
type Config struct {
	SiteName string
}

// New returns an http.Handler with routes to support
// the API for this application.
func New(bs *blog.Store, cfg *Config) http.Handler {
	api := api{
		blogStore: bs,
		config:    cfg,
	}
	api.loadTemplates()

	router := chi.NewRouter()

	// add middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// add routes
	router.Get("/*", api.root)
	router.Get("/", api.root)
	router.Get("/blog", api.blog)
	router.Get("/blog/entry/{slug}", api.blogPost)

	router.Get("/favicon.ico", faviconHandler)

	fs := http.FileServer(http.Dir("./static/"))
	router.Get("/static/*", http.HandlerFunc(http.StripPrefix("/static/", fs).ServeHTTP))

	return router
}

type api struct {
	blogStore *blog.Store
	config    *Config
	templates *template.Template
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func (a *api) loadTemplates() {
	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}

	templates, err := template.ParseFiles(allFiles...)
	if err != nil {
		log.Fatal(err)
	}

	templates.Funcs(template.FuncMap{
		"html": func(value string) template.HTML {
			return template.HTML(value)
		},
	})
	a.templates = templates
}

func (a *api) blog(w http.ResponseWriter, r *http.Request) {
	authed := true
	pagetitle := a.config.SiteName + " - Blog"
	posts, err := a.blogStore.GetPosts()
	if err != nil {
		log.Println(errors.Wrap(err, "getting posts from database"))
	}

	data := struct {
		Authenticated *bool
		PageTitle     *string
		Posts         []blog.Post
	}{
		Authenticated: &authed,
		PageTitle:     &pagetitle,
		Posts:         posts,
	}
	a.loadTemplates()
	web.Render(w, a.templates.Lookup("blog.html"), data)
}

func (a *api) blogPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := a.blogStore.GetPostBySlug(slug)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Print(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	data := struct {
		Authenticated bool
		PageTitle     string
		Posts         []blog.Post
	}{
		Authenticated: true,
		PageTitle:     a.config.SiteName + " - Blog",
		Posts:         []blog.Post{*post},
	}
	a.loadTemplates()
	web.Render(w, a.templates.Lookup("blog.html"), data)
}

func (a *api) root(w http.ResponseWriter, r *http.Request) {
	a.loadTemplates()
	data := struct {
		Authenticated bool
		PageTitle     string
	}{
		Authenticated: false,
		PageTitle:     a.config.SiteName + " - Home",
	}
	web.Render(w, a.templates.Lookup("home.html"), data)
}
