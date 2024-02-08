package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetPlayerActivityData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mc_server := r.URL.Query().Get("server")
	username := r.URL.Query().Get("username")
	uuid := r.URL.Query().Get("uuid")

	usingUuid := false

	if username == "" && uuid == "" {
		http.Error(w, "Invalid 'username' or 'uuid' parameter required.", http.StatusBadRequest)
		return
	}

	if mc_server == "" {
		http.Error(w, "Invalid 'server' parameter required.", http.StatusBadRequest)
		return
	}

	if uuid != "" {
		usingUuid = true
		username = uuid
	}

	playerActivityData, err := c.Database.GetAllPlayerActivity(mc_server, username, usingUuid)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, playerActivityData)

}
