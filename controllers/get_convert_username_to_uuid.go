package controllers

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) ConvertUsernameToUUID(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Invalid 'username' parameter required", http.StatusBadRequest)
		return
	}

	uuid, err := c.Database.ConvertUsernameToUUID(username)
	if err != nil {
		http.Error(w, "Internal database error - contact febzey", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	// var GET_UUID_FROM_USERNAME_QUERY = "SELECT DISTINCT uuid FROM users WHERE username = ?"
	// rows, err := c.Database.Query(GET_UUID_FROM_USERNAME_QUERY, username)
	// if err != nil {
	// 	http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
	// 	c.Logger.Error(err.Error())
	// 	return
	// }

	// defer rows.Close()

	// type UUID struct {
	// 	UUID sql.NullString `json:"uuid"`
	// }

	// uuid := UUID{}

	// for rows.Next() {
	// 	err := rows.Scan(&uuid.UUID)
	// 	if err != nil {
	// 		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
	// 		c.Logger.Error(err.Error())
	// 		return
	// 	}
	// }

	utils.RespondWithJSON(w, http.StatusOK, uuid)

}
