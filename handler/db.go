package handler

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "history.db?cache=shared&mode=rwc&_journal_mode=WAL&_busy_timeout=10000&_txlock=immediate&_sync=1")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(0)

	return db, nil
}
