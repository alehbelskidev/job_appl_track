// Package repository
package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Repo struct {
	db *sql.DB
}

func NewRepo() (*Repo, error) {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		return nil, err
	}

	return &Repo{db: db}, nil
}

func (s *Repo) Close() error {
	return s.db.Close()
}
