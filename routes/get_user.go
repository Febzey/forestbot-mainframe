package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
	"github.com/gorilla/mux"
)

// Getting a user by their name
func (c *Controller) GetUserByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	server := vars["server"]

	players, err := c.Database.GetUserByName(name, server)
	if err != nil {
		c.Logger.Error(err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, players)
}

// Getting a user by their UUID
func (c *Controller) GetUserByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	server := vars["server"]

	players, err := c.Database.GetUserByUUID(uuid, server)
	if err != nil {
		c.Logger.Error(err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, players)
}
