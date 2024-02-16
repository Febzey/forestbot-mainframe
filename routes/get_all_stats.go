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

	uuid, err := c.Database.ConvertUsernameToUUID(username)
	if err != nil {
		http.Error(w, "Internal database error - contact febzey", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	if uuid == nil || !uuid.UUID.Valid || uuid.UUID.String == "" {
		http.Error(w, "Cant find the uuid.", http.StatusInternalServerError)
		return
	}

	//TODO: change this to uuid based.
	userstats, err := c.Database.GetAllPlayerStatisticsByUUID(uuid.UUID.String)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userstats)
}
