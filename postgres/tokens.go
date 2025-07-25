package postgres

import "database/sql"

type Tokens struct {
	db *sql.DB
}

func (t Tokens) Exists(id string) (bool, error) {
	row := t.db.QueryRow("SELECT COUNT(*) FROM tokens WHERE id = $1", id)
	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (t Tokens) GetRefreshHash(id string) (string, error) {
	row := t.db.QueryRow("SELECT refresh FROM tokens WHERE id = $1", id)
	var hash string
	err := row.Scan(&hash)
	return hash, err
}

func (t Tokens) Delete(id string) error {
	_, err := t.db.Exec("DELETE FROM tokens WHERE id = $1", id)
	return err
}
