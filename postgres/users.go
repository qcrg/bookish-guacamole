package postgres

import (
	"database/sql"
)

type Users struct {
	db *sql.DB
}

func (t *Users) Exists(id string) (bool, error) {
	row := t.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id)
	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (t *Users) Add(id string) error {
	_, err := t.db.Exec("INSERT INTO (users) VALUES ($1)", id)
	return err
}
