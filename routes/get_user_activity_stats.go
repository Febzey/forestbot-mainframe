package routes

import (
	"encoding/json"
	"net/http"
)

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
