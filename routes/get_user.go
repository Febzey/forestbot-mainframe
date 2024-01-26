package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

// Getting a user by their name
// example: http://localhost:5000/api/v1/user?name=1&server=2
// query: name, server
// description: Gets a user by their name
func (c *Controller) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	server := r.URL.Query().Get("server")

	players, err := c.Database.GetUserByName(name, server)
	if err != nil {
		c.Logger.Error(err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, players)
}

// Getting a user by their UUID
// example: http://localhost:5000/api/v1/user?uuid=1&server=2
// query: uuid, server
// description: Gets a user by their UUID
func (c *Controller) GetUserByUUID(w http.ResponseWriter, r *http.Request) {
	// uuid := vars["uuid"]
	// server := vars["server"]
	uuid := r.URL.Query().Get("uuid")
	server := r.URL.Query().Get("server")

	players, err := c.Database.GetUserByUUID(uuid, server)
	if err != nil {
		c.Logger.Error(err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, players)
}
