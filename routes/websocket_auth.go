package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/febzey/ForestBot-Mainframe/utils"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true }, //TODO: change this to only allow connections from our domain
}

func InitializeWebsocketClient(conn *websocket.Conn, api_key string, mc_server string) *WebsocketClient {

	client_id, err := utils.RandomUUID()
	if err != nil {
		fmt.Println(err.Error())
		sendMessageByStructure(conn, types.WebsocketMessage{
			Client_id: client_id,
			Action:    "error",
			Data:      "Error generating client_id - Internal server error",
		})
		conn.Close()
		return nil
	}

	if mc_server == "" || api_key == "" {
		errorMessage := "Missing required headers: client-id, x-api-key"
		sendMessageByStructure(conn, types.WebsocketMessage{
			Client_id: client_id,
			Action:    "error",
			Data:      errorMessage,
		})
		//close connection
		conn.Close()
	}

	sendMessageByStructure(conn, types.WebsocketMessage{
		Client_id: client_id,
		Action:    "id",
		Data:      client_id,
	})

	client_permissions := utils.CheckApiKey(api_key)

	return &WebsocketClient{
		ClientID:    client_id,
		Conn:        conn,
		Permissions: client_permissions,
		Mc_server:   mc_server,
	}

}

func (c *Controller) handleWebSocketAuth(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(err.Error())
		return
	}

	mc_server := r.URL.Query().Get("server")
	api_key := r.URL.Query().Get("x-api-key")

	client := InitializeWebsocketClient(conn, api_key, mc_server)

	//Eventually put this in a function that saves it to a logfile when program exits
	c.Logger.Info(fmt.Sprintf("Client connected to websocket: %s %s", client.ClientID, api_key))

	c.Mutex.Lock()
	c.Clients[client.ClientID] = client
	c.Mutex.Unlock()

	defer func() {
		c.Mutex.Lock()
		delete(c.Clients, client.ClientID)
		c.Mutex.Unlock()
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			c.Logger.Error(err.Error())
			return
		}

		//check if client has permissions to write data
		if !client.Permissions.Write {
			c.Logger.Error("Client does not have permission to write data")
			//send message to client
			errorMessage := "Client does not have permission to write data"
			sendMessageByStructure(conn, types.WebsocketMessage{
				Client_id: client.ClientID,
				Action:    "error",
				Data:      errorMessage,
			})
			//close connection
			conn.Close()
			return
		}

		var recievedMessage types.WebsocketMessage

		if err := json.Unmarshal(p, &recievedMessage); err != nil {
			c.Logger.Error(err.Error())

			//send a message back to the client telling them their message structure was invalid
			errorMessage := "Invalid message structure"
			sendMessageByStructure(conn, types.WebsocketMessage{
				Client_id: client.ClientID,
				Action:    "error",
				Data:      errorMessage,
			})
			continue
		}

		c.MessageChan <- MessageChannel{
			ClientID: client.ClientID,
			Message:  recievedMessage,
		}

	}
}

