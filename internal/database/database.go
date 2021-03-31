package database

import (
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// New returns a conection to the database.
func New(cfg *mysql.Config) (*sqlx.DB, error) {
	log.Printf("Connecting to database: %s/%s", cfg.Addr, cfg.DBName)

	db, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	if err := db.Ping(); err != nil {
		log.Print(err)
		return nil, err
	}

	log.Print("Connected.")

	return db, nil
}
