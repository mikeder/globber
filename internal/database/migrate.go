package database

import "database/sql"

import "context"

import "log"

// Migrate runs migrations on the provide db connection
func Migrate(db *sql.DB) error {
	migrations := []string{}

	patch0 := `CREATE TABLE authors (
		id int(11) NOT NULL AUTO_INCREMENT,
		email varchar(100) NOT NULL,
		name varchar(100) NOT NULL,
		hashed_password varchar(100) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY email (email)
	) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8`

	migrations = append(migrations, patch0)

	patch1 := `CREATE TABLE entries (
		id int(11) NOT NULL AUTO_INCREMENT,
		author_id int(11) NOT NULL,
		slug varchar(100) NOT NULL,
		title varchar(512) NOT NULL,
		markdown mediumtext NOT NULL,
		html mediumtext NOT NULL,
		published datetime NOT NULL,
		updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,                               
		PRIMARY KEY (id),
		UNIQUE KEY slug (slug),
		KEY published (published)
	  ) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8`

	migrations = append(migrations, patch1)

	for i, patch := range migrations {
		ctx := context.Background() // TODO: do something useful w/ ctx
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		log.Printf("Applying patch %d", i)
		tx.ExecContext(ctx, patch)
		if err := tx.Commit(); err != nil {
			return tx.Rollback()
		}
	}

	return nil
}
