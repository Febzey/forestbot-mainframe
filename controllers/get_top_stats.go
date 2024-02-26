package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

// METHOD: GET
// PATH: /top-statistics
// QUERIES: username, server, limit, statistic
// RESPONSE: JSON
// DESCRIPTION: Get the top statistics for a user on a server
// example http://localhost:5000/api/v1/topstatistics?username=Febzey&server=newtest_new1&limit=5&statistic=playtime
func (c *Controller) GetTopStatistics(w http.ResponseWriter, r *http.Request) {

	//get the query for username
	server := r.URL.Query().Get("server")
	limit := r.URL.Query().Get("limit")
	statistic := r.URL.Query().Get("statistic")

	if server == "" || limit == "" || statistic == "" {
		http.Error(w, "Invalid 'server' AND 'limit' AND 'statistic' parameter required", http.StatusBadRequest)
		return
	}

	// Ensure the user-specified statistic is a valid column name to prevent SQL injection
	validStatistics := map[string]bool{"playtime": true, "joins": true, "kills": true, "deaths": true}
	if !validStatistics[statistic] {
		http.Error(w, "Invalid 'statistic' parameter", http.StatusBadRequest)
		return
	}

	// Construct the SQL query dynamically
	SELECT_TOP_STATISTICS_QUERY := fmt.Sprintf("SELECT username, %s FROM users WHERE mc_server = ? ORDER BY %s DESC LIMIT ?", statistic, statistic)

	rows, err := c.Database.Query(SELECT_TOP_STATISTICS_QUERY, server, limit)
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	defer rows.Close()

	type TopStatistic struct {
		Username  string      `json:"username"`
		Statistic interface{} `json:"statistic"`
	}

	var topStatistics []TopStatistic

	for rows.Next() {
		var ts TopStatistic
		var stat sql.RawBytes

		err := rows.Scan(&ts.Username, &stat)
		if err != nil {
			http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
			c.Logger.Error(err.Error())
			return
		}

		// Convert RawBytes to the appropriate data type based on the selected statistic
		switch statistic {
		case "playtime", "joins", "kills", "deaths":
			ts.Statistic, err = strconv.Atoi(string(stat))
			if err != nil {
				http.Error(w, "Internal Database Error - Unable to convert statistic to integer", http.StatusInternalServerError)
				c.Logger.Error(err.Error())
				return
			}
		// Add additional cases if needed for other data types

		default:
			http.Error(w, "Internal Database Error - Unknown statistic type", http.StatusInternalServerError)
			return
		}

		topStatistics = append(topStatistics, ts)
	}

	utils.RespondWithJSON(w, http.StatusOK, topStatistics)

	// //SELECT username, description FROM whois WHERE username=? is our query

	// rows, err := c.Database.Query("SELECT username, description FROM whois WHERE username=?", username)
	// if err != nil {
	// 	http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
	// 	c.Logger.Error(err.Error())
	// 	return
	// }

	// //we need to create a slice of strings to store the usernames in.
	// var usernames []string
	// var descriptions []string

	// //we need to loop through the rows and append the usernames to the slice.
	// for rows.Next() {
	// 	var username string
	// 	var description string
	// 	err := rows.Scan(&username, &description)
	// 	if err != nil {
	// 		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
	// 		c.Logger.Error(err.Error())
	// 		return
	// 	}
	// 	usernames = append(usernames, username)
	// 	descriptions = append(descriptions, description)
	// }

	// //we need to check if the slice is empty.
	// if len(usernames) == 0 {
	// 	http.Error(w, "No usernames found", http.StatusNotFound)
	// 	return
	// }

	// utils.RespondWithJSON(w, http.StatusOK, descriptions)

}
