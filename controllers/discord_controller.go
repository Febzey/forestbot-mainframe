package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

/**
Method: POST
Posting a new Discord Guild to the database
**/
func (c *Controller) PostDiscordGuild(w http.ResponseWriter, r *http.Request) {

	c.Logger.Info("PostDiscordGuild called")

	var Guild types.Guild

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&Guild)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	if err := c.Database.SaveDiscordGuild(Guild); err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error saving guild", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Guild)

	//we want to get: channel_id, created_at, guild_id, guild_name, mc_server, setup_by, setup_by_name

}

/**
METHOD: GET
Getting all discord guiilds stored in database.
**/
func (c *Controller) GetDiscordGuilds(w http.ResponseWriter, r *http.Request) {

	rows, err := c.Database.Query("SELECT * FROM guilds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var guilds []types.Guild
	for rows.Next() {

		var guild types.Guild
		err = rows.Scan(
			&guild.Guild_id,
			&guild.Channel_id,
			&guild.Mc_server,
			&guild.Setup_by,
			&guild.Created_at,
			&guild.Guild_name,
		)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, "Error with database", http.StatusInternalServerError)
			return
		}
		guilds = append(guilds, guild)
	}

	utils.RespondWithJSON(w, http.StatusOK, guilds)

}

/**
Method: Delete
Delete a discord guild, (when the bot leaves a discord server for example)
**/
func (c *Controller) DeleteDiscordGuild(w http.ResponseWriter, r *http.Request) {
	c.Logger.Info("DeleteDiscordGuild called")

	guild_id := r.URL.Query().Get("guild_id")

	if guild_id == "" {
		http.Error(w, "Invalid 'guild_id' parameter required", http.StatusBadRequest)
		return
	}

	err := c.Database.DeleteDiscordGuild(guild_id)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error deleting guild", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, guild_id)
}

/**
Method: GET
Getting all of our active LiveChats from our database
**/
func (c *Controller) GetDiscordLiveChatChannels(w http.ResponseWriter, r *http.Request) {

	rows, err := c.Database.Query("SELECT * FROM livechats")
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Internal Database Error.", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var livechats []types.LivechatChannel
	for rows.Next() {

		var livechat types.LivechatChannel
		err = rows.Scan(
			&livechat.GuildName,
			&livechat.GuildID,
			&livechat.ChannelID,
			&livechat.Setupby,
			&livechat.Date,
			&livechat.Mc_server,
		)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, "Internal Database Error.", http.StatusInternalServerError)
			return
		}
		livechats = append(livechats, livechat)
	}

	utils.RespondWithJSON(w, http.StatusOK, livechats)

}

/**
Method: Delete
Deleting a discord live chat
**/
func (c *Controller) DeleteDiscordLiveChat(w http.ResponseWriter, r *http.Request) {
	c.Logger.Info("DeleteDiscordLiveChat called")

	channel_id := r.URL.Query().Get("channel_id")

	if channel_id == "" {
		http.Error(w, "Invalid 'channel_id' parameter required", http.StatusBadRequest)
		return
	}

	err := c.Database.DeleteDiscordLiveChat(channel_id)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error deleting livechat", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, channel_id)
}

/**
Method: POST
post a discord livechat for forestbot discord client
**/
func (c *Controller) PostDiscordLiveChat(w http.ResponseWriter, r *http.Request) {
	c.Logger.Info("PostDiscordLiveChat called")

	var Livechat types.LivechatChannel

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&Livechat)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	if err := c.Database.SaveDiscordLiveChat(Livechat); err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error saving livechat", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Livechat)
}
