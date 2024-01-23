package database

import (
	"github.com/febzey/ForestBot-Mainframe/types"
)

func (d *Database) SavePlayerLeave(args types.MinecraftPlayerLeaveMessage) error {
	uuid := args.Uuid
	server := args.Server
	timestamp := args.Timestamp

	_, err := d.pool.Exec("UPDATE users set leaves = leaves + 1, lastseen = ? WHERE uuid = ? AND mc_server = ?", timestamp, uuid, server)
	if err != nil {
		return err
	}

	return nil
}
