package routes

import (
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

// METHOD: GET
// PATH: /advancements
// QUERIES: uuid, server, limit, order
// RESPONSE: JSON
// DESCRIPTION: Gets the advancements of a player
func (c *Controller) getAdvancements(w http.ResponseWriter, r *http.Request) {
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

	// Default page size
	pageSize := 100

	// Parse limit as an integer
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
		return
	}

	// Calculate the number of pages based on the limit and page size
	numPages := (limitInt + pageSize - 1) / pageSize

	advancements := []types.MinecraftAdvancementMessage{}

	// Loop through pages
	for page := 1; page <= numPages; page++ {
		// Calculate the offset for the current page
		offset := (page - 1) * pageSize

		// Query database for the current page
		rows, err := c.Database.Query("SELECT * FROM advancements WHERE mc_server = ? AND uuid = ? ORDER BY time "+order+" LIMIT ? OFFSET ?", server, uuid, pageSize, offset)

		if err != nil {
			// Log the error and send a 500 to the client
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}

		// Process the current page of rows
		for rows.Next() {
			var advancement types.MinecraftAdvancementMessage
			err := rows.Scan(
				&advancement.Username,
				&advancement.Advancement,
				&advancement.Time,
				&advancement.Mc_server,
				&advancement.Id,
				&advancement.Uuid,
			)
			if err != nil {
				// Log the error and send a 500 to the client
				http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
				c.Logger.Error(err.Error())
				return
			}
			advancements = append(advancements, advancement)

		}

		// Close the rows for the current page
		rows.Close()

		// Process or send the current page of advancements as needed
		// (e.g., send it to the client or perform further processing)
		// ...

		// If needed, you can break out of the loop early based on certain conditions
		// (e.g., if you reach a maximum number of processed rows)
	}

	utils.RespondWithJSON(w, http.StatusOK, advancements)
}
