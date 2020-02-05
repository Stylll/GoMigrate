package main

import (
	"database/sql"
)

func UP(db *sql.DB) error {
	var query = `CREATE TABLE IF NOT EXISTS MIGRATIONS_TABLE_NAME (
		id SERIAL PRIMARY KEY,
		migration_name VARCHAR UNIQUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		deleted_at TIMESTAMP DEFAULT NULL
	)`

	_, err := db.Exec(query)

	return err
}

func DOWN(db *sql.DB) error {
	var query = `DROP TABLE IF EXISTS MIGRATIONS_TABLE_NAME;`

	_, err := db.Exec(query)

	return err
}
