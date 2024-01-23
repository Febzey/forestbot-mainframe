package database

import "github.com/febzey/ForestBot-Mainframe/types"

func (d *Database) SaveMinecraftChatMessage(message types.MinecraftChatMessage) error {
	_, err := d.pool.Exec("INSERT INTO messages (name, message, date, mc_server, uuid) VALUES (?, ?, ?, ?, ?)", message.Name, message.Message, message.Date, message.Mc_server, message.Uuid)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetMinecraftChatMessages() error {
	//implement this

	return nil
}
