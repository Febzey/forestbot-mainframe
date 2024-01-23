package routes

import (
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
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// might put player lists for each server in here also
// or create a new struct for it

//lets make a map of minecraft servers and their player lists using the types.Player struct. Player is a struct that is defined in types/user_types.go

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
	ClientID string
	Message  types.WebsocketMessage
}

type Controller struct {
	Database    *database.Database
	Logger      *logger.Logger
	MessageChan chan MessageChannel
	Clients     map[string]*WebsocketClient
	PlayerLists map[string][]types.Player
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
		Mutex:       sync.Mutex{},
	}
	//this is a sick ideayes
	//so this function will run always :D //always watching the MessageChan
	//lets see if it works!
	go ProcessWebsocketMessage(controller)

	var apiUrl = "/api/v1"

	var routes = []Route{

		//we should be able to get most stats just from the get user method.
		//Like joindate, kills, deaths, lastdeath, lastseen, etc.

		//Get Request to implement:
		//Advancements -- Done
		//Messages
		//Random Quote
		//Tablist
		//Bulk deaths
		//Bulk kills

		//THIS REALLY IS NOT HARD STOP OVER THINKING AND JUST GET IT FUCKING DONE!
		//STOP PROCASTONATING

		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/userbyname/{name}/{server}",
			HandlerFunc: controller.GetUserByName,
		},
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/player/{uuid}/{server}",
			HandlerFunc: controller.GetUserByUUID,
		},
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/websocket/connect",
			HandlerFunc: controller.handleWebSocketAuth,
		},
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/advancements",
			HandlerFunc: controller.getAdvancements,
		},
		{
			Method:      http.MethodGet,
			Pattern:     apiUrl + "/messages",
			HandlerFunc: controller.GetMessages,
		},
	}

	for _, route := range routes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}

}
