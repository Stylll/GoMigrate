package main

import (
	"database/sql"
)

/**
IGNORE THE ERROR STATING THAT UP OR DOWN IS REDECLARED. AT RUN TIME,
THE GO SCRIPTS WILL RUN INDIVIDUALLY SO THE COMPILER WON'T THROW AN ERROR
*/

func UP(db *sql.DB) error {
	var query = `CREATE TABLE IF NOT EXISTS SAMPLE_TABLE_NAME (
		id SERIAL PRIMARY KEY,
		sample_column VARCHAR UNIQUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL
	)`

	_, err := db.Exec(query)

	return err
}
