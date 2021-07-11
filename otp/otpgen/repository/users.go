package repository

import (
	"math/rand"
	"strconv"

	"github.com/fsmiamoto/system_security/otp/otpgen/db"
	"github.com/fsmiamoto/system_security/otp/otpgen/hash"
)

type User struct {
	Username string
	Password string
	Seed     string
	Salt     string
}

func Get(username string) (*User, error) {
	stmt := `SELECT username, password, seed, salt FROM users WHERE username = ?`

	row := db.DB.QueryRow(stmt, username)

	u := &User{}
	if err := row.Scan(&u.Username, &u.Password, &u.Seed, &u.Salt); err != nil {
		return nil, err
	}

	return u, nil
}

func Add(username, password string) error {
	stmt := `INSERT INTO users (username, password, seed, salt) VALUES (?,?,?,?)`

	salt := hash.Sha256(strconv.Itoa(rand.Int()))
	seed := hash.Sha256(strconv.Itoa(rand.Int()))

	hashedPassword := hash.Sha256(password + salt)

	_, err := db.DB.Exec(stmt, username, hashedPassword, seed, salt)
	return err
}

func Remove(username string) error {
	stmt := `DELETE FROM users WHERE username = ?`
	_, err := db.DB.Exec(stmt, username)
	return err
}
