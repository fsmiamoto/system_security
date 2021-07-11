package otp

import (
	"time"

	"github.com/fsmiamoto/system_security/otp/app/hash"
)

const length = 6

func NewList(n int, seed, salt string) []string {
	now := time.Now().Format("2006-01-02 15:04")
	return newList(n, seed, salt, now)
}

func newList(n int, seed, salt, timestamp string) []string {
	if n == 0 {
		return nil
	}

	otps := make([]string, n)
	otps[0] = hash.Sha256(seed + salt + timestamp)[:length]
	for i := 1; i < len(otps); i++ {
		otps[i] = hash.Sha256(otps[i-1])[:length]
	}

	return otps
}
