package database

import (
	"github.com/febzey/ForestBot-Mainframe/types"
)

// GetUserByUUID returns a user from the database by their UUID
func (d *Database) GetUserByUUID(uuid string, server string) (types.User, error) {
	var user types.User

	rows, err := d.Query("SELECT * FROM users WHERE uuid = ? AND mc_server = ?", uuid, server)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		err := rows.Scan(
			&user.Username,
			&user.Kills,
			&user.Deaths,
			&user.Joindate,
			&user.LastSeen,
			&user.UUID,
			&user.Playtime,
			&user.Joins,
			&user.Leaves,
			&user.LastDeathTime,
			&user.LastDeathString,
			&user.MCServer,
		)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (d *Database) GetUserByName(username string, server string) (types.User, error) {
	var user types.User

	rows, err := d.Query("SELECT * FROM users WHERE username = ? AND mc_server = ?", username, server)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		err := rows.Scan(
			&user.Username,
			&user.Kills,
			&user.Deaths,
			&user.Joindate,
			&user.LastSeen,
			&user.UUID,
			&user.Playtime,
			&user.Joins,
			&user.Leaves,
			&user.LastDeathTime,
			&user.LastDeathString,
			&user.MCServer,
		)
		if err != nil {
			return user, err
		}

	}

	return user, nil
}
