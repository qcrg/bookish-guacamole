package postgres

import "database/sql"

type Sessions struct {
	db *sql.DB
}

func (t Sessions) Exists(id string) (bool, error) {
	row := t.db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = $1", id)
	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (t Sessions) ExistsFromTokenId(token_id string) (bool, error) {
	row := t.db.QueryRow(
		"SELECT COUNT(*) FROM sessions WHERE token_id = $1",
		token_id,
	)
	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (t Sessions) GetUserId(id string) (string, error) {
	row := t.db.QueryRow("SELECT user_id FROM sessions WHERE id = $1", id)
	var user_id string
	err := row.Scan(&user_id)
	return user_id, err
}

func (t Sessions) GetUserIdFromTokenId(token_id string) (string, error) {
	row := t.db.QueryRow(
		"SELECT user_id FROM sessions WHERE token_id = $1",
		token_id,
	)
	var user_id string
	err := row.Scan(&user_id)
	return user_id, err
}

func (t Sessions) GetIdFromTokenId(token_id string) (string, error) {
	row := t.db.QueryRow(
		"SELECT id FROM sessions WHERE token_id = $1",
		token_id,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

func (t Sessions) GetUserAgent(id string) (string, error) {
	row := t.db.QueryRow("SELECT user_agent FROM sessions WHERE id = $1", id)
	var user_agent string
	err := row.Scan(&user_agent)
	return user_agent, err
}

func (t Sessions) GetInitIp(id string) (string, error) {
	row := t.db.QueryRow("SELECT init_ip FROM sessions WHERE id = $1", id)
	var init_ip string
	err := row.Scan(&init_ip)
	return init_ip, err
}

func (t Sessions) Delete(id string) error {
	_, err := t.db.Exec("DELETE FROM sessions WHERE id = $1", id)
	return err
}
