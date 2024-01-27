package routes

import (
	"database/sql"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

type Guild struct {
	Guild_id   string `json:"guild_id"`
	Channel_id sql.NullString
	Mc_server  string `json:"mc_server"`
	Setup_by   string `json:"setup_by"`
	Created_at int    `json:"created_at"`
	Guild_name string `json:"guild_name"`
}

func (c *Controller) GetDiscordGuilds(w http.ResponseWriter, r *http.Request) {

	rows, err := c.Database.Query("SELECT * FROM guilds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var guilds []Guild
	for rows.Next() {

		var guild Guild
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
