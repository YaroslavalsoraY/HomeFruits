package jwt

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	refreshToken := make([]byte, 32)
	rand.Read(refreshToken)
	return hex.EncodeToString(refreshToken)
}