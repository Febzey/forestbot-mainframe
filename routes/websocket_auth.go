package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true }, //TODO: change this to only allow connections from our domain
}

func (c *Controller) handleWebSocketAuth(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(err.Error())
		return
	}

	client_id := r.URL.Query().Get("client-id")
	api_key := r.URL.Query().Get("x-api-key")

	if client_id == "" || api_key == "" {
		c.Logger.Error("Client did not provide client-id, or x-api-key headers")
		//send message to client
		errorMessage := "Missing required headers: client-id, x-api-key"
		err = conn.WriteMessage(websocket.TextMessage, []byte(errorMessage))
		if err != nil {
			c.Logger.Error(err.Error())
		}
		//close connection
		conn.Close()
		return
	}

	//Eventually put this in a function that saves it to a logfile when program exits
	c.Logger.Info(fmt.Sprintf("Client connected to websocket: %s %s", client_id, api_key))

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

		var recievedMessage types.WebsocketMessage

		if err := json.Unmarshal(p, &recievedMessage); err != nil {
			c.Logger.Error(err.Error())

			//send a message back to the client telling them their message structure was invalid
			errorMessage := "Invalid message structure"
			if err = conn.WriteMessage(websocket.TextMessage, []byte(errorMessage)); err != nil {
				c.Logger.Error(err.Error())
			}
			continue
		}

		c.MessageChan <- recievedMessage

	}
}

// This function is handling messages sent to our websocket server
// example: messages sent from our minecraft bots or discord live chat!
func ProcessWebsocketMessage(c *Controller) {
	for {
		message := <-c.MessageChan

		switch message.Action {

		//Handling discord chat messages and sending them to our clients
		case "inbound_discord_chat":
			var discordMessage types.DiscordMessage
			if err := mapstructure.Decode(message.Data, &discordMessage); err != nil {
				c.Logger.Error(err.Error())
				continue
			}
			c.Logger.Info("Discord chat message recieved from client: " + fmt.Sprintf("%v", discordMessage))
			c.BroadcastMessageToClients(message)
			continue

		//Handling minecraft chat messages and sending them to our clients and saving in database.
		case "inbound_minecraft_chat":
			var minecraftChatMessage types.MinecraftChatMessage
			if err := mapstructure.Decode(message.Data, &minecraftChatMessage); err != nil {
				c.Logger.Error(err.Error()) //i think its a bug
				continue
			}

			c.Logger.Info("Minecraft chat message recieved from client: " + fmt.Sprintf("%v", minecraftChatMessage))
			err := c.Database.SaveMinecraftChatMessage(minecraftChatMessage)
			if err != nil {
				c.Logger.Error(err.Error())
			}

			c.BroadcastMessageToClients(message)
			continue

		//Handling minecraft advancement messages and sending them to our clients and saving in database.
		case "minecraft_advancement":
			var minecraftAdvancementMessage types.MinecraftAdvancementMessage
			if err := mapstructure.Decode(message.Data, &minecraftAdvancementMessage); err != nil {
				c.Logger.Error(err.Error())
				continue
			}
			c.Logger.Info("Minecraft advancement message recieved from client: " + fmt.Sprintf("%v", minecraftAdvancementMessage))

			err := c.Database.SaveMinecraftAdvancementMessage(minecraftAdvancementMessage)
			if err != nil {
				c.Logger.Error(err.Error())
			}

			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player join messages and sending them to our clients and saving in database.
		case "minecraft_player_join":
			var minecraftPlayerJoinMessage types.MinecraftPlayerJoinMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerJoinMessage); err != nil {
				c.Logger.Error(err.Error())
				continue
			}
			c.Logger.Info("Minecraft player join message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerJoinMessage))

			data, err := c.Database.SavePlayerJoin(minecraftPlayerJoinMessage)
			if err != nil {
				c.Logger.Error(err.Error())
			}

			switch data.Action {
			case "new_name":
				c.SendMessageToClient(minecraftPlayerJoinMessage.Server, types.WebsocketMessage{
					Action: "new_name",
					Data:   map[string]interface{}{"username": data.Data},
				})
				continue

			case "new_user":
				c.SendMessageToClient(minecraftPlayerJoinMessage.Server, types.WebsocketMessage{
					Action: "new_user",
					Data:   map[string]interface{}{"username": data.Data},
				})
				continue

			case "none":
				continue

			default:
				c.Logger.Error("Invalid action from database")
			}

			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player leave messages and sending them to our clients and saving in database.
		case "minecraft_player_leave":
			var minecraftPlayerLeaveMessage types.MinecraftPlayerLeaveMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerLeaveMessage); err != nil {
				c.Logger.Error(err.Error())
				continue
			}
			c.Logger.Info("Minecraft player leave message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerLeaveMessage))
			c.Database.SavePlayerLeave(minecraftPlayerLeaveMessage)
			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player death messages and sending them to our clients and saving in database.
		case "minecraft_player_death":
			var minecraftPlayerDeathMessage types.MinecraftPlayerDeathMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerDeathMessage); err != nil {
				c.Logger.Error(err.Error())
				continue
			}
			c.Logger.Info("Minecraft player death message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerDeathMessage))

			err := c.Database.InsertPlayerDeathOrKill(minecraftPlayerDeathMessage)
			if err != nil {
				c.Logger.Error(err.Error())
			}
			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player list, updating each players playtime and saving in database. while updating our playerlist map
		case "send_update_player_list":
			dataMap, ok := message.Data.(map[string]interface{})
			if !ok {
				c.Logger.Error("Expected 'data' field to be a map[string]interface{}")
				continue
			}

			// Extract the "players" array from the map
			playersArray, ok := dataMap["players"].([]interface{})
			if !ok {
				c.Logger.Error("Expected 'players' field to be a []interface{}")
				continue
			}

			// Directly decode []interface{} into []types.Player
			var minecraftPlayerListArray []types.Player
			if err := mapstructure.Decode(playersArray, &minecraftPlayerListArray); err != nil {
				c.Logger.Error(err.Error())
				continue
			}

			//we want to run a function that will update player playtime by 60000 milliseconds which is 1 minute
			for _, player := range minecraftPlayerListArray {
				err := c.Database.UpdatePlayerPlaytime(player.Uuid, player.Server)
				if err != nil {
					c.Logger.Error(err.Error())
				}
			}

			c.Mutex.Lock()
			c.PlayerLists[minecraftPlayerListArray[0].Server] = minecraftPlayerListArray
			c.Mutex.Unlock()

		default:
			{
				c.Logger.Error("Invalid message action")
				continue

			}

		}

	}
}

// This is a function that allows us to send client specific messages.
// As in we can send a message to a specific client by their client_id
func (c *Controller) SendMessageToClient(client_id string, message types.WebsocketMessage) {
	c.Mutex.Lock()
	conn, ok := c.Clients[client_id]
	c.Mutex.Unlock()
	if !ok {
		c.Logger.Error("Client with client_id: " + client_id + " does not exist")
	}

	if err := conn.WriteJSON(message); err != nil {
		c.Logger.Error(err.Error())
	}

}

// This function is used to broadcast a message
// to every single client connected to our server
func (c *Controller) BroadcastMessageToClients(message types.WebsocketMessage) {
	for _, conn := range c.Clients {
		if err := conn.WriteJSON(message); err != nil {
			c.Logger.Error(err.Error())
		}
	}
}

func SendWebsocketError(conn *websocket.Conn, errorMessage string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(errorMessage))
}
