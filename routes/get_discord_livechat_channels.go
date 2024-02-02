package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
)

func (c *Controller) GetDiscordLiveChatChannels(w http.ResponseWriter, r *http.Request) {

	rows, err := c.Database.Query("SELECT * FROM livechats")
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Internal Database Error.", http.StatusInternalServerError)
		return
	}

	var livechats []types.LivechatChannel
	for rows.Next() {

		var livechat types.LivechatChannel
		err = rows.Scan(
			&livechat.GuildName,
			&livechat.GuildID,
			&livechat.ChannelID,
			&livechat.Setupby,
			&livechat.Date,
			&livechat.Mc_server,
		)
		if err != nil {
			c.Logger.Error(err.Error())
			http.Error(w, "Internal Database Error.", http.StatusInternalServerError)
			return
		}
		livechats = append(livechats, livechat)
	}

	utils.RespondWithJSON(w, http.StatusOK, livechats)

}