// This function is handling messages sent to our websocket server
// example: messages sent from our minecraft bots or discord live chat!
func ProcessWebsocketMessage(c *Controller) {
	for {
		messageChannel := <-c.MessageChan

		message := messageChannel.Message
		realClientID := messageChannel.ClientID

		if _, ok := c.Clients[message.Client_id]; !ok {

			c.Logger.Error("Client with client_id: " + message.Client_id + " does not exist")

			sendMessageByStructure(c.Clients[realClientID].Conn, types.WebsocketMessage{
				Client_id: realClientID,
				Action:    "error",
				Data:      "Client with client_id: " + message.Client_id + " does not exist",
			})

			continue

		}

		switch message.Action {

		//Handling discord chat messages and sending them to our clients
		case "inbound_discord_chat":
			var discordMessage types.DiscordMessage
			if err := mapstructure.Decode(message.Data, &discordMessage); err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for inbound_discord_chat",
				})
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
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for inbound_minecraft_chat",
				})
				continue
			}

			c.Logger.Info("Minecraft chat message recieved from client: " + fmt.Sprintf("%v", minecraftChatMessage))
			err := c.Database.SaveMinecraftChatMessage(minecraftChatMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft chat message to database",
				})
			}

			c.BroadcastMessageToClients(message)
			continue

		//Handling minecraft advancement messages and sending them to our clients and saving in database.
		case "minecraft_advancement":
			var minecraftAdvancementMessage types.MinecraftAdvancementMessage
			if err := mapstructure.Decode(message.Data, &minecraftAdvancementMessage); err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_advancement",
				})
				continue
			}
			c.Logger.Info("Minecraft advancement message recieved from client: " + fmt.Sprintf("%v", minecraftAdvancementMessage))

			err := c.Database.SaveMinecraftAdvancementMessage(minecraftAdvancementMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft advancement message to database",
				})
			}

			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player join messages and sending them to our clients and saving in database.
		case "minecraft_player_join":
			var minecraftPlayerJoinMessage types.MinecraftPlayerJoinMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerJoinMessage); err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_player_join",
				})
				continue
			}
			c.Logger.Info("Minecraft player join message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerJoinMessage))

			data, err := c.Database.SavePlayerJoin(minecraftPlayerJoinMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft player join message to database",
				})
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
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_player_leave",
				})

				continue
			}
			c.Logger.Info("Minecraft player leave message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerLeaveMessage))
			err := c.Database.SavePlayerLeave(minecraftPlayerLeaveMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft player leave message to database",
				})
			}
			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player death messages and sending them to our clients and saving in database.
		case "minecraft_player_death":
			var minecraftPlayerDeathMessage types.MinecraftPlayerDeathMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerDeathMessage); err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_player_death",
				})
				//we need a way to send this to the client that sent this message
				//we can do this by sending the client_id in the message

				continue
			}
			c.Logger.Info("Minecraft player death message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerDeathMessage))

			err := c.Database.InsertPlayerDeathOrKill(minecraftPlayerDeathMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft player death message to database",
				})
			}
			c.BroadcastMessageToClients(message)

			continue

		//Handling minecraft player list, updating each players playtime and saving in database. while updating our playerlist map
		case "send_update_player_list":
			dataMap, ok := message.Data.(map[string]interface{})
			if !ok {
				c.Logger.Error("Expected 'data' field to be a map[string]interface{}")
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for send_update_player_list",
				})
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
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for send_update_player_list",
				})
				continue
			}

			//we want to run a function that will update player playtime by 60000 milliseconds which is 1 minute
			for _, player := range minecraftPlayerListArray {
				err := c.Database.UpdatePlayerPlaytime(player.Uuid, player.Server)
				if err != nil {
					c.Logger.Error(err.Error())
					sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
						Client_id: message.Client_id,
						Action:    "error",
						Data:      "Error updating player playtime in database",
					})
				}
			}

			c.Mutex.Lock()
			c.PlayerLists[minecraftPlayerListArray[0].Server] = minecraftPlayerListArray
			c.Mutex.Unlock()

		default:
			{
				c.Logger.Error("Invalid message action")
				sendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message action - not a valid action",
				})
				continue

			}

		}

	}
}

// This is a function that allows us to send client specific messages.
// As in we can send a message to a specific client by their client_id
func (c *Controller) SendMessageToClient(client_id string, message types.WebsocketMessage) {
	c.Mutex.Lock()
	client, ok := c.Clients[client_id]
	c.Mutex.Unlock()
	if !ok {
		c.Logger.Error("Client with client_id: " + client_id + " does not exist")
	}

	if err := client.Conn.WriteJSON(message); err != nil {
		c.Logger.Error(err.Error())
	}

}

// This function is used to broadcast a message
// to every single client connected to our server
func (c *Controller) BroadcastMessageToClients(message types.WebsocketMessage) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, client := range c.Clients {
		if client.Permissions.Read {
			if err := client.Conn.WriteJSON(message); err != nil {
				c.Logger.Error(err.Error())
			}
		}
	}
}

func SendWebsocketError(conn *websocket.Conn, errorMessage string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(errorMessage))
}

func sendMessageByStructure(conn *websocket.Conn, message types.WebsocketMessage) error {
	err := conn.WriteJSON(message)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
