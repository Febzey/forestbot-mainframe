package controllers

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

// TODO! fix this!!!!!!!!!

// METHOD: GET
// PATH: /quote
// QUERIES: name, server
// RESPOND: JSON
// DESCRIPTION: Get a random quote from a user on a server
// example http://localhost:5000/api/v1/quote?name=Febzey&server=newtest_new1
func (c *Controller) GetRandomQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	name := r.URL.Query().Get("name")
	server := r.URL.Query().Get("server")

	if name == "" || server == "" {
		http.Error(w, "Invalid 'name' and 'server' parameter required", http.StatusBadRequest)
		return
	}

	rows, err := c.Database.Query("SELECT name,message,date,mc_server,uuid FROM messages WHERE mc_server = ? AND name = ? AND LENGTH(message) > 10 ORDER BY RAND() LIMIT 1", server, name)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Internal Database Error - Contact Febzey or IncognitoMode", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var message types.MinecraftChatMessage

	for rows.Next() {

		err := rows.Scan(
			&message.Name,
			&message.Message,
			&message.Date,
			&message.Mc_server,
			&message.Uuid,
		)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, "Internal Database Error - Contact Febzey or IncognitoMode", http.StatusInternalServerError)
			return
		}

	}

	utils.RespondWithJSON(w, http.StatusOK, message)

}
