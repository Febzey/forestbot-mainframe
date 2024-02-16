package database

import (
	"database/sql"
	"errors"
)

type UUID struct {
	UUID sql.NullString `json:"uuid"`
}

func (d *Database) ConvertUsernameToUUID(username string) (*UUID, error) {
	if username == "" {
		return nil, errors.New("invalid 'username' parameter required")
	}

	var GET_UUID_FROM_USERNAME_QUERY = "SELECT DISTINCT uuid FROM users WHERE username = ?"
	rows, err := d.Query(GET_UUID_FROM_USERNAME_QUERY, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var uuid UUID

	for rows.Next() {
		err := rows.Scan(&uuid.UUID)
		if err != nil {
			return nil, err
		}
	}

	return &uuid, nil
}
