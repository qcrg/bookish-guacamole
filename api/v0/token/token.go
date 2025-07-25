package token

import "crypto/ed25519"

var skey ed25519.PrivateKey
var pkey ed25519.PublicKey

const (
	TypeAccess   = "access"
	TypeRefresh = "refresh"
)
