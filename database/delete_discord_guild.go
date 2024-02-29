package database

func (d *Database) DeleteDiscordGuild(guildID string) error {
	_, err := d.Execute("DELETE FROM guilds WHERE guild_id = ?", guildID)
	if err != nil {
		return err
	}
	return nil
}
