package main

import "database/sql"

type ShortURL struct {
	ID         int    `json:"id"`
	ShortCode  string `json:"short_code"`
	LongURL    string `json:"long_url"`
	VisitCount int    `json:"visit_count"`
}

func setupDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./urls.db")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_code TEXT UNIQUE,
			long_url TEXT,
			visit_count INTEGER DEFAULT 0
		)
	`)

	if err != nil {
		panic(err)
	}

	return db
}
