package handlers

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/web"
	"github.com/pkg/errors"
)

// New returns an http.Handler with routes to support
// the API for this application.
func New(bs *blog.Store) http.Handler {
	api := api{
		blogStore: bs,
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
	router.Get("/", api.root)
	router.Get("/blog", api.blog)
	router.Get("/blog/post/{postID}", api.blogPost)

	fs := http.FileServer(http.Dir("./static/"))
	router.Get("/static/*", http.HandlerFunc(http.StripPrefix("/static/", fs).ServeHTTP))

	return router
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

type api struct {
	blogStore *blog.Store
	templates *template.Template
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
	log.Printf("loaded templates: %+v", templates.Name())
}

func (a *api) blog(w http.ResponseWriter, r *http.Request) {
	posts, err := a.blogStore.GetPosts()
	if err != nil {
		log.Println(errors.Wrap(err, "getting posts from database"))
	}

	data := struct {
		PageTitle string
		Posts     []blog.Post
	}{
		PageTitle: "blog",
		Posts:     posts,
	}
	web.Render(w, a.templates.Lookup("blog.html"), data)
}

func (a *api) blogPost(w http.ResponseWriter, r *http.Request) {
	idin := chi.URLParam(r, "postID")
	if idin == "" {
		w.Write([]byte("postID not found"))
	}
	id, err := strconv.Atoi(idin)
	post, err := a.blogStore.GetPost(id)
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
		PageTitle string
		Posts     []blog.Post
	}{
		PageTitle: "blog",
		Posts:     []blog.Post{*post},
	}
	web.Render(w, a.templates.Lookup("blog.html"), data)
}

func (a *api) root(w http.ResponseWriter, r *http.Request) {
	web.Render(w, a.templates.Lookup("home.html"), nil)
}
