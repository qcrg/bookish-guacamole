package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Type string `json:"typ"`
	jwt.RegisteredClaims
}

func registered_claims(id string, timeout time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(timeout)},
		ID:        id,
	}
}

func NewAccessClaims(id string, timeout time.Duration) Claims {
	return Claims{
		Type:             TypeAccess,
		RegisteredClaims: registered_claims(id, timeout),
	}
}

func NewRefreshClaims(id string, timeout time.Duration) Claims {
	return Claims{
		Type:             TypeRefresh,
		RegisteredClaims: registered_claims(id, timeout),
	}
}
