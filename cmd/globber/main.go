package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"

	"github.com/mikeder/globber/cmd/globber/internal/handlers"
	"github.com/mikeder/globber/internal/blog"
	"github.com/mikeder/globber/internal/database"

	_ "net/http/pprof"
)

func main() {
	if err := run(); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}

func run() error {
	log.SetOutput(os.Stdout)
	go func() {
		http.ListenAndServe(":3030", nil)
	}()

	cfg := struct {
		DbUser   string `default:"root" desc:"Username for database connection."`
		DbPass   string `default:"root" desc:"Password for database connection."`
		DbHost   string `default:"db" desc:"Hostname for database connection."`
		DbName   string `default:"blog" desc:"Database schema name."`
		SiteName string `default:"TestBlog" desc:"Name to be used for Title tags."`
	}{}

	if err := envconfig.Process("myapp", &cfg); err != nil {
		return err
	}

	helpFlag := flag.Bool("help", false, "print help info")

	flag.Parse()

	if *helpFlag {
		return envconfig.Usage("", &cfg)
	}

	// Set a database connection string, this needs improvement.
	dbCFG := mysql.NewConfig()

	dbCFG.Addr = cfg.DbHost
	dbCFG.DBName = cfg.DbName
	dbCFG.Net = "tcp"
	dbCFG.Passwd = cfg.DbPass
	dbCFG.User = cfg.DbUser

	db, err := database.New(dbCFG)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	blogStore := blog.New(db)

	handlerCFG := handlers.Config{SiteName: cfg.SiteName}
	handler := handlers.New(blogStore, &handlerCFG)

	server := &http.Server{
		Addr:         ":3000",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
