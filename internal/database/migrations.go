package database

// Migrations are applied to a database when requested.
var Migrations = []Migration{
	Migration{
		author: "mikeder",
		query: `CREATE TABLE IF NOT EXISTS version (
			id int(11) NOT NULL AUTO_INCREMENT,
			applied TIME NOT NULL,
			author varchar(25) NOT NULL,
			version float(4,2) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8`,
		version: 0.1,
	},
	Migration{
		author: "mikeder",
		query: `CREATE TABLE IF NOT EXISTS authors (
			id int(11) NOT NULL AUTO_INCREMENT,
			email varchar(100) NOT NULL,
			name varchar(100) NOT NULL,
			hashed_password varchar(100) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY email (email)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8`,
		version: 0.2,
	},
	Migration{
		author: "mikeder",
		query: `CREATE TABLE IF NOT EXISTS entries (
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
		  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8`,
		version: 0.3,
	},
	Migration{
		author: "mikeder",
		query: `ALTER TABLE entries 
			ADD highlight int(8)`,
		version: 0.4,
	},
}
