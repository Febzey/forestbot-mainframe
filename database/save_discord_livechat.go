package database

import (
	"regexp"

	"github.com/febzey/ForestBot-Mainframe/types"
)

func (d *Database) SaveDiscordLiveChat(args types.LivechatChannel) error {

	regex := regexp.MustCompile(`/[^a-zA-Z0-9\s]/g`)
	cleanGuildName := regex.ReplaceAllString(args.GuildName, "")

	var query = `
	INSERT INTO livechats (guildName, guildID, channelID, setupBy, date, mc_server)
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	guildName = VALUES(guildName),
	setupBy = VALUES(setupBy),
	date = VALUES(date),
	mc_server = VALUES(mc_server);
	`

	_, err := d.Execute(
		query,
		cleanGuildName,
		args.GuildID,
		args.ChannelID,
		args.Setupby,
		args.Date,
		args.Mc_server,
	)
	if err != nil {
		return err
	}

	return nil
}
