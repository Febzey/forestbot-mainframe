package routes

import (
	"encoding/json"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

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
