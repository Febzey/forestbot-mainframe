package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetUserOnlineCheck(w http.ResponseWriter, r *http.Request) {

	//get the query for username
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Invalid 'username' parameter required", http.StatusBadRequest)
		return
	}

	//look for that username inside of our plaerList map
	for _, playerList := range c.PlayerLists {
		for _, player := range playerList {
			if player.Username == username {
				utils.RespondWithJSON(w, http.StatusOK, map[string]string{"online": "true", "server": player.Server})
				return
			}
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"online": "false"})
}
