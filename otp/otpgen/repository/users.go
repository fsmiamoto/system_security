package repository

import (
	"math/rand"
	"strconv"

	"github.com/fsmiamoto/system_security/otp/otpgen/db"
	"github.com/fsmiamoto/system_security/otp/otpgen/hash"
)

const maxHashLength = 32

func Add(username, password string) error {
	stmt := `INSERT INTO users (username, password, seed, salt) VALUES (?,?,?,?)`

	salt := hash.Sha256(strconv.Itoa(rand.Int()))[:maxHashLength]
	seed := hash.Sha256(strconv.Itoa(rand.Int()))[:maxHashLength]

	hashedPassword := hash.Sha256(password + salt)[:maxHashLength]

	_, err := db.DB.Exec(stmt, username, hashedPassword, seed, salt)
	return err
}

func Remove(username string) error {
	stmt := `DELETE FROM users WHERE username = ?`
	_, err := db.DB.Exec(stmt, username)
	return err
}
