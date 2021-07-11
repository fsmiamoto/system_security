package repository

import (
	"github.com/fsmiamoto/system_security/otp/app/db"
)

type User struct {
	Username string
	Seed     string
	Salt     string
}

func Get(username string) (*User, error) {
	stmt := `SELECT username, seed, salt FROM users WHERE username = ?`

	row := db.DB.QueryRow(stmt, username)

	u := &User{}
	if err := row.Scan(&u.Username, &u.Seed, &u.Salt); err != nil {
		return nil, err
	}

	return u, nil
}

func Add(username, seed, salt string) error {
	stmt := `INSERT INTO users (username, seed, salt) VALUES (?,?,?)`

	_, err := db.DB.Exec(stmt, username, seed, salt)
	return err
}

func Remove(username string) error {
	stmt := `DELETE FROM users WHERE username = ?`
	_, err := db.DB.Exec(stmt, username)
	return err
}
