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

	defer rows.Close()

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

	defer rows.Close() // Defer the closing of rows

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

// Getting all stats for a user from database on all servers theyve been seen on by username
func (d *Database) GetAllPlayerStatisticsByUsername(username string) ([]types.User, error) {
	var users []types.User

	rows, err := d.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var user types.User
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
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (d *Database) GetAllPlayerStatisticsByUUID(uuid string) ([]types.User, error) {
	var users []types.User

	rows, err := d.Query("SELECT * FROM users WHERE uuid = ?", uuid)
	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var user types.User
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
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}
