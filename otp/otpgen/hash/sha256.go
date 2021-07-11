package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

const maxLength = 32

func Sha256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])[:maxLength]
}
