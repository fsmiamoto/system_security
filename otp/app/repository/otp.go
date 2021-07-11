package repository

import (
	"context"
	"database/sql"

	"github.com/fsmiamoto/system_security/otp/app/db"
)

func Exists(otp string) bool {
	stmt := `SELECT otp FROM invalid_otps WHERE otp = ?`
	row := db.DB.QueryRow(stmt, otp)
	return row.Scan() != sql.ErrNoRows
}

func Invalidate(otps []string) error {
	stmt := `INSERT INTO invalid_otps (otp) VALUES (?)`

	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for i := range otps {
		_, err := db.DB.Exec(stmt, otps[i])
		if err != nil {
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
