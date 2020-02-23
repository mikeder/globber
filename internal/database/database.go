package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// New returns a conection to the database.
func New(cfg *mysql.Config) (*sql.DB, error) {
	log.Printf("Connecting to database: %s/%s", cfg.Addr, cfg.DBName)

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
