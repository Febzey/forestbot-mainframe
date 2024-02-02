package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

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
