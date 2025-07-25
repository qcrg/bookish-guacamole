package token

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

const (
	_TOKEN_KEY_PREFIX = "BHGL_TOKEN_"
	SEED_KEY          = _TOKEN_KEY_PREFIX + "SEED"
)

func Init() error {
	seed_str, exists := os.LookupEnv(SEED_KEY)
	if !exists {
		return errors.New(SEED_KEY + " is empty")
	}
	seed, err := base64.StdEncoding.DecodeString(seed_str)
	if err != nil {
		return fmt.Errorf("Couldn't decode the encoded seed in base64: %w", err)
	}
	if len(seed) != ed25519.SeedSize {
		return fmt.Errorf(
			"Decoded seed size(%d bytes) is not equal to %d bytes",
			len(seed),
			ed25519.SeedSize,
		)
	}
	skey = ed25519.NewKeyFromSeed(seed)
	pkey = skey.Public().(ed25519.PublicKey)
	return nil
}
