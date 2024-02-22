package routes

import (
	"encoding/json"
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

/*
**

	METHOD: GET
	DESCRIPTION: Gets all statistics for a username on all servers they've been seen on.

**
*/
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

/*
*

	METHOD: GET
	Queires: username
	Description: Get whois description for a specific username.

*
*/
func (c *Controller) GetWhoIs(w http.ResponseWriter, r *http.Request) {

	//get the query for username
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Invalid 'username' parameter required", http.StatusBadRequest)
		return
	}

	//SELECT username, description FROM whois WHERE username=? is our query

	rows, err := c.Database.Query("SELECT username, description FROM whois WHERE username=?", username)
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	defer rows.Close()

	//we need to create a slice of strings to store the usernames in.
	var usernames []string
	var descriptions []string

	//we need to loop through the rows and append the usernames to the slice.
	for rows.Next() {
		var username string
		var description string
		err := rows.Scan(&username, &description)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}
		usernames = append(usernames, username)
		descriptions = append(descriptions, description)
	}

	//we need to check if the slice is empty.
	if len(usernames) == 0 {
		http.Error(w, "No usernames found", http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, descriptions)

}

/*
**

	METHOD: POST
	BODY: json
	RESPONSE: json
	DESCRIPTION: Sets the description for whois for specific username

**
*/
func (c *Controller) POSTIamDescription(w http.ResponseWriter, r *http.Request) {

	var IamBody struct {
		Username    string `json:"username"`
		Description string `json:"description"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&IamBody)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	err = c.Database.INSERT_player_whois_description(IamBody.Username, IamBody.Description)
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, IamBody)
}
