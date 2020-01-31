package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func New(cfg *mysql.Config) (*sql.DB, error) {
	log.Printf("Connecting to database: %s/%s", cfg.Addr, cfg.DBName)

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
