package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gosimple/slug"
	"github.com/mikeder/globber/internal/auth"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/models"
	"github.com/mikeder/globber/internal/web"
	"github.com/pkg/errors"
)

const maxMemory int64 = 1 * 1.049e+6 // 1MB

type siteData struct {
	Authenticated bool
	SiteName      string
	Username      string
}

type blogPageData struct {
	siteData
	Entries []blog.Entry
}

type blogComposeData struct {
	siteData
	Entry *blog.Entry
}

func (s *site) blogArchive(w http.ResponseWriter, r *http.Request) {
	entries, err := s.blogStore.GetArchive(r.Context())
	if err != nil {
		log.Println(errors.Wrap(err, "getting archive posts from database"))
	}

	validTkn, username := auth.ValidateCtx(r.Context())

	sd := siteData{validTkn, s.config.SiteName, username}
	data := blogPageData{sd, entries}

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

	validTkn, username := auth.ValidateCtx(r.Context())

	sd := siteData{validTkn, s.config.SiteName, username}
	data := blogPageData{sd, []blog.Entry{entry}}

	s.loadTemplates()
	web.Render(w, s.templates.Lookup("blog.html"), data)
}

func (s *site) blogCompose(w http.ResponseWriter, r *http.Request) {
	entryID := r.URL.Query().Get("id")
	var entry *blog.Entry

	if entryID != "" {
		id, err := strconv.Atoi(entryID)
		if err != nil {
			log.Print(err)
			return
		}
		entry, err = s.blogStore.GetEntryByID(r.Context(), int(id))
		if err != nil && err != sql.ErrNoRows {
			log.Print(err)
			return
		}
	}

	if r.Method == http.MethodPost {
		if entry == nil {
			entry = newEntry(r)
		} else {
			updateEntry(entry, r)
		}
		if err := s.blogStore.PostEntry(r.Context(), entry); err != nil {
			log.Println(err)
		}

		validTkn, username := auth.ValidateCtx(r.Context())
		sd := siteData{validTkn, s.config.SiteName, username}
		data := blogPageData{sd, []blog.Entry{*entry}}

		s.loadTemplates()
		web.Render(w, s.templates.Lookup("blog.html"), data)
		return
	}

	validTkn, username := auth.ValidateCtx(r.Context())
	sd := siteData{validTkn, s.config.SiteName, username}
	data := blogComposeData{sd, entry}

	s.loadTemplates()
	web.Render(w, s.templates.Lookup("compose.html"), data)
}

func newEntry(r *http.Request) *blog.Entry {
	// load form data and pass to store
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Print(err)
		return nil
	}

	title := r.Form.Get("title")
	html := r.Form.Get("editordata")

	entry := &models.Entry{
		AuthorID: 0,
		Slug:     slug.Make(title),
		Title:    title,
		HTML:     html,
		// Markdown: md,
	}

	return &blog.Entry{entry}
}

func updateEntry(e *blog.Entry, r *http.Request) {
	// load form data and pass to store
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Print(err)
		return
	}

	title := r.Form.Get("title")
	html := r.Form.Get("editordata")

	e.Title = title
	e.HTML = html
}

func (s *site) blogPage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		log.Println(err)
	}

	entries, err := s.blogStore.GetEntriesByPage(r.Context(), pageNum)
	if err != nil {
		log.Println(errors.Wrap(err, "getting posts from database"))
	}

	validTkn, username := auth.ValidateCtx(r.Context())
	sd := siteData{validTkn, s.config.SiteName, username}
	data := blogPageData{sd, entries}

	s.loadTemplates()
	web.Render(w, s.templates.Lookup("blog.html"), data)
}

func (s *site) root(w http.ResponseWriter, r *http.Request) {
	validTkn, username := auth.ValidateCtx(r.Context())
	sd := siteData{validTkn, s.config.SiteName, username}

	s.loadTemplates()
	web.Render(w, s.templates.Lookup("home.html"), sd)
}
