package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/mikeder/globber/cmd/globber/internal/handlers"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Set a database connection string, this needs improvement.
	dbuser := flag.String("dbuser", "", "database username")
	dbpass := flag.String("dbpass", "", "database password")
	dbhost := flag.String("dbhost", "", "database hostname")
	dbname := flag.String("dbname", "", "database name")
	sitename := flag.String("sitename", "", "website name used in titles")

	flag.Parse()

	if *sitename == "" {
		*sitename = "Test Site"
	}

	dbCFG := mysql.NewConfig()

	dbCFG.User = *dbuser
	dbCFG.Passwd = *dbpass
	dbCFG.Net = "tcp"
	dbCFG.Addr = *dbhost
	dbCFG.DBName = *dbname

	db := database.New(dbCFG)

	blogStore := blog.New(db)

	server := &http.Server{
		Addr:         ":3000",
		Handler:      handlers.New(blogStore, &handlers.Config{SiteName: *sitename}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
