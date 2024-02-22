package types

// Types for api permissions.
type APIPermissions struct {
	Read  bool
	Write bool
}

// // A Route for an HTTP endpoint.
// type Route struct {

// 	//The HTTP method to use for the route.
// 	Method string

// 	//The pattern to use for the route.
// 	Pattern string

// 	//The handler function to use for the route.
// 	HandlerFunc http.HandlerFunc
// }

// // This is what a connected websocket clients struct will be.
// type WebsocketClient struct {

// 	//The unqie ID for the websocket client connected.
// 	ClientID string

// 	//Here is going to be our Permissions. One api key will have all permissions.
// 	//The other api key will only have read data permissions.
// 	//Permssion types can be: read_data, write_data
// 	Permissions APIPermissions

// 	//The minecraft server the websocket is being used for.
// 	Mc_server string

// 	//The websocket connection for the client.
// 	Conn *websocket.Conn
// }

// // The struct the websockets main message channel will have.
// type MessageChannel struct {

// 	//The websocket client ID for the message.
// 	ClientID string

// 	//The websocket message for the message.
// 	Message WebsocketMessage
// }

// type Controller struct {
// 	//Our main database instance with helper functions.
// 	Database *database.Database

// 	//Logger utility function for nice console logging.
// 	Logger *logger.Logger

// 	//a message channel where each websocket message is sent to.
// 	//Global channel for anything websocket.
// 	MessageChan chan MessageChannel

// 	//List of connected websocket clients
// 	//key is their unique ID given when they connect.
// 	Clients map[string]*WebsocketClient

// 	//Player list for each server connected
// 	//key would be the name of server
// 	//Array is player
// 	PlayerLists map[string][]Player

// 	//Caching images for playerlist / tablist
// 	ImageCache ImageCache

// 	Mutex sync.Mutex
// }
