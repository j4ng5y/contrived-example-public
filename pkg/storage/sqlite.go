package storage

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type (
	Sqlite struct {
		ctx context.Context
		db  *sql.DB
	}
)

func (r *Sqlite) Connect(ctx context.Context, uri string) error {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS resp (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT,
		error TEXT		
	);`); err != nil {
		return err
	}

	r.ctx = ctx
	r.db = db
	return nil
}

func (r *Sqlite) Insert(message, err_message string) error {
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(r.ctx, "INSERT INTO resp (message, error) VALUES (?, ?)", message, err_message)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
