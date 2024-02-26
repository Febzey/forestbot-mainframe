package controllers

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

// METHOD: GET
// PATH: /messagecount
// QUERIES: username, server

func (c *Controller) GetUserMessageCount(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	server := r.URL.Query().Get("server")

	SELECT_MESSAGE_COUNT_QUERY := "SELECT name,COUNT(name) AS cnt FROM messages WHERE name=? AND mc_server = ? HAVING cnt > 1"
	rows, err := c.Database.Query(SELECT_MESSAGE_COUNT_QUERY, username, server)
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	defer rows.Close()

	type MessageCount struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	messageCount := MessageCount{}

	for rows.Next() {
		err := rows.Scan(&messageCount.Name, &messageCount.Count)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, messageCount)

}
