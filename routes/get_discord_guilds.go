package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetDiscordGuilds(w http.ResponseWriter, r *http.Request) {

	rows, err := c.Database.Query("SELECT * FROM guilds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var guilds []types.Guild
	for rows.Next() {

		var guild types.Guild
		err = rows.Scan(
			&guild.Guild_id,
			&guild.Channel_id,
			&guild.Mc_server,
			&guild.Setup_by,
			&guild.Created_at,
			&guild.Guild_name,
		)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, "Error with database", http.StatusInternalServerError)
			return
		}
		guilds = append(guilds, guild)
	}

	utils.RespondWithJSON(w, http.StatusOK, guilds)

}
