package controllers

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetUserWordOccurence(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	server := r.URL.Query().Get("server")
	word := r.URL.Query().Get("word")

	SELECT_WORD_COUNT_QUERY := "SELECT name, message, COUNT(message) AS cnt FROM messages WHERE name=? AND mc_server = ? AND message LIKE ? GROUP BY name"
	rows, err := c.Database.Query(SELECT_WORD_COUNT_QUERY, username, server, "%"+word+"%")
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	defer rows.Close()

	type WordCount struct {
		Name    string `json:"name"`
		Count   int    `json:"count"`
		Message string `json:"message"`
	}

	wordCount := WordCount{}

	for rows.Next() {
		err := rows.Scan(&wordCount.Name, &wordCount.Message, &wordCount.Count)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}

	}

	utils.RespondWithJSON(w, http.StatusOK, wordCount)

}
