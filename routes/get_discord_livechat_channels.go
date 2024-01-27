package routes

import (
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/utils"
)

type LivechatChannels struct {
	GuildName string `json:"guildName"`
	GuildID   string `json:"guildID"`
	ChannelID string `json:"channelID"`
	Setupby   string `json:"setupBy"`
	Date      string `json:"date"`
	Mc_server string `json:"mcServer"`
}

func (c *Controller) GetDiscordLiveChatChannels(w http.ResponseWriter, r *http.Request) {

	rows, err := c.Database.Query("SELECT * FROM livechats")
	if err != nil {
		c.Logger.Error(err.Error())
		http.Error(w, "Internal Database Error.", http.StatusInternalServerError)
		return
	}

	var livechats []LivechatChannels
	for rows.Next() {

		var livechat LivechatChannels
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
