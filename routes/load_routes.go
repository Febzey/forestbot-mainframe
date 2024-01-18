package routes

import (
	"net/http"
	"sync"

	"github.com/febzey/ForestBot-Mainframe/database"
	"github.com/febzey/ForestBot-Mainframe/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Controller struct {
	Database    *database.Database
	Logger      *logger.Logger
	MessageChan chan interface{}
	Clients     map[string]*websocket.Conn
	Mutex       sync.Mutex
}

func LoadAndHandleRoutes(router *mux.Router, db *database.Database, logger *logger.Logger) {

	//main controller that is passed to all of our routes.
	controller := &Controller{
		Database:    db,
		Logger:      logger,
		MessageChan: make(chan interface{}),
		Clients:     make(map[string]*websocket.Conn),
		Mutex:       sync.Mutex{},
	}

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
