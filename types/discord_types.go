package types

import "database/sql"

type Guild struct {
	Guild_id   string `json:"guild_id"`
	Channel_id sql.NullString
	Mc_server  string `json:"mc_server"`
	Setup_by   string `json:"setup_by"`
	Created_at int    `json:"created_at"`
	Guild_name string `json:"guild_name"`
}

type LivechatChannel struct {
	GuildName string `json:"guildName"`
	GuildID   string `json:"guildID"`
	ChannelID string `json:"channelID"`
	Setupby   string `json:"setupBy"`
	Date      string `json:"date"`
	Mc_server string `json:"mcServer"`
}
