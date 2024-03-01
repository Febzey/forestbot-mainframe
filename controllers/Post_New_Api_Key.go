package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

type NewApiKeyRequest struct {
	ContactEmail string `json:"contactEmail"`
	Permissions  struct {
		Write bool `json:"write"`
		Read  bool `json:"read"`
	} `json:"Permissions"`
	RateLimit string `json:"rateLimit"`
	TokenType string `json:"tokentype"`
}

func (c *Controller) PostNewApiKey(w http.ResponseWriter, r *http.Request) {

	var req NewApiKeyRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		// Handle decoding error
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	rateLimit, err := strconv.Atoi(req.RateLimit)
	if err != nil {
		http.Error(w, "error converting ratelimit to int", http.StatusInternalServerError)
		return
	}

	plainTextKey, err := c.KeyService.NewApiKey(req.Permissions.Write, req.Permissions.Read, req.ContactEmail, rateLimit, req.TokenType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Logger.Success(fmt.Sprintf("Created and saved new API Key for email: %s, Read: %t, Write: %t ", req.ContactEmail, req.Permissions.Read, req.Permissions.Write))

	w.Header().Set("Content-Type", "application/json")

	//TODO:: instead of sending key back to client as just a string,
	// lets structure the response like { key: string }
	utils.RespondWithJSON(w, http.StatusOK, plainTextKey)

	return

}
