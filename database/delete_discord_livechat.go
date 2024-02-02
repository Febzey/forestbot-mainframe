package database

func (d *Database) DeleteDiscordLiveChat(channelID string) error {
	_, err := d.pool.Exec("DELETE FROM livechats WHERE channelID = ?", channelID)
	if err != nil {
		return err
	}
	return nil
}
