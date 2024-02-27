package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
	"github.com/gorilla/websocket"
)

/*
*
A struct for our individual websocket clients.
created when a client connects to our server.
*
*/
type WebsocketClient struct {

	//The unqie ID for the websocket client connected.
	ClientID string

	//Here is going to be our Permissions. One api key will have all permissions.
	//The other api key will only have read data permissions.
	//Permssion types can be: read_data, write_data
	Permissions types.APIPermissions

	//The minecraft server the websocket is being used for.
	//this will only be populated if the client connecting is a bot-client
	//if a client wants messages for a specific server, they can just filter for the server they need.
	//since all messages are sent from one pipeline anywas
	Mc_server string

	//Perhaps add a list of servers a regular client would want to listen for?
	//if the option is "all", we will loop through active connections, and add each mc_server to this list
	//or the client can specifiy individual servers to listen for, we would then do a check to see if that server exists.
	//Design choice: Option A: add this as a optional query parameter. Option B: add a event option the user can send later as a "command".
	//we will implement this if the websocket server becomes overwhelemed with broadcasts. (which it will with our current setup.)
	//ListeningServers []string

	//Determine if this client is an active minecraft bot.
	//Mc Clients are essential for this entire project,
	//they will take first priority when needed.
	//these are the actual bots that will be in minecraft servers collecting data
	//There can only be 1 mc-client to 1 mc_server
	//if a user wants to read messages for a certain mc server without being a bot,
	//simply just do not register your self as a bot client in the url queries
	IsMcClient bool

	//The websocket connection for the client.
	Conn *websocket.Conn

	//egress channel for outbound websocket messages (messages that we send back to the client)
	//we have a go routine running that listens to this channel and will send messages accordingly.
	Egress chan WebsocketEvent

	//A pointer to our programs Controller
	Controller *Controller
}

/*
Some settings for our websocket behaviour
*/
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

/*
*

initializing our websocket client
when the user connects, this function will
generate their unique client_id, check permissions based on the given API key,
then send the client_id back to the user for the user to store and send with each message.

*
*/
func NewWebsocketClient(conn *websocket.Conn, api_key string, mc_server string, isBot string, c *Controller) *WebsocketClient {

	//
	//Incoming connection is missing either server or x-api-key queries

	//we will want to switch to check for the api_key in a map.
	//
	if api_key == "" {
		err := conn.WriteJSON(WebsocketEvent{
			Client_id: "",
			Action:    "error",
			Data:      "Missing required headers: x-api-key - closing the connection",
		})
		if err != nil {
			c.Logger.WebsocketError("Error sending message to client after incorrect queries provided.")
		}
		conn.Close()
		return nil
	}

	//
	//Checking the clients permissions based of the given api key
	//
	client_permissions, err := utils.CheckApiKey(api_key)
	if err != nil {
		if err := conn.WriteJSON(WebsocketEvent{
			Client_id: "",
			Action:    "error",
			Data:      "You provided an invalid api key",
		}); err != nil {
			c.Logger.WebsocketError("client sent invalid api key, but we failed to tell them that.")
		}
		conn.Close()
		return nil
	}

	//
	//Checking if the connecting client is a mc bot
	//
	if isBot == "true" {
		if mc_server == "" {
			if err := conn.WriteJSON(WebsocketEvent{
				Client_id: "",
				Action:    "error",
				Data:      "When registering as a bot-client, you must provide a server your bot associates with.",
			}); err != nil {
				c.Logger.WebsocketError("Error sending message to client.")
			}
			conn.Close()
			return nil
		}
		for _, client := range c.Clients {
			if client.Mc_server == mc_server && client.IsMcClient {
				//a user is trying to connect as a bot client, but a
				//bot client already exists for this mc_server
				if err := conn.WriteJSON(WebsocketEvent{
					Client_id: "",
					Action:    "error",
					Data:      "A mc client already exists and is running for this minecraft server, if you still want to listen to the traffic, then take away the is-bot-client from your query",
				}); err != nil {
					c.Logger.WebsocketError("Error sending message to client.")
				}
				conn.Close()
				return nil
			}
		}
	}

	//
	//Generating unique client id
	//
	client_id, err := utils.RandomUUID()
	if err != nil {
		fmt.Println(err.Error())
		err := conn.WriteJSON(WebsocketEvent{
			Client_id: client_id,
			Action:    "error",
			Data:      "Error generating client_id - Internal server error",
		})
		if err != nil {
			c.Logger.WebsocketError("Error sending message to client after uuid failed to generate.")
		}
		conn.Close()
		return nil
	}

	//
	//Sending the recenetly generated client_id back to the client so they can store it to use in future messages
	//
	if err := conn.WriteJSON(WebsocketEvent{
		Client_id: client_id,
		Action:    "id",
		Data:      client_id,
	}); err != nil {
		c.Logger.Error("Failed to send generated uuid back to client. Closing connection.")
		conn.Close()
		return nil
	}

	//
	//Creating our custom websocket client instance.
	//
	client := &WebsocketClient{
		ClientID:    client_id,
		Conn:        conn,
		Permissions: client_permissions,
		Mc_server:   mc_server,
		Egress:      make(chan WebsocketEvent),
		Controller:  c,
		IsMcClient:  isBot == "true",
	}

	//
	//Adding the websocket client to our clients map inside of the Controller Struct
	//
	c.Mutex.Lock()
	c.Clients[client_id] = client
	c.Mutex.Unlock()

	//
	//Returning our newly creating client instance since all checks have passed!
	//
	return client

}

