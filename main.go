package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func main() {

	// Set a database connection string, this needs improvement.
	defaultDbString := "user:password@tcp(host.domain:3306)/database"
	dbstring := flag.String("dbconnstring", defaultDbString, "database connection string")
	flag.Parse()
	if dbstring == nil {
		log.Printf("You must provide a -dbconstring\n")
		flag.Usage()
		os.Exit(1)
	}

	var db = Database{
		Conn: *dbstring,
		Type: "mysql",
	}

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"html": func(value string) template.HTML {
			return template.HTML(value)
		},
	}).ParseFiles("layout.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.getPosts()
		if err != nil {
			log.Println(err)
		}

		log.Println(r.RemoteAddr, "-", r.URL.Path, "got ", len(posts), " posts.")

		data := struct {
			PageTitle string
			Posts     []Post
		}{
			PageTitle: "Test Blog",
			Posts:     posts,
		}
		tmpl.ExecuteTemplate(w, "layout.html", data)
	})
	http.ListenAndServe(":80", nil)
}

type Database struct {
	Conn string
	Type string
}

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

func (d *Database) getPosts() ([]Post, error) {
	db, err := sql.Open(d.Type, d.Conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM entries ORDER BY published DESC")
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
