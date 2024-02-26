package controllers

import (
	"fmt"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

//METHOD: GET
//PATH: /tablist
//QUERIES: server
//RESPONSE: ARRAY
//Example: http://localhost:5000/api/v1/tablist?server=newtest_new1

// PlayerLists map[string][]types.Player
// string will be server name, []types.Player will be the list of players on that server
func (c *Controller) GetTablist(w http.ResponseWriter, r *http.Request) {
	server := r.URL.Query().Get("server")

	if server == "" {
		http.Error(w, "Invalid 'server' parameter required", http.StatusBadRequest)
		return
	}

	playerList := c.PlayerLists[server]
	if playerList == nil {
		http.Error(w, "Server does not have an active player list", http.StatusBadRequest)
		return
	}

	dc := utils.RenderTab(playerList, &c.ImageCache)

	filePath := "tablists/" + server + ".png"

	if err := dc.SavePNG(filePath); err != nil {
		fmt.Println("Error saving PNG:", err)
		return
	}

	//send the image to the client
	http.ServeFile(w, r, filePath)

}
