package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// MinecraftMessage represents the top-level message structure.

// The Data inside of websocket message can be one of the other types below.
// Type will be either "discord" or "minecraft"
// Action will be like "savechat" "chat" "advancement" "playerjoin" "playerleave" "playerdeath" etc
type WebsocketMessage struct {
	Type   string      `json:"type"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type DiscordInBoundMessage struct {
	GuildName string `json:"guild_name"`
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	SetupBy   string `json:"setup_by"`
	TimeStamp string `json:"time_stamp"`
	MCServer  string `json:"mc_server"`
}

// MinecraftChatMessage represents the structure for chat messages.
type MinecraftChatMessage struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	MCServer string `json:"mc_server"`
	Date     int64  `json:"timestamp"`
}

// MinecraftChatAdvancement represents the structure for advancement messages.
type MinecraftChatAdvancement struct {
	Username    string `json:"username"`
	Advancement string `json:"advancement"`
	Time        int64  `json:"time"`
	MCServer    string `json:"mc_server"`
	UUID        string `json:"uuid,omitempty"`
}

// MinecraftPlayerJoinArgs represents the structure for player join messages.
type MinecraftPlayerJoinArgs struct {
	User     string `json:"user"`
	UUID     string `json:"uuid"`
	MCServer string `json:"mc_server"`
	Time     string `json:"time"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true }, //TODO: change this to only allow connections from our domain
}

//TODO: Implement authentication of some sort
//TODO: Implement a way to send data to the client
//TODO: Implement a way to receive data from the client
//TODO: keep track of clients, add to map when connected with their id, remove when disconnected
//TODO: Implement a way to send data to a specific client
//TODO keep organzied create individual functions for each task
//TODO: create list of every sort of data we need to pull from database, then create query file
//create message channel, send all messages to that channel from each client, then have a goroutine that listens to that channel and sends the messages to the client

func (c *Controller) handleWebSocketAuth(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(err.Error())
		return
	}

	//our websocket url will look like this ws://localhost:5000/authenticate?client-id=newtest&client-type=minecraft-bot&x-api-key=123456789

	client_id := r.URL.Query().Get("client-id")
	client_type := r.URL.Query().Get("client-type")
	api_key := r.URL.Query().Get("x-api-key")

	//if there is no client id or clientType or apiKey headers, then close the connection
	if client_id == "" || client_type == "" || api_key == "" {
		c.Logger.Error("Client did not provide client-id, client-type, or x-api-key headers")
		//send message to client
		errorMessage := "Missing required headers: client-id, client-type, x-api-key"
		err = conn.WriteMessage(websocket.TextMessage, []byte(errorMessage))
		if err != nil {
			c.Logger.Error(err.Error())
		}
		//close connection
		conn.Close()
		return
	}

	//Eventually put this in a function that saves it to a logfile when program exits
	c.Logger.Info(fmt.Sprintf("Client connected to websocket: %s %s %s", client_id, client_type, api_key))

	c.Mutex.Lock()
	c.Clients[client_id] = conn
	c.Mutex.Unlock()

	defer func() {
		c.Mutex.Lock()
		delete(c.Clients, client_id)
		c.Mutex.Unlock()
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			c.Logger.Error(err.Error())
			return
		}

		// Process the raw WebSocket message
		var recievedMessage WebsocketMessage
		if err := json.Unmarshal(p, &recievedMessage); err != nil {
			c.Logger.Error(err.Error())
			continue
		}
		//c.MessageChan <- recievedMessage

		//log our messages to console
		c.Logger.Info("Message recieved from client: " + fmt.Sprintf("%v", recievedMessage))

		switch recievedMessage.Type {
		case "discord":
			switch recievedMessage.Action {
			case "chat":
				//implement chat for inbound discord messages we want to umarshall the data "chat" will be the discordinbound message struct
				//var discordInboundMessage DiscordInBoundMessage
				data, ok := recievedMessage.Data.(DiscordInBoundMessage)
				if !ok {
					c.Logger.Error("Failed to assert data as DiscordInBoundMessage")
					continue
				}
				c.Logger.Info("Discord chat message recieved from client: " + fmt.Sprintf("%v", data))
				// jsonData, err := json.Marshal(data)
				// if err != nil {
				// 	c.Logger.Error(err.Error())
				// 	continue
				// }
				// if err := json.Unmarshal(jsonData, &discordInboundMessage); err != nil {
				// 	c.Logger.Error(err.Error())
				// 	continue
				// }
			}

		case "minecraft":
			switch recievedMessage.Action {
			case "savechat":

				var minecraftChatMessage MinecraftChatMessage
				dataMap, ok := recievedMessage.Data.(map[string]interface{})
				if !ok {
					c.Logger.Error("Failed to assert data as map[string]interface{}")
					continue
				}

				jsonData, err := json.Marshal(dataMap)
				if err != nil {
					c.Logger.Error(err.Error())
					sendWebsocketError(conn, "Some Properties may not be correct")
					continue
				}

				if err := json.Unmarshal(jsonData, &minecraftChatMessage); err != nil {
					c.Logger.Error(err.Error())
					sendWebsocketError(conn, "Some Properties may not be correct")
					continue
				}

				c.Logger.Info("Minecraft chat message recieved from client: " + fmt.Sprintf("%v", minecraftChatMessage))

				c.Logger.Info(minecraftChatMessage.Message)

				//implement saving message to database

				//implement sending message to discord!

			}
		}

		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	c.Logger.Error(err.Error())
		// 	return
		// }
	}

}

func sendWebsocketError(conn *websocket.Conn, errorMessage string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(errorMessage))
}
