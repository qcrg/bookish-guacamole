package token

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	ErrStrTokenExpired     = "Token expired"
	ErrStrInvalidSignature = "Invalid signature"
)

func Parse(token string) (token_id string, typ string, err error) {
	var claims Claims
	_, err = jwt.ParseWithClaims(
		token,
		&claims,
		func(*jwt.Token) (any, error) {
			return pkey, nil
		},
	)
	if err != nil {
		return
	}
	token_id = claims.ID
	typ = claims.Type
	return
}
