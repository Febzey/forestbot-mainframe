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

		//Quries: name, server
		//Description: Gets a user by their name
		//example url: http://localhost:5000/api/v1/userbyname?name=febzey&server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/userbyname",
			HandlerFunc: controller.GetUserByName,
		},
		//Quries: uuid, server
		//Description: Gets a user by their UUID
		//example url: http://localhost:5000/api/v1/userbyuuid?uuid=30303-addwdwd-222=3333&server=simplyvanilla
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/player",
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
		//Quries: uuid, server, limit, order
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
	}

	for _, route := range routes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}

}
