package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetAllPlayerStatisticsByUsername(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Invalid 'username' parameter required.", http.StatusBadRequest)
		return
	}

	userstats, err := c.Database.GetAllPlayerStatisticsByUsername(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userstats)
}
