/******

	This is the websocket controller.
	Handling everything websocket.

******/

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

func (c *Controller) InitializeWebsocketClient(conn *websocket.Conn, api_key string, mc_server string) *WebsocketClient {

	client_id, err := utils.RandomUUID()
	if err != nil {
		fmt.Println(err.Error())
		c.SendMessageByStructure(conn, types.WebsocketMessage{
			Client_id: client_id,
			Action:    "error",
			Data:      "Error generating client_id - Internal server error",
		})
		conn.Close()
		return nil
	}

	if mc_server == "" || api_key == "" {
		errorMessage := "Missing required headers: client-id, x-api-key - closing the connection"
		c.SendMessageByStructure(conn, types.WebsocketMessage{
			Client_id: client_id,
			Action:    "error",
			Data:      errorMessage,
		})
		//close connection
		conn.Close()
		return nil
	}

	c.SendMessageByStructure(conn, types.WebsocketMessage{
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

func (c *Controller) handleUserDisconnect(client_id string) {
	c.Mutex.Lock()
	delete(c.Clients, client_id)
	c.Mutex.Unlock()
}

func (c *Controller) handleWebSocketAuth(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(err.Error())
		return
	}

	mc_server := r.URL.Query().Get("server")
	api_key := r.URL.Query().Get("x-api-key")

	client := c.InitializeWebsocketClient(conn, api_key, mc_server)

	c.Logger.WebsocketConnect((fmt.Sprintf("Websocket Client Connected For Minecraft Server: %s | ID: %s ", mc_server, client.ClientID)))

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
			if closeErr, ok := err.(*websocket.CloseError); ok {
				c.Logger.WebsocketDisconnect(fmt.Sprintf("Client: %s, Code: %d, Text: %s", client.ClientID, closeErr.Code, closeErr.Text))
			} else {
				c.Logger.Error(err.Error())
			}

			c.handleUserDisconnect(client.ClientID)
			break
		}
		// if err != nil {
		// 	c.Logger.Error(err.Error())

		// 	return
		// }

		if !client.Permissions.Write {
			c.Logger.Error("Client does not have permission to write data")
			//send message to client
			errorMessage := "Client does not have permission to write data"
			c.SendMessageByStructure(conn, types.WebsocketMessage{
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
			c.SendMessageByStructure(conn, types.WebsocketMessage{
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

/*
*
Handling all inbound websocket messages.
The client must send their client_id with each message.
*
*/
func ProcessWebsocketMessage(c *Controller) {
	for {
		messageChannel := <-c.MessageChan

		message := messageChannel.Message
		realClientID := messageChannel.ClientID

		if _, ok := c.Clients[message.Client_id]; !ok {
			c.Logger.WebsocketError("Client with client_id: " + message.Client_id + " does not exist")
			c.SendMessageByStructure(c.Clients[realClientID].Conn, types.WebsocketMessage{
				Client_id: realClientID,
				Action:    "error",
				Data:      "Client with client_id: " + message.Client_id + " does not exist",
			})
			continue
		}

		switch message.Action {

		/************************


		Handling all inbound Discord Chat Messages
		Saving to database and then broadcasting the message back
		to each client incase they want the message also.


		************************/
		case "inbound_discord_chat":
			var discordMessage types.DiscordMessage
			if err := mapstructure.Decode(message.Data, &discordMessage); err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for inbound_discord_chat",
				})
				continue
			}
			c.Logger.WebsocketInfo("Discord chat message recieved from client: " + fmt.Sprintf("%v", discordMessage))
			c.BroadcastMessageToClients(message)
			continue
		/************************


		Handling all inbound Minecraft Chat Messages
		Saving to database and then broadcasting the message back
		to each client incase they want the message also.


		************************/
		case "inbound_minecraft_chat":
			var minecraftChatMessage types.MinecraftChatMessage
			if err := mapstructure.Decode(message.Data, &minecraftChatMessage); err != nil {
				c.Logger.Error(err.Error()) //i think its a bug
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for inbound_minecraft_chat",
				})
				continue
			}

			c.Logger.WebsocketInfo("Minecraft chat message recieved from client: " + fmt.Sprintf("%v", minecraftChatMessage))
			err := c.Database.SaveMinecraftChatMessage(minecraftChatMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft chat message to database",
				})
			}

			c.BroadcastMessageToClients(message)
			continue

		/************************


		Handling all inbound Minecraft Advancement Messages,
		we save to database and then broadcast.


		************************/
		case "minecraft_advancement":
			var minecraftAdvancementMessage types.MinecraftAdvancementMessage
			if err := mapstructure.Decode(message.Data, &minecraftAdvancementMessage); err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_advancement",
				})
				continue
			}
			c.Logger.WebsocketInfo("Minecraft advancement message recieved from client: " + fmt.Sprintf("%v", minecraftAdvancementMessage))

			err := c.Database.SaveMinecraftAdvancementMessage(minecraftAdvancementMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft advancement message to database",
				})
			}

			c.BroadcastMessageToClients(message)

			continue

		/************************


		Handling all inbound Minecraft Player Join Messages,
		we save to database, update their lastseen, join count,
		we also check to see if the user has changed their name,
		we also check if the user has never been seen before,
		we send appropriate message actions if those cases are found.
		we then broadcast the join message to all clients.


		************************/
		case "minecraft_player_join":
			var minecraftPlayerJoinMessage types.MinecraftPlayerJoinMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerJoinMessage); err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_player_join",
				})
				continue
			}
			c.Logger.WebsocketInfo("Minecraft player join message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerJoinMessage))

			data, err := c.Database.SavePlayerJoin(minecraftPlayerJoinMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft player join message to database",
				})
			}

			c.addUserToPlayerList(minecraftPlayerJoinMessage.Server, types.Player{
				Username: minecraftPlayerJoinMessage.Username,
				Uuid:     minecraftPlayerJoinMessage.Uuid,
				Latency:  minecraftPlayerJoinMessage.Latency,
				Server:   minecraftPlayerJoinMessage.Server,
				Head_url: head_url + minecraftPlayerJoinMessage.Username + "/16",
			})

			switch data.Action {
			case "new_name":
				c.BroadcastMessageToClients(types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "new_name",
					Data:      map[string]interface{}{"user": data.Data, "server": c.Clients[message.Client_id].Mc_server},
				})
				continue

			case "new_user":
				c.BroadcastMessageToClients(types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "new_user",
					Data:      map[string]interface{}{"user": data.Data, "server": c.Clients[message.Client_id].Mc_server},
				})
				continue

			case "none":
				c.BroadcastMessageToClients(message)
				continue

			default:
				c.Logger.Error("Invalid action from database")
			}

			continue

		/************************


		Handling inbound player leave messages,
		we update lastseen and leave count in database,
		here we also will remove the player from the servers player list
		that they just left from.
		we then broadcast the message to all clients.


		************************/
		case "minecraft_player_leave":
			var minecraftPlayerLeaveMessage types.MinecraftPlayerLeaveMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerLeaveMessage); err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_player_leave",
				})

				continue
			}
			c.Logger.WebsocketInfo("Minecraft player leave message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerLeaveMessage))
			err := c.Database.SavePlayerLeave(minecraftPlayerLeaveMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft player leave message to database",
				})
			}

			c.removeUserFromPlayerList(minecraftPlayerLeaveMessage.Server, minecraftPlayerLeaveMessage.Username)

			c.BroadcastMessageToClients(message)

			continue

		/************************


		Handling player death messages, we save to database,
		either PVP or PVE death messages.
		we then broadcast back to all clients.


		************************/
		case "minecraft_player_death":
			var minecraftPlayerDeathMessage types.MinecraftPlayerDeathMessage
			if err := mapstructure.Decode(message.Data, &minecraftPlayerDeathMessage); err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message structure for minecraft_player_death",
				})
				//we need a way to send this to the client that sent this message
				//we can do this by sending the client_id in the message

				continue
			}
			c.Logger.WebsocketInfo("Minecraft player death message recieved from client: " + fmt.Sprintf("%v", minecraftPlayerDeathMessage))

			err := c.Database.InsertPlayerDeathOrKill(minecraftPlayerDeathMessage)
			if err != nil {
				c.Logger.Error(err.Error())
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Error saving minecraft player death message to database",
				})
			}
			c.BroadcastMessageToClients(message)

			continue

		/************************


		This action allows our minecraft bots to send a
		updated player list to our websocket every 60000 milliseconds (60 seconds),
		we will update the c.PlayerLists with the server that sent this action,
		we also will run a function that updates all users playtime by 60000 seconds.
		For this action to work properly it is up to the bot owner to send this action every 60 seconds.


		************************/
		case "send_update_player_list":
			dataMap, ok := message.Data.(map[string]interface{})
			if !ok {
				c.Logger.Error("Expected 'data' field to be a map[string]interface{}")
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
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
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
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
					c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
						Client_id: message.Client_id,
						Action:    "error",
						Data:      "Error updating player playtime in database",
					})
				}

				player.Head_url = head_url + player.Username + "/16"

				c.addUserToPlayerList(player.Server, player)
			}

			continue

		/************************


		We will send this back to the client if they sent
		an unknown message action.


		************************/
		default:
			{
				c.SendMessageByStructure(c.Clients[message.Client_id].Conn, types.WebsocketMessage{
					Client_id: message.Client_id,
					Action:    "error",
					Data:      "Invalid message action - not a valid action",
				})
				continue

			}

		}

	}
}

/*
With this function we are able to send a message
to every websocket client connected to our server.
*/
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

/*
This function allows us to send a message to specific
websocket connection, while following our websocket message structure.
*/
func (c *Controller) SendMessageByStructure(conn *websocket.Conn, message types.WebsocketMessage) error {
	err := conn.WriteJSON(message)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
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
