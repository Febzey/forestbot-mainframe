package routes

import (
	"encoding/json"
	"net/http"
)

func (c *Controller) GetHourlyPlayerActivityStats(w http.ResponseWriter, r *http.Request) {

	mc_server := r.URL.Query().Get("server")

	playerActivityByHour, err := c.Database.PlayerActivityHourlyResults(mc_server)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	responseJSON, err := json.Marshal(map[string]interface{}{"player_activity_by_hour": playerActivityByHour})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	// response := map[string]interface{}{
	// 	"player_activity_by_hour": playerActivityByHour,
	// }

	// responseJSON, err := json.Marshal(response)
	// if err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	c.Logger.Error(err.Error())
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}
