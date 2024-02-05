package routes

import (
	"encoding/json"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

// METHOD: POST
// PATH: /iam_description
// BODY: JSON
// RESPONSE: JSON
// DESCRIPTION: Sets the description of a user
// example: http://localhost:5000/api/v1/whois_description
func (c *Controller) POSTIamDescription(w http.ResponseWriter, r *http.Request) {

	var IamBody struct {
		Username    string `json:"username"`
		Description string `json:"description"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&IamBody)
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	err = c.Database.INSERT_player_whois_description(IamBody.Username, IamBody.Description)
	if err != nil {
		http.Error(w, "Internal Database Error - Please contact Febzey on Discord", http.StatusInternalServerError)
		c.Logger.Error(err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, IamBody)
}
