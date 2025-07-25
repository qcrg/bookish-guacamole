package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DB struct {
	db *sql.DB

	debug string
}

func (t *DB) CreateSession(
	user_id string,
	token_id string,
	refresh_token_hash string,
	user_agent_hash string,
	init_ip_hash string,
	expires_at time.Time,
) (string, error) {
	tr, err := t.db.Begin()
	if err != nil {
		return "", err
	}
	defer tr.Rollback()
	_, err = tr.Exec(
		"INSERT INTO tokens (id, refresh) VALUES ($1, $2)",
		token_id,
		refresh_token_hash,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert tokens: %w", err)
	}
	session_id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("Failed to create UUIDv7 for new session: %w", err)
	}
	_, err = tr.Exec(
		"INSERT INTO sessions (id, user_id, token_id, user_agent, init_ip, expires_at) VALUES ($1, $2, $3, $4, $5, $6)",
		session_id.String(),
		user_id,
		token_id,
		user_agent_hash,
		init_ip_hash,
		expires_at,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert new session: %w", err)
	}

	err = tr.Commit()
	if err != nil {
		return "", fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return session_id.String(), nil
}

func (t *DB) UpdateSession(
	session_id string,
	old_token_id string,
	new_token_id string,
	refresh_token_hash string,
	expires_at time.Time,
) error {
	tr, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tr.Rollback()

	_, err = tr.Exec(
		"INSERT INTO tokens (id, refresh) VALUES ($1, $2)",
		new_token_id,
		refresh_token_hash,
	)
	if err != nil {
		return fmt.Errorf("Failed to create new token: %d", err)
	}
	_, err = tr.Exec(
		"UPDATE sessions SET token_id = $1, expires_at = $2 WHERE id = $3",
		new_token_id,
		expires_at,
		session_id,
	)
	if err != nil {
		return fmt.Errorf("Failed to update session token: %d", err)
	}
	_, err = tr.Exec("DELETE FROM tokens WHERE id = $1", old_token_id)
	if err != nil {
		return fmt.Errorf("Failed to delete old token: %d", err)
	}

	err = tr.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return nil
}

func (t *DB) PurgeSessionFromTokenId(token_id string) error {
	tr, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tr.Rollback()

	session_id, err := t.Sessions().GetIdFromTokenId(token_id)
	if err != nil {
		return fmt.Errorf("Failed to get session_id: %w", err)
	}

	_, err = tr.Exec("DELETE FROM sessions WHERE id = $1", session_id)
	if err != nil {
		return fmt.Errorf("Failed to delete session: %w", err)
	}
	_, err = tr.Exec("DELETE FROM tokens WHERE id = $1", token_id)
	if err != nil {
		return fmt.Errorf("Failed to delete token: %w", err)
	}

	err = tr.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return nil
}

func (t *DB) Users() Users {
	return Users{t.db}
}

func (t *DB) Sessions() Sessions {
	return Sessions{t.db}
}

func (t *DB) Tokens() Tokens {
	return Tokens{t.db}
}

func (t *DB) Close() error {
	return t.db.Close()
}

func (t *DB) GetDebugStr() string {
	return t.debug
}

func NewDatabase(conf Config) (*DB, error) {
	// FIXME
	if "disable" != conf.GetTLSMod() {
		log.Fatal().Msgf("Only 'disable' TLS mod is implemented")
	}

	db, err := sql.Open("postgres", conf.GetConnectionString())
	if err != nil {
		return nil, err
	}

	return &DB{db: db, debug: conf.GetConnectionString()}, nil
}
