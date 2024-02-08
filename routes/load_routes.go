package routes

import (
	"image"
	"net/http"
	"sync"

	"github.com/febzey/ForestBot-Mainframe/database"
	"github.com/febzey/ForestBot-Mainframe/logger"
	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Route struct {

	//The HTTP method to use for the route.
	Method string

	//The pattern to use for the route.
	Pattern string

	//The handler function to use for the route.
	HandlerFunc http.HandlerFunc
}

type WebsocketClient struct {

	//The unqie ID for the websocket client connected.
	ClientID string

	//Here is going to be our Permissions. One api key will have all permissions.
	//The other api key will only have read data permissions.
	//Permssion types can be: read_data, write_data
	Permissions utils.APIPermissions

	//The minecraft server the websocket is being used for.
	Mc_server string

	//The websocket connection for the client.
	Conn *websocket.Conn
}

type MessageChannel struct {

	//The websocket client ID for the message.
	ClientID string

	//The websocket message for the message.
	Message types.WebsocketMessage
}

type Controller struct {
	Database    *database.Database
	Logger      *logger.Logger
	MessageChan chan MessageChannel
	Clients     map[string]*WebsocketClient
	PlayerLists map[string][]types.Player
	ImageCache  types.ImageCache
	Mutex       sync.Mutex
}

func LoadAndHandleRoutes(router *mux.Router, db *database.Database, logger *logger.Logger) {

	//main controller that is passed to all of our routes.
	controller := &Controller{
		Database:    db,
		Logger:      logger,
		MessageChan: make(chan MessageChannel),
		Clients:     make(map[string]*WebsocketClient),
		PlayerLists: make(map[string][]types.Player),
		ImageCache: types.ImageCache{
			HeadImages: make(map[string]image.Image),
		},
		Mutex: sync.Mutex{},
	}

	//this is a sick ideayes
	//so this function will run always :D //always watching the MessageChan
	//lets see if it works!
	go ProcessWebsocketMessage(controller)

	var apiUrl = "/api/v1"

	var routes = []Route{

		//Get Request to implement:
		//Advancements -- Done
		//Messages -- Done
		//Random Quote -- Done
		//Tablist -- Done
		//Bulk deaths -- DONE
		//Bulk kills -- done

		//queries: server username or uuid
		//Description: Gets the player activity data for a server
		//example url: http://localhost:5000/api/v1/player-activity-weekly-report?server=simplyvanilla?username=febzey
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/specific-player-activity-weekly-report",
			HandlerFunc: controller.GetPlayerActivityData,
		},

		//queries: server
		//Description: Gets the player activity by week day
		//example url: http://localhost:5000/api/v1/player-activity-by-week-day?server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/player-activity-by-week-day",
			HandlerFunc: controller.GetPlayerActivityByWeekDay,
		},

		//queries: server
		//Description: Gets the player activity by hour
		//example url: http://localhost:5000/api/v1/player-activity-by-hour?server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/player-activity-by-hour",
			HandlerFunc: controller.GetHourlyPlayerActivityStats,
		},

		//queries: username
		//Description: Gets the player statistics for a user for all servers theyve been see on
		//example url: http://localhost:5000/api/v1/all-player-stats?username=febzey
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/all-player-stats",
			HandlerFunc: controller.GetAllPlayerStatisticsByUsername,
		},

		//Quries: name, server
		//Description: Gets a user by their name
		//example url: http://localhost:5000/api/v1/playername?name=febzey&server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/playername",
			HandlerFunc: controller.GetUserByName,
		},
		//Quries: uuid, server
		//Description: Gets a user by their UUID
		//example url: http://localhost:5000/api/v1/playeruuid?uuid=30303-addwdwd-222=3333&server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/playeruuid",
			HandlerFunc: controller.GetUserByUUID,
		},

		//This is a websocket for handling data between the server and the client.
		//This is used for getting data from the server in real time.
		//This is also used for sending data to the server in real time.
		//playtime, chat other things are sent through this websocket.
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/websocket/connect",
			HandlerFunc: controller.handleWebSocketAuth,
		},

		//Quries: uuid, server, limit, order
		//Description: Gets the advancements of a player
		//example url: http://localhost:5000/api/v1/advancements?uuid=30303-addwdwd-222=3333&server=simplyvanilla&limit=100&order=DESC
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/advancements",
			HandlerFunc: controller.getAdvancements,
		},

		//Quries: name, server, limit, order
		//Description: Gets the messages for a player
		//example url: http://localhost:5000/api/v1/messages?name=febzey&server=simplyvanilla&limit=100&order=DESC
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/messages",
			HandlerFunc: controller.GetMessages,
		},
		//Quries: name, server
		//Description: Get a random quote from a user on a server
		//example url: http://localhost:5000/api/v1/quote?name=febzey&server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/quote",
			HandlerFunc: controller.GetRandomQuotes,
		},
		//Quries: server
		//Description: Get the tablist for a server
		//example url: http://localhost:5000/api/v1/tablist?server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/tablist",
			HandlerFunc: controller.GetTablist,
		},
		//Quries: uuid, server, limit, order, type (pvp, pve, all)
		//Description: Gets the deaths of a player
		//example url: http://localhost:5000/api/v1/deaths?uuid=30303-addwdwd-222=3333&server=simplyvanilla&limit=100&order=DESC
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/deaths",
			HandlerFunc: controller.GetMinecraftDeaths,
		},
		//Quries: uuid, server, limit, order
		//Description: Gets the kills of a player
		//example url: http://localhost:5000/api/v1/kills?uuid=30303-addwdwd-222=3333&server=simplyvanilla&limit=100&order=DESC

		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/kills",
			HandlerFunc: controller.GetMinecraftKills,
		},
		//queries username
		//description: checks if a user is online or not and returns server and true or false
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/online",
			HandlerFunc: controller.GetUserOnlineCheck,
		},

		//queries username, server
		//description: returns back 6 closest names to the username
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/namesearch",
			HandlerFunc: controller.GetNameSearch,
		},

		//queries username
		//description: returns back the whois data for the username
		//example url: http://localhost:5000/api/v1/whois?username=febzey
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/whois",
			HandlerFunc: controller.GetWhoIs,
		},

		//queries username
		//description: returns back the UUID for the username
		//example url: http://localhost:5000/api/v1/convert-username-to-uuid?username=febzey
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/convert-username-to-uuid",
			HandlerFunc: controller.ConvertUsernameToUUID,
		},

		//queries username, server
		//description: returns back the playtime for the username
		//example url: http://localhost:5000/api/v1/messagecount?username=febzey&server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/messagecount",
			HandlerFunc: controller.GetUserMessageCount,
		},

		//queries username, server, word
		//description: returns back the word count for the username
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/wordcount",
			HandlerFunc: controller.GetUserWordOccurence,
		},

		//get all the top statistics for a server
		//queries server, limit, statistic, (stastic can be playtime, joins, kills, deaths)
		//example url: http://localhost:5000/api/v1/topstatistics?server=simplyvanilla&limit=5&statistic=playtime
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/top-statistic",
			HandlerFunc: controller.GetTopStatistics,
		},

		//Get all the guilds forestbot is in for discord
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/discord/guilds",
			HandlerFunc: controller.GetDiscordGuilds,
		},

		//get all live chat channels for our discord bot
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/discord/livechats",
			HandlerFunc: controller.GetDiscordLiveChatChannels,
		},

		//////POST REQUESTS//////

		//body: {"guild_id": "123", "channel_id": "123", "mc_server": "simplyvanilla", "setup_by": "123", "created_at": "123", "guild_name": "simplyvanilla_discord_server"}
		//description: adds a guild to the database
		//example url: http://localhost:5000/api/v1/discord/addguild
		{
			Method:      http.MethodPost,
			Pattern:     apiUrl + "/discord/addguild",
			HandlerFunc: controller.PostDiscordGuild,
		},
		//body: {"guildName": "simplyvanilla_discord_server", "guildID": "123", "channelID": "123", "setupBy": "123", "date": "123", "mcServer": "simplyvanilla"}
		//description: adds a live chat channel to the database
		//example url: http://localhost:5000/api/v1/discord/addlivechat
		{
			Method:      http.MethodPost,
			Pattern:     apiUrl + "/discord/addlivechat",
			HandlerFunc: controller.PostDiscordLiveChat,
		},

		//body: {"username": "febzey", "description": "I am a cool guy"}
		//description: Sets the description of a user
		//example url: http://localhost:5000/api/v1/whois_description
		{
			Method:      http.MethodPost,
			Pattern:     apiUrl + "/whois-description",
			HandlerFunc: controller.POSTIamDescription,
		},

		/////// DELETE REQUESTS ///////

		//queries: guild_id
		//description: deletes a guild from the database
		//example url: http://localhost:5000/api/v1/discord/deleteguild?guild_id=123
		{
			Method:      http.MethodDelete,
			Pattern:     apiUrl + "/discord/deleteguild",
			HandlerFunc: controller.DeleteDiscordGuild,
		},

		//queries: channel_id
		//description: deletes a live chat channel from the database
		//example url: http://localhost:5000/api/v1/discord/deletelivechat?channel_id=123
		{
			Method:      http.MethodDelete,
			Pattern:     apiUrl + "/discord/deletelivechat",
			HandlerFunc: controller.DeleteDiscordLiveChat,
		},
	}

	for _, route := range routes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}

}
