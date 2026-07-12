package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	refresh_token := make([]byte, 32)
	if _, err := rand.Read(refresh_token); err != nil {
		panic(err)
	}

	return hex.EncodeToString(refresh_token)
}
