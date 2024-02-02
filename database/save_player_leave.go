package database

import (
	"time"

	"github.com/febzey/ForestBot-Mainframe/types"
)

func (d *Database) SavePlayerLeave(args types.MinecraftPlayerLeaveMessage) error {
	username := args.Username
	uuid := args.Uuid
	server := args.Server
	timestamp := args.Timestamp

	logoutEventActivity := &types.PlayerActivity{
		UUID:      uuid,
		Username:  username,
		Date:      time.Now().UnixNano() / int64(time.Millisecond),
		Type:      "logout",
		Mc_server: server,
	}

	_, err := d.pool.Exec("INSERT INTO playerActivity(uuid, username, date, type, mc_server) VALUES (?,?,?,?,?)", logoutEventActivity.UUID, logoutEventActivity.Username, logoutEventActivity.Date, logoutEventActivity.Type, logoutEventActivity.Mc_server)
	if err != nil {
		return err
	}

	_, err = d.pool.Exec("UPDATE users set leaves = leaves + 1, lastseen = ? WHERE uuid = ? AND mc_server = ?", timestamp, uuid, server)
	if err != nil {
		return err
	}

	return nil
}
