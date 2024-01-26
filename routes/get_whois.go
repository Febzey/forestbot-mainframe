package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

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
