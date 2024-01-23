package routes

import (
	"net/http"
	"sync"

	"github.com/febzey/ForestBot-Mainframe/database"
	"github.com/febzey/ForestBot-Mainframe/logger"
	"github.com/febzey/ForestBot-Mainframe/types"
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

type Controller struct {
	Database    *database.Database
	Logger      *logger.Logger
	MessageChan chan types.WebsocketMessage
	Clients     map[string]*websocket.Conn
	PlayerLists map[string][]types.Player
	Mutex       sync.Mutex
}

func LoadAndHandleRoutes(router *mux.Router, db *database.Database, logger *logger.Logger) {

	//main controller that is passed to all of our routes.
	controller := &Controller{
		Database:    db,
		Logger:      logger,
		MessageChan: make(chan types.WebsocketMessage),
		Clients:     make(map[string]*websocket.Conn),
		PlayerLists: make(map[string][]types.Player),
		Mutex:       sync.Mutex{},
	}
	//this is a sick ideayes
	//so this function will run always :D //always watching the MessageChan
	//lets see if it works!
	go ProcessWebsocketMessage(controller)

	var routes = []Route{
		{
			Method:      http.MethodGet,
			Pattern:     "/api/v1/userbyname/{name}/{server}",
			HandlerFunc: controller.GetUserByName,
		},
		{
			Method:      http.MethodGet,
			Pattern:     "/api/v1/player/{uuid}/{server}",
			HandlerFunc: controller.GetUserByUUID,
		},
		{
			Method:      http.MethodGet,
			Pattern:     "/authenticate",
			HandlerFunc: controller.handleWebSocketAuth,
		},
	}

	for _, route := range routes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}

}
