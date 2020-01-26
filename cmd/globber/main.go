package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func noescape(str string) template.HTML {
	return template.HTML(str)
}

type API struct {
	db *database
}

func (a *API) Root(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"html": func(value string) template.HTML {
			return template.HTML(value)
		},
	}).ParseFiles("templates/layout.html"))

	posts, err := a.db.getPosts()
	if err != nil {
		log.Println(err)
	}

	data := struct {
		PageTitle string
		Posts     []Post
	}{
		PageTitle: "Test Blog",
		Posts:     posts,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func main() {
	// Set a database connection string, this needs improvement.
	dbuser := flag.String("dbuser", "", "database username")
	dbpass := flag.String("dbpass", "", "database password")
	dbhost := flag.String("dbhost", "", "database hostname")
	dbname := flag.String("dbname", "", "database name")

	flag.Parse()

	config := mysql.NewConfig()

	config.User = *dbuser
	config.Passwd = *dbpass
	config.Net = "tcp"
	config.Addr = *dbhost
	config.DBName = *dbname

	api := API{
		db: newDB(config),
	}

	http.Handle("/", wrapHandler(http.HandlerFunc(api.Root)))

	http.ListenAndServe(":3000", nil)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func wrapHandler(wrapped http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lrw := NewLoggingResponseWriter(w)
		wrapped.ServeHTTP(lrw, r)

		log.Printf("Requst handled: %s %s - %d", r.RemoteAddr, r.URL.Path, lrw.statusCode)

	})
}

func newDB(cfg *mysql.Config) *database {
	log.Printf("Connecting to database: %s/%s", cfg.Addr, cfg.DBName)

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println(err.Error())
	}

	return &database{db}
}

type database struct {
	*sql.DB
}

// Post contains the fields for a blog post.
type Post struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Markdown  string `json:"markdown"`
	HTML      string `json:"html"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

func (d *database) getPosts() ([]Post, error) {
	results, err := d.Query("SELECT * FROM entries ORDER BY published DESC")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var posts []Post
	for results.Next() {
		var post Post
		err = results.Scan(
			&post.ID,
			&post.Author,
			&post.Slug,
			&post.Title,
			&post.Markdown,
			&post.HTML,
			&post.Published,
			&post.Updated,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
