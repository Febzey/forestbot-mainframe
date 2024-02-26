/******

	This is the websocket controller.
	Handling everything websocket.

******/

package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
)

/**
Main message Channel, all inbound client messages will defined
by this struct.
**/
type MessageChannel struct {

	//The websocket client ID for the message.
	//The generated uuid will be stored here, not the uuid the client sends.
	//The uuid that we recieve from the client can be found in the 'Message' struct.
	ClientID string

	//The websocket message for the message.
	Message WebsocketEvent
}

/*
Function used for removing a websocket client from our Controller Client map.
*/
func (c *Controller) removeWebSocketClient(clientID string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if client, ok := c.Clients[clientID]; ok {
		if client.Conn != nil {
			client.Conn.Close()
		}

		delete(c.Clients, clientID)
	}
}

/**
This is our websocket controller function,
this will go off anytime someone connects to the websocket,
it is a http endpoint but is upgraded to websocket.
**/
func (c *Controller) websocketController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//
	//Upgrading http connection to websocket
	//
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(err.Error())
		return
	}

	//
	//Getting our server and x-api-key url queries.
	//
	mc_server, api_key, isBot := r.URL.Query().Get("server"), r.URL.Query().Get("x-api-key"), r.URL.Query().Get("is-bot-client")

	//
	//Creating a websocket instance with or given queries.
	//
	client := NewWebsocketClient(conn, api_key, mc_server, isBot, c)
	if client == nil {
		c.Logger.WebsocketError(fmt.Sprintf("Client failed to connect."))
		return
	}

	//
	//Let the console know a websocket client has connected.
	//
	c.Logger.WebsocketConnect(fmt.Sprintf("Websocket Client Connected For Minecraft Server: %s | ID: %s | isMcBot: %t", mc_server, client.ClientID, isBot == "true"))

	go client.readMessages()
	go client.writeMessages()

}

/*
Removing a specific user from a specific servers c.PlayerLists
if the user and the server are present in the map
*/
func (c *Controller) removeUserFromPlayerList(serverName string, usernameToRemove string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if playerList, ok := c.PlayerLists[serverName]; ok {
		var updatedPlayerList []types.Player

		for _, player := range playerList {
			if player.Username != usernameToRemove {
				updatedPlayerList = append(updatedPlayerList, player)
			}
		}

		c.PlayerLists[serverName] = updatedPlayerList
	}
}

/*
Adding a user to the c.PlayerLists map.
If the server doesnt exist as a key in c.PlayerLists,
we will create a new one and add the user.
*/
func (c *Controller) addUserToPlayerList(serverName string, newPlayer types.Player) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if playerList, ok := c.PlayerLists[serverName]; ok {
		for i, player := range playerList {
			if player.Username == newPlayer.Username {
				c.PlayerLists[serverName][i] = newPlayer
				return
			}
		}

		c.PlayerLists[serverName] = append(playerList, newPlayer)
	} else {
		c.PlayerLists[serverName] = []types.Player{newPlayer}
	}
}

/*
With this function we are able to send a message
to every websocket client connected to our server by
sending a message to each of our connected clients egress channels
*/
func (c *Controller) BroadcastMessageToClients(message WebsocketEvent) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, client := range c.Clients {
		if client.Permissions.Read {
			client.Egress <- message
		}
	}
}

/*
This function allows us to send a message to specific
websocket connection, while following our websocket message structure.
*/
func (c *Controller) sendMessageByStructure(id string, message WebsocketEvent) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	//change to send to the clients egress. then the egress handler will handle the msg
	client, ok := c.Clients[id]
	if !ok {
		return errors.New("No client found")
	}
	client.Egress <- message
	return nil
}

func (c *Controller) sendErrorMessage(id string, message string) {
	c.sendMessageByStructure(id, WebsocketEvent{
		Client_id: id,
		Action:    "error",
		Data:      message,
	})
}