/*
Go routine for reading messages per client
for the server to even read their message content, they need to have a proper write api key.
*/
func (ws *WebsocketClient) readMessages() {
	defer func() {
		ws.Controller.removeWebSocketClient(ws.ClientID)
	}()

	for {
		//
		//Read raw incoming data.
		//
		_, p, err := ws.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ws.Controller.Logger.WebsocketDisconnect(fmt.Sprintf("Websocket disconnected: Error: %s, client_id: %s \n", err, ws.ClientID))
			}
			break
		}

		//
		//Checking permissions for our client.
		//Even though ther user is trying to write data with an invalid write key, we still keep the connection open,
		//without forwarding their message event to the proccessor. (this is subject to change)
		//
		if !ws.Permissions.Write {
			ws.Controller.sendMessageByStructure(ws.ClientID, WebsocketEvent{
				Client_id: ws.ClientID,
				Action:    "error",
				Data:      "Client does not have permission to write data",
			})
			continue
		}

		//
		//Attempting to decode the websocket message into our desired message structure.
		//
		var recievedMessage WebsocketEvent
		if err := json.Unmarshal(p, &recievedMessage); err != nil {
			ws.Controller.Logger.WebsocketError(err.Error())

			//
			//send a message back to the client telling them their message structure was invalid
			//
			ws.Controller.sendMessageByStructure(ws.ClientID, WebsocketEvent{
				Client_id: ws.ClientID,
				Action:    "error",
				Data:      "Invalid message structure",
			})
			continue
		}

		//
		//Send data to our websocket message channel
		//for later proccessing
		//
		ws.Controller.MessageChan <- MessageChannel{
			ClientID: ws.ClientID,
			Message:  recievedMessage,
		}

	}
}

/*
Go routine for sending messages.
we will constantly read our egress channel and send messages
if content is picked up in our chanenl, each client will have their own egress.
for now all of our messages will be sent as json
*/
func (ws WebsocketClient) writeMessages() {
	defer func() {
		ws.Controller.removeWebSocketClient(ws.ClientID)
	}()

	for message := range ws.Egress {
		if message.Action == "" {
			if err := ws.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
				break
			}
			return
		}

		if err := ws.Conn.WriteJSON(message); err != nil {
			ws.Controller.Logger.WebsocketError(err.Error())
		}
	}

}
