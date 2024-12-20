package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

//todo implement all history requests heres.

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

	fmt.Println(limit, " limit")

	// need to get username instead fuck!, or do something like that incase we cant get the uuid.

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
	numPages := (limitInt + limitInt - 1) / limitInt

	advancements := []types.MinecraftAdvancementMessage{}

	// Loop through pages
	for page := 1; page <= numPages; page++ {
		// Calculate the offset for the current page
		offset := (page - 1) * limitInt

		fmt.Println(pageSize, " pageSize")
		fmt.Println(offset, " offset")

		// Query database for the current page
		rows, err := c.Database.Query("SELECT * FROM advancements WHERE mc_server = ? AND uuid = ? ORDER BY time "+order+" LIMIT ? OFFSET ?", server, uuid, limitInt, offset)

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
		rows.Close()

	}

	utils.RespondWithJSON(w, http.StatusOK, advancements)
}

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

	if limitInt > 100 {
		pageSize := 100
		offset := 0
		messages := []types.MinecraftChatMessage{}

		for offset < limitInt {
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
			offset += pageSize
		}

		utils.RespondWithJSON(w, http.StatusOK, messages)

	} else {

		rows, err := c.Database.Query("SELECT name,message,date,mc_server,uuid FROM messages WHERE mc_server = ? AND name = ? ORDER BY date "+order+" LIMIT ?", server, name, limit)
		if err != nil {
			http.Error(w, "Internal Database Error", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}
		defer rows.Close()

		messages := []types.MinecraftChatMessage{}

		for rows.Next() {
			var message types.MinecraftChatMessage

			err := rows.Scan(
				&message.Name,
				&message.Message,
				&message.Date,
				&message.Mc_server,
				&message.Uuid,
			)
			if err != nil {
				http.Error(w, "Internal Database Error", http.StatusInternalServerError)
				c.Logger.Error(err.Error())
				return
			}

			messages = append(messages, message)
		}

		utils.RespondWithJSON(w, http.StatusOK, messages)
	}
}

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

		// Close the rows for the current page
		rows.Close()

	}

	utils.RespondWithJSON(w, http.StatusOK, deaths)

}
