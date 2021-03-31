package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Migration is used to apply schema patches to a database.
type Migration struct {
	author  string
	query   string
	version float32
}

// Migrate runs migrations on the provide db connection
func Migrate(ctx context.Context, db *sqlx.DB) error {
	patched := false
	current := currentVersion(ctx, db)

	log.Printf("Current schema version: %.2f", current)

	for _, patch := range Migrations {
		if patch.version > current {
			log.Printf("Applying patch version %.2f by %s", patch.version, patch.author)

			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				return err
			}
			defer func() {
				if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
					log.Println(err)
				}
			}()

			_, err = tx.ExecContext(ctx, patch.query)
			if err != nil {
				return logAndRollback(err, tx)
			}
			if err := tx.Commit(); err != nil {
				return logAndRollback(err, tx)
			}

			// finally update version table
			if err := updateVersion(ctx, db, patch.author, patch.version); err != nil {
				return err
			}
			patched = true
		}
	}
	if !patched {
		log.Println("No patches needed, database up-to-date.")
	}

	return nil
}

func currentVersion(ctx context.Context, db *sqlx.DB) float32 {
	const q = `SELECT version from version ORDER BY id DESC LIMIT 0, 1`
	var version float32 = -1.0

	row, err := db.QueryContext(ctx, q)
	if err != nil {
		log.Println(err)
		return version
	}

	if row.Err() != nil {
		log.Println(err)
		return version
	}

	for row.Next() {
		if err := row.Scan(&version); err != nil {
			log.Println(err)
			return version
		}
	}
	return version
}

func updateVersion(ctx context.Context, db *sqlx.DB, author string, version float32) error {
	const q = `INSERT INTO version(applied, author, version) VALUES (?, ?, ?);`

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Println(err)
		}
	}()

	_, err = tx.ExecContext(ctx, q, time.Now(), author, version)
	if err != nil {
		return logAndRollback(err, tx)
	}

	if err := tx.Commit(); err != nil {
		return logAndRollback(err, tx)
	}

	return nil
}

func logAndRollback(err error, tx *sql.Tx) error {
	log.Println(err)
	return tx.Rollback()
}
