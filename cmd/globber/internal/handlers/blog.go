package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/web"
	"github.com/pkg/errors"
)

func (s *site) blogArchive(w http.ResponseWriter, r *http.Request) {
	authed := true
	pagetitle := s.config.SiteName + " - Blog"

	posts, err := s.blogStore.GetArchive(r.Context())
	if err != nil {
		log.Println(errors.Wrap(err, "getting archive posts from database"))
	}

	data := struct {
		Authenticated *bool
		PageTitle     *string
		Entries       []blog.Entry
	}{
		Authenticated: &authed,
		PageTitle:     &pagetitle,
		Entries:       posts,
	}
	s.loadTemplates()
	web.Render(w, s.templates.Lookup("archive.html"), data)
}

func (s *site) blogEntry(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	entry, err := s.blogStore.GetEntryBySlug(r.Context(), slug)
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
		Posts         []blog.Entry
	}{
		Authenticated: true,
		PageTitle:     s.config.SiteName + " - Blog",
		Posts:         []blog.Entry{entry},
	}
	s.loadTemplates()
	web.Render(w, s.templates.Lookup("blog.html"), data)
}

func (s *site) blogPage(w http.ResponseWriter, r *http.Request) {
	authed := true
	pagetitle := s.config.SiteName + " - Blog"
	page, ok := r.URL.Query()["page"]
	if !ok {
		page = []string{"0"}
	}

	pageNum, err := strconv.Atoi(page[0])
	if err != nil {
		log.Println(err)
	}

	posts, err := s.blogStore.GetEntriesByPage(r.Context(), pageNum)
	if err != nil {
		log.Println(errors.Wrap(err, "getting posts from database"))
	}

	data := struct {
		Authenticated *bool
		PageTitle     *string
		Entries       []blog.Entry
	}{
		Authenticated: &authed,
		PageTitle:     &pagetitle,
		Entries:       posts,
	}
	s.loadTemplates()
	web.Render(w, s.templates.Lookup("blog.html"), data)
}

func (s *site) root(w http.ResponseWriter, r *http.Request) {
	s.loadTemplates()
	data := struct {
		Authenticated bool
		PageTitle     string
	}{
		Authenticated: false,
		PageTitle:     s.config.SiteName + " - Home",
	}
	web.Render(w, s.templates.Lookup("home.html"), data)
}
