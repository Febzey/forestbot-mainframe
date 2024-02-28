package controllers

import (
	"image"
	"net/http"
	"sync"

	"github.com/febzey/ForestBot-Mainframe/database"
	"github.com/febzey/ForestBot-Mainframe/keyservice"
	"github.com/febzey/ForestBot-Mainframe/logger"
	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/gorilla/mux"
)

var (
	apiUrl   = "/api/v1"
	head_url = "https://mc-heads.net/avatar/"
)

type Route struct {

	//The HTTP method to use for the route.
	Method string

	//The pattern to use for the route.
	Pattern string

	//The handler function to use for the route.
	HandlerFunc http.HandlerFunc
}

type Controller struct {
	//Our main database instance with helper functions.
	Database *database.Database

	//Logger utility function for nice console logging.
	Logger *logger.Logger

	//a message channel where each websocket message is sent to.
	//Global channel for anything websocket.
	MessageChan chan MessageChannel

	//List of connected websocket clients
	//key is their unique ID given when they connect.
	Clients map[string]*WebsocketClient

	//event handlers for websocket.
	//string will be action type for the event.
	Handlers map[string]Handler

	//Player list for each server connected
	//key would be the name of server
	//Array is player
	PlayerLists map[string][]types.Player

	//Caching images for playerlist / tablist
	ImageCache types.ImageCache

	//Key service for authentication
	KeyService keyservice.APIKeyService

	//a mutex to keep our Controller in sync.
	Mutex sync.Mutex
}

func NewController(db *database.Database, logger *logger.Logger, keyService *keyservice.APIKeyService) *Controller {
	return &Controller{
		Database:    db,
		Logger:      logger,
		MessageChan: make(chan MessageChannel),
		Clients:     make(map[string]*WebsocketClient),
		Handlers:    make(map[string]Handler),
		PlayerLists: make(map[string][]types.Player),
		ImageCache: types.ImageCache{
			HeadImages: make(map[string]image.Image),
		},
		KeyService: *keyService,
		Mutex:      sync.Mutex{},
	}
}

func LoadAndHandleRoutes(router *mux.Router, controller *Controller) {

	controller.setupWebsocketEventHandlers()

	//Continous running function that processes all Websocket Messages.
	go ProcessWebsocketEvent(controller)

	var routes = []Route{
		//Gets all available servers forestbot has been on
		//example url: http://localhost:5000/api/v1/all-server
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/all-servers",
			HandlerFunc: controller.GetAvailableServers,
		},

		//queries: server username or uuid
		//Description: Gets the player activity data for a server
		//example url: http://localhost:5000/api/v1/player-activity-weekly-report?server=simplyvanilla?username=febzey
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/specific-player-activity-weekly-report",
			HandlerFunc: controller.GetPlayerActivityData,
		},

		//queries: server
		//Description: Gets the total number of players logged in on each day of the week for a specific server.
		//example url: http://localhost:5000/api/v1/player-activity-by-week-day?server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/player-activity-by-week-day",
			HandlerFunc: controller.GetPlayerActivityByWeekDay,
		},

		//queries: server
		//Description: Gets the player count for a specific server by the hour and weekday.
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
			HandlerFunc: controller.websocketController,
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
