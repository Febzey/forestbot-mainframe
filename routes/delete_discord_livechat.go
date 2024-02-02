package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

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
