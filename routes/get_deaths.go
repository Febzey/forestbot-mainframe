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

	uuid := r.URL.Query().Get("uuid")
	server := r.URL.Query().Get("server")
	limit := r.URL.Query().Get("limit")
	order := r.URL.Query().Get("order")
	killType := r.URL.Query().Get("type")

	//if any of these are empty, return a bad request
	if uuid == "" || server == "" {
		http.Error(w, "Invalid 'uuid', 'server' parameter required. limit & order are optional.", http.StatusBadRequest)
		return
	}

	if killType == "" {
		killType = "all"
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

	numPages := (limitInt + pageSize - 1) / pageSize

	deaths := []types.MinecraftPlayerDeathMessage{}

	for page := 1; page <= numPages; page++ {
		offset := (page - 1) * pageSize

		SELECT_DEATHS_QUERY := ""

		if killType == "all" {
			SELECT_DEATHS_QUERY = "SELECT * FROM deaths where mc_server = ? AND victimUUID = ? ORDER BY time " + order + " LIMIT ? OFFSET ?"
		}

		if killType == "pvp" {
			SELECT_DEATHS_QUERY = "SELECT * FROM deaths where mc_server = ? AND victimUUID = ? AND type = 'pvp' ORDER BY time " + order + " LIMIT ? OFFSET ?"
		}

		if killType == "pve" {
			SELECT_DEATHS_QUERY = "SELECT * FROM deaths where mc_server = ? AND victimUUID = ? AND type = 'pve' ORDER BY time " + order + " LIMIT ? OFFSET ?"
		}

		rows, err := c.Database.Query(SELECT_DEATHS_QUERY, server, uuid, pageSize, offset)
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
