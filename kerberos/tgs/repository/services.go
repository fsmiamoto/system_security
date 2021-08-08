package repository

import (
	"github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/fsmiamoto/system_security/kerberos/tgs/db"
)

func Get(id string) (*contracts.Service, error) {
	const stmt = `SELECT id, secret_key, iv FROM services WHERE id = ?`

	c := &contracts.Service{}

	row := db.DB.QueryRow(stmt, id)

	if err := row.Scan(&c.ID, &c.SecretKey, &c.InitVector); err != nil {
		return nil, err
	}

	return c, nil
}
