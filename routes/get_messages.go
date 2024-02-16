package routes

import (
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

//https://localhost:5000/messages?name=Febzey&server=newtest_new1&limit=100&order=DESC

// Getting a user by their name
// PATH: /messages
// QUERIES: name, server, limit, order
// RESPONSE: JSON
// DESCRIPTION: Gets the messages of a player
func (c *Controller) GetMessages(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	server := r.URL.Query().Get("server")
	limit := r.URL.Query().Get("limit")
	order := r.URL.Query().Get("order")

	//if any of these are empty, return a bad request
	if name == "" || server == "" {
		http.Error(w, "Invalid 'name', 'server' parameter required. limit & order are optional.", http.StatusBadRequest)
		return
	}

	if order == "" {
		order = "DESC"
	}

	if limit == "" {
		limit = "40"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
		return
	}

	pageSize := 100
	if limitInt < pageSize {
		pageSize = limitInt
	}

	numPages := (limitInt + pageSize - 1) / pageSize

	messages := []types.MinecraftChatMessage{}

	for page := 1; page <= numPages; page++ {
		offset := (page - 1) * pageSize

		rows, err := c.Database.Query("SELECT name,message,date,mc_server,uuid FROM messages WHERE mc_server = ? AND name = ? ORDER BY date "+order+" LIMIT ? OFFSET ?", server, name, pageSize, offset)

		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey or IncognitoMode on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}

		for rows.Next() {

			if (len(messages) + 1) > limitInt {
				break
			}

			var message types.MinecraftChatMessage

			err := rows.Scan(
				&message.Name,
				&message.Message,
				&message.Date,
				&message.Mc_server,
				&message.Uuid,
			)
			if err != nil {
				http.Error(w, "Internal Database Error - Please contact Febzey or IncognitoMode on Discord", http.StatusInternalServerError)
				c.Logger.Error(err.Error())
				return
			}

			messages = append(messages, message)
		}

		rows.Close()
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)

}
