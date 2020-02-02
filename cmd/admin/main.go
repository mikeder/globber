package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"

	"github.com/mikeder/globber/internal/database"
)

func main() {
	if err := run(); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}

func run() error {
	log.SetOutput(os.Stdout)

	cfg := struct {
		DbUser string `default:"root" desc:"Username for database connection."`
		DbPass string `default:"root" desc:"Password for database connection."`
		DbHost string `default:"db" desc:"Hostname for database connection."`
		DbName string `default:"blog" desc:"Database schema name."`
	}{}

	if err := envconfig.Process("myapp", &cfg); err != nil {
		return err
	}

	helpFlag := flag.Bool("help", false, "print help info")
	migrate := flag.Bool("migrate", false, "perform database migrations on startup.")

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

	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(30*time.Second))
	defer cancel()

	if *migrate {
		if err := database.Migrate(ctx, db); err != nil {
			cancel()
			return err
		}
	}

	return nil
}
