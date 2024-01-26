package routes

import (
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

// METHOD: GET
// PATH: /deaths
// QUERIES: uuid, server, limit, order
// RESPONSE: JSON
// DESCRIPTION: Gets the deaths of a player
// example: http://localhost:5000/api/v1/deaths?uuid=1&server=2&limit=3&order=DESC
func (c *Controller) GetMinecraftDeaths(w http.ResponseWriter, r *http.Request) {
	// server := r.URL.Query().Get("server")
	uuid := r.URL.Query().Get("uuid")
	limit := r.URL.Query().Get("limit")
	order := r.URL.Query().Get("order")
	server := r.URL.Query().Get("server")

	//if any of these are empty, return a bad request
	if uuid == "" || server == "" || limit == "" || order == "" {
		http.Error(w, "Invalid 'uuid', 'server', 'limit', and 'order' parameter required", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
		return
	}

	pageSize := 100

	numPages := (limitInt + pageSize - 1) / pageSize

	deaths := []types.MinecraftPlayerDeathMessage{}

	for page := 1; page <= numPages; page++ {
		offset := (page - 1) * pageSize

		rows, err := c.Database.Query("SELECT * FROM deaths where mc_server = ? AND victimUUID = ? ORDER BY time "+order+" LIMIT ? OFFSET ?", server, uuid, pageSize, offset)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}

		for rows.Next() {

			if (len(deaths) + 1) > limitInt {
				break
			}

			var deathMessage types.MinecraftPlayerDeathMessage

			err := rows.Scan(
				&deathMessage.Victim,
				&deathMessage.Death_message,
				&deathMessage.Murderer,
				&deathMessage.Time,
				&deathMessage.Type,
				&deathMessage.Mc_server,
				&deathMessage.Id,
				&deathMessage.VictimUUID,
				&deathMessage.MurdererUUID,
			)
			if err != nil {
				http.Error(w, "Internal Database Error - Please contact Febzey", http.StatusInternalServerError)
				c.Logger.Error(err.Error())
				return
			}

			deaths = append(deaths, deathMessage)

		}

	}

	utils.RespondWithJSON(w, http.StatusOK, deaths)

}
