package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	DefaultAccessTokenTimeout  time.Duration = 15 * time.Minute
	DefaultRefreshTokenTimeout time.Duration = 30 * 24 * time.Hour
)

func GenPair(
	access_token_timeout time.Duration,
	refresh_token_timeout time.Duration,
) (access string, refresh string, id string, err error) {
	uuidv7, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("Failed to generate UUIDv7: %w", err)
		return
	}
	id = uuidv7.String()

	access_claims := NewAccessClaims(id, access_token_timeout)
	refresh_claims := NewRefreshClaims(id, refresh_token_timeout)

	atoken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, access_claims)
	if atoken == nil {
		err = fmt.Errorf("Failed to create access token")
		return
	}
	rtoken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, refresh_claims)
	if rtoken == nil {
		err = fmt.Errorf("Failed to create refresh token")
		return
	}
	access, err = atoken.SignedString(skey)
	if err != nil {
		err = fmt.Errorf("Failed to sign access token: %w", err)
		return
	}
	refresh, err = rtoken.SignedString(skey)
	if err != nil {
		err = fmt.Errorf("Failed to sign refresh token: %w", err)
		return
	}
	return
}
