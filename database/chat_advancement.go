package database

import "github.com/febzey/ForestBot-Mainframe/types"

func (d *Database) SaveMinecraftAdvancementMessage(message types.MinecraftAdvancementMessage) error {
	_, err := d.pool.Exec("INSERT INTO advancements (username, advancement, time, mc_server, uuid) VALUES (?, ?, ?, ?, ?)",
		message.Username, message.Advancement, message.Time, message.Mc_server, message.Uuid)
	if err != nil {
		return err
	}

	return nil
}
