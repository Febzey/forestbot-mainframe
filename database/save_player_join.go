package database

import (
	"time"

	"github.com/febzey/ForestBot-Mainframe/types"
)

type Result struct {
	Action string
	Data   map[string]interface{}
}

func (d *Database) SavePlayerJoin(message types.MinecraftPlayerJoinMessage) (*Result, error) {

	user := message.Username
	server := message.Server
	uuid := message.Uuid
	timestamp := message.Timestamp

	no_action := &Result{
		Action: "none",
		Data:   map[string]interface{}{},
	}

	//Getting the user to see if they already exist in the database:
	rows, err := d.pool.Query("SELECT * FROM users WHERE uuid = ? AND mc_server = ?", uuid, server)
	if err != nil {
		return no_action, err
	}
	defer rows.Close()

	//check if the user does not exist
	if !rows.Next() {
		//if the user does not exist, create them
		_, err := d.pool.Exec("INSERT INTO users(username, joindate, uuid, joins, mc_server) VALUES (?,?,?,?,?)", user, timestamp, uuid, 1, server)
		if err != nil {
			return no_action, err
		}

		return &Result{
			Action: "new_user",
			Data: map[string]interface{}{
				"username": user,
			},
		}, nil

	} else {

		var userFromDatabase types.User
		err := rows.Scan(
			&userFromDatabase.Username,
			&userFromDatabase.Kills,
			&userFromDatabase.Deaths,
			&userFromDatabase.Joindate,
			&userFromDatabase.LastSeen,
			&userFromDatabase.UUID,
			&userFromDatabase.Playtime,
			&userFromDatabase.Joins,
			&userFromDatabase.Leaves,
			&userFromDatabase.LastDeathTime,
			&userFromDatabase.LastDeathString,
			&userFromDatabase.MCServer,
		)
		if err != nil {
			return no_action, err
		}

		//if the user does exist, update their join count
		_, err = d.pool.Exec("UPDATE users SET joins = joins + 1, lastseen = ? WHERE uuid = ? AND mc_server = ?", timestamp, uuid, server)
		if err != nil {
			return no_action, err
		}

		loginEventData := types.PlayerActivity{
			UUID:      uuid,
			Username:  user,
			Date:      time.Now().UnixNano() / int64(time.Millisecond),
			Type:      "login",
			Mc_server: server,
		}

		insertLoginActivity := "INSERT INTO playerActivity(uuid, username, date, type, mc_server) VALUES (?,?,?,?,?)"
		_, err = d.pool.Exec(insertLoginActivity, loginEventData.UUID, loginEventData.Username, loginEventData.Date, loginEventData.Type, loginEventData.Mc_server)
		if err != nil {
			return no_action, err
		}

		//if the username is different from the one in the database, update it
		if user != userFromDatabase.Username {
			_, err := d.pool.Exec("UPDATE users SET username = ? WHERE username = ? AND uuid = ? AND mc_server = ?", user, userFromDatabase.Username, uuid, server)
			if err != nil {
				return no_action, err
			}

			return &Result{
				Action: "new_name",
				Data: map[string]interface{}{
					"old_name": userFromDatabase.Username,
					"new_name": user,
				},
			}, nil

		}

	}

	return no_action, nil

}
