package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func New(cfg *mysql.Config) *sql.DB {
	log.Printf("Connecting to database: %s/%s", cfg.Addr, cfg.DBName)

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println(err.Error())
	}

	return db
}
