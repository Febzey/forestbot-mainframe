package routes

import (
	"encoding/json"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

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
