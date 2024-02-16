package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetNameSearch(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	mcServer := r.URL.Query().Get("server")

	if username == "" || mcServer == "" {
		http.Error(w, "Invalid 'username' and 'server' parameter required", http.StatusBadRequest)
		return
	}

	//we need to search for the user in our database.
	//what this command does is allow a user to find a player by their username
	//will only a few letters of their name. if they forget the rest of their name.
	//this is useful for finding players that have a lot of characters in their name.
	//we will return about 6 results.

	rows, err := c.Database.Query("SELECT username FROM users WHERE username LIKE ? AND mc_server = ? ORDER BY ABS(CHAR_LENGTH(username) - CHAR_LENGTH(?)), lastseen DESC LIMIT 6", "%"+username+"%", mcServer, username)
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	defer rows.Close()

	//we need to create a slice of strings to store the usernames in.
	var usernames []string

	//we need to loop through the rows and append the usernames to the slice.
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}
		usernames = append(usernames, username)
	}

	//we need to check if the slice is empty.
	if len(usernames) == 0 {
		http.Error(w, "No usernames found", http.StatusNotFound)
		return
	}

	//we need to respond with the json.
	utils.RespondWithJSON(w, http.StatusOK, usernames)

}
