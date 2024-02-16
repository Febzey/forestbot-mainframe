package routes

import (
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

// METHOD: GET
// PATH: /kills
// QUERIES: uuid, server, limit, order
// RESPONSE: JSON
// DESCRIPTION: Gets the kills of a player
// example: http://localhost:5000/api/v1/kills?uuid=1&server=2&limit=3&order=DESC
func (c *Controller) GetMinecraftKills(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	server := r.URL.Query().Get("server")
	limit := r.URL.Query().Get("limit")
	order := r.URL.Query().Get("order")

	//if any of these are empty, return a bad request
	if uuid == "" || server == "" {
		http.Error(w, "Invalid 'uuid', 'server' parameter required. limit & order are optional.", http.StatusBadRequest)
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
	numPages := (limitInt + pageSize - 1) / pageSize

	kills := []types.MinecraftPlayerDeathMessage{}

	for page := 1; page <= numPages; page++ {

		offset := (page - 1) * pageSize

		rows, err := c.Database.Query("SELECT * FROM deaths where mc_server = ? AND murdererUUID = ? ORDER BY time "+order+" LIMIT ? OFFSET ?", server, uuid, pageSize, offset)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}

		for rows.Next() {

			if (len(kills) + 1) > limitInt {
				break
			}

			var killMessage types.MinecraftPlayerDeathMessage

			err := rows.Scan(
				&killMessage.Victim,
				&killMessage.Death_message,
				&killMessage.Murderer,
				&killMessage.Time,
				&killMessage.Type,
				&killMessage.Mc_server,
				&killMessage.Id,
				&killMessage.VictimUUID,
				&killMessage.MurdererUUID,
			)

			if err != nil {
				http.Error(w, "Internal Database Error - Please contact Febzey", http.StatusInternalServerError)
				c.Logger.Error(err.Error())
				return
			}

			kills = append(kills, killMessage)

		}

		rows.Close()

	}

	utils.RespondWithJSON(w, http.StatusOK, kills)

}
