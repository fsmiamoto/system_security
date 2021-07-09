package hash

import "crypto/sha256"

func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}
