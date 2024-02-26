package controllers

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetAvailableServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	servers, err := c.Database.UniqueServers()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, servers)
}
