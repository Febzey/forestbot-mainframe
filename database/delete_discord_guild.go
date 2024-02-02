package database

func (d *Database) DeleteDiscordGuild(guildID string) error {
	_, err := d.pool.Exec("DELETE FROM guilds WHERE guild_id = ?", guildID)
	if err != nil {
		return err
	}
	return nil
}
