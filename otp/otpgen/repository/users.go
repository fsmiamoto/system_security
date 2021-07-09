package repository

import (
	"math/rand"
	"strconv"

	"github.com/fsmiamoto/system_security/otp/otpgen/db"
	"github.com/fsmiamoto/system_security/otp/otpgen/hash"
)

func Add(username, password string) error {
	stmt := `INSERT INTO users (username, password, seed, salt) VALUES (?,?,?,?)`

	salt := rand.Int()
	seed := rand.Int()

	hashedSeed := hash.Sha256(strconv.Itoa(seed))
	hashedPassword := hash.Sha256(password)

	_, err := db.DB.Exec(stmt, username, hashedPassword, hashedSeed, salt)
	return err
}

func Remove(username string) error {
	stmt := `DELETE FROM users WHERE username = ?`
	_, err := db.DB.Exec(stmt, username)
	return err
}
