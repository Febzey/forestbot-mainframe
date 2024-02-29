package database

import (
	"regexp"

	"github.com/febzey/ForestBot-Mainframe/types"
)

func (d *Database) SaveDiscordGuild(guild types.Guild) error {
	// Clean guild_name using regex
	regex := regexp.MustCompile(`/[^a-zA-Z0-9\s]/g`)
	cleanGuildName := regex.ReplaceAllString(guild.Guild_name, "")

	var query = `
	INSERT INTO guilds (guild_id, channel_id, mc_server, setup_by, created_at, guild_name)
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE channel_id = ?, mc_server = ?, setup_by = ?, created_at = ?
	`

	_, err := d.Execute(
		query,
		guild.Guild_id,
		guild.Channel_id,
		guild.Mc_server,
		guild.Setup_by,
		guild.Created_at,
		cleanGuildName,
		guild.Channel_id,
		guild.Mc_server,
		guild.Setup_by,
		guild.Created_at,
	)
	if err != nil {
		return err
	}

	return nil
}
