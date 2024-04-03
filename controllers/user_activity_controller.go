/***

This controller handles any endpoints for interacting with our
playerActivity table in the database.

***/

package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

/**
+------------+-------------+--------------+
| user_count | day_of_week | hour_of_day  |
+------------+-------------+--------------+
| 5          | 1           | 12           |  -- Sunday, 12:00 AM
| 8          | 1           | 13           |  -- Sunday, 1:00 PM
| 10         | 2           | 14           |  -- Monday, 2:00 PM
| ...        | ...         | ...          |
+------------+-------------+--------------+

We use this function for our websites specific server graph.
**/

func (c *Controller) GetHourlyServerActivityStats(w http.ResponseWriter, r *http.Request) {

	mc_server := r.URL.Query().Get("server")

	playerActivityByHour, err := c.Database.ServerActivityHourlyResults(mc_server)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	responseJSON, err := json.Marshal(playerActivityByHour)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}

/**
 +--------------+-------------+
 | COUNT(UUID)  | day_of_week |
 +--------------+-------------+
 | 15           | 1           |  -- Sunday
 | 25           | 2           |  -- Monday
 | 20           | 3           |  -- Tuesday
 | 18           | 4           |  -- Wednesday
 | 30           | 5           |  -- Thursday
 | 22           | 6           |  -- Friday
 | 17           | 7           |  -- Saturday
 +--------------+-------------+
 In this example
 On Sunday (day_of_week = 1), 15 players logged in.
 On Monday (day_of_week = 2), 25 players logged in.

 A more spread our version of the above function. instead of hourly. it is weekly.
**/

func (c *Controller) GetPlayerActivityByWeekDay(w http.ResponseWriter, r *http.Request) {

	mc_server := r.URL.Query().Get("server")

	playerActivityByWeekDay, err := c.Database.PlayerActivityWeekResults(mc_server)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	response := map[string]interface{}{
		"player_activity_by_week_day": playerActivityByWeekDay,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}

/*
*
| UUID            | Date               | type   |
|-----------------|--------------------|--------|
| abc123          | 1707673202084      | login  |
| def456          | 1707673195801      | login  |
| ghi789          | 1707673126855      | logout |
| jkl012          | 1707673100524      | login  |
| ...             | ...                | ...    |
Getting raw activity data for a specific user on a specific minecraft server
useful for creating custom graphs for a user.
*
*/
func (c *Controller) GetPlayerActivityData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mc_server := r.URL.Query().Get("server")
	username := r.URL.Query().Get("username")
	uuid := r.URL.Query().Get("uuid")

	usingUuid := false

	if username == "" && uuid == "" {
		http.Error(w, "Invalid 'username' or 'uuid' parameter required.", http.StatusBadRequest)
		return
	}

	if mc_server == "" {
		http.Error(w, "Invalid 'server' parameter required.", http.StatusBadRequest)
		return
	}

	if uuid != "" {
		usingUuid = true
		username = uuid
	}

	playerActivityData, err := c.Database.GetAllPlayerActivity(mc_server, username, usingUuid)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, playerActivityData)

}

// // GetServerStats is a function that returns the weekly server stats for a specific server.
// func (c *Controller) GetServerStats(w http.ResponseWriter, r *http.Request) {
// 	mc_server := r.URL.Query().Get("server")

// 	serverStats, err := c.Database.SELECT_server_stats(mc_server)
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		c.Logger.Error(err.Error())
// 		return
// 	}

// 	responseJSON, err := json.Marshal(serverStats)
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		c.Logger.Error(err.Error())
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(responseJSON)

// }

// getting top 5 leaderboards for various things from the past 7 days exactly.
func (c *Controller) GetTop5Leaderboard(w http.ResponseWriter, r *http.Request) {
	mc_server := r.URL.Query().Get("server")

	top5Leaderboards, err := c.Database.SELECT_top_5_player_stats(mc_server)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	responseJSON, err := json.Marshal(top5Leaderboards)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}

func (c *Controller) GetServerTotalSavedDataCount(w http.ResponseWriter, r *http.Request) {

	mc_server := r.URL.Query().Get("server")

	totalSavedDataCount, err := c.Database.SELECT_server_stats_total_overall(mc_server)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	responseJSON, err := json.Marshal(totalSavedDataCount)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}
