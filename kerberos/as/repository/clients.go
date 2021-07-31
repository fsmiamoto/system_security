package repository

import (
	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/as/db"
)

func Get(id string) (*contracts.Client, error) {
	const stmt = `SELECT id, secret_key, iv FROM clients WHERE id = ?`

	c := &contracts.Client{}

	row := db.DB.QueryRow(stmt, id)

	if err := row.Scan(&c.ID, &c.SecretKey, &c.InitVector); err != nil {
		return nil, err
	}

	return c, nil
}
