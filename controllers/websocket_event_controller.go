package controllers

import (
	"fmt"

	"github.com/febzey/ForestBot-Mainframe/types"
	"github.com/mitchellh/mapstructure"
)

type WebsocketEvent struct {
	//Client id given by the user, this should match the uuid generated
	//at the start of their session
	Client_id string `json:"client_id"`

	//Action for the message example:
	//minecraft_chat_message, minecraft_player_join, minecraft_player_leave, minecraft_player_kill, minecraft_player_death, discord_message
	Action string `json:"action"`

	//The data for the message.
	//The data can be many of our structs below.
	Data interface{} `json:"data"`
}

//Individual Handler structure.
//I think this struct is only used when we need to create the initial map.
type Handler struct {
	action  string
	handler func(*Controller, WebsocketEvent)
}

/*
Load our events
*/
func (c *Controller) setupWebsocketEventHandlers() {
	events := []Handler{
		{
			action:  "inbound_discord_chat",
			handler: handleInboundDiscordChat,
		},
		{
			action:  "inbound_minecraft_chat",
			handler: handleInboundMinecraftChat,
		},
		{
			action:  "minecraft_advancement",
			handler: handleMinecraftAdvancement,
		},
		{
			action:  "minecraft_player_join",
			handler: handleMinecraftPlayerJoin,
		},
		{
			action:  "minecraft_player_leave",
			handler: handleMinecraftPlayerLeave,
		},
		{
			action:  "minecraft_player_death",
			handler: handleMinecraftPlayerDeath,
		},
		{
			action:  "send_update_player_list",
			handler: handleUpdatePlayerList,
		},
	}

	for _, event := range events {
		c.Handlers[event.action] = event
	}
}

/*
Handling all inbound websocket messages.
The client must send their client_id with each message.
This function is ran as a go routine and will be running continously and is called inside of the controller file.
This function can be running even before any clients connect.
*/
func ProcessWebsocketEvent(c *Controller) {
	for {
		messageChannel := <-c.MessageChan

		message, realClientID := messageChannel.Message, messageChannel.ClientID

		if _, ok := c.Clients[message.Client_id]; !ok {
			c.sendMessageByStructure(realClientID, WebsocketEvent{
				Client_id: realClientID,
				Action:    "error",
				Data:      "The client_id you gave is not valid. or unexpected error.",
			})
			continue
		}

		if message.Client_id != realClientID {
			c.sendMessageByStructure(realClientID, WebsocketEvent{
				Client_id: realClientID,
				Action:    "error",
				Data:      "It seems you sent a client_id that does not match the one assigned to you!",
			})
		}

		event, ok := c.Handlers[message.Action]
		if !ok {
			c.sendMessageByStructure(realClientID, WebsocketEvent{
				Client_id: realClientID,
				Action:    "error",
				Data:      "Invalid event action type.",
			})
		}

		if event.handler != nil {
			event.handler(c, message)
		}
	}
}

/**
*
*
* Handling all websocket messages mainly for our mc-client,
* we often broadcast the message back to each client for them to use.
*
*
**/

/**
*Handling Inbound discord chat messages from our websocket
**/
func handleInboundDiscordChat(c *Controller, message WebsocketEvent) {
	var discordMessage types.DiscordMessage
	if err := mapstructure.Decode(message.Data, &discordMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "Invalid message structure for inbound_discord_chat")
		return
	}

	c.Logger.WebsocketInfo("Discord chat message received from client: " + fmt.Sprintf("%v", discordMessage))
	c.BroadcastMessageToClients(message)
}

/*
* Handling minecraft inbound chat messages from our websocket.
 */
func handleInboundMinecraftChat(c *Controller, message WebsocketEvent) {
	var minecraftChatMessage types.MinecraftChatMessage
	if err := mapstructure.Decode(message.Data, &minecraftChatMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "Invalid message structure for inbound_minecraft_chat")
		return
	}

	c.Logger.WebsocketInfo("Minecraft chat message received from client: " + fmt.Sprintf("%v", minecraftChatMessage))
	err := c.Database.SaveMinecraftChatMessage(minecraftChatMessage)
	if err != nil {
		c.sendErrorMessage(message.Client_id, "Error saving minecraft chat message to database")
	}

	c.BroadcastMessageToClients(message)
}

/*
* Handling inbound minecraft chat advancements from our websocket.
 */
func handleMinecraftAdvancement(c *Controller, message WebsocketEvent) {
	var minecraftAdvancementMessage types.MinecraftAdvancementMessage
	if err := mapstructure.Decode(message.Data, &minecraftAdvancementMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "Invalid message structure for minecraft_advancement")
		return
	}

	c.Logger.WebsocketInfo("Minecraft advancement message received from client: " + fmt.Sprintf("%v", minecraftAdvancementMessage))

	err := c.Database.SaveMinecraftAdvancementMessage(minecraftAdvancementMessage)
	if err != nil {
		c.sendErrorMessage(message.Client_id, "Error saving minecraft advancement message to database")
	}

	c.BroadcastMessageToClients(message)
}

/*
* Handing minecraft player join messages from our websocket.
 */
func handleMinecraftPlayerJoin(c *Controller, message WebsocketEvent) {
	var minecraftPlayerJoinMessage types.MinecraftPlayerJoinMessage
	if err := mapstructure.Decode(message.Data, &minecraftPlayerJoinMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "minecraft_player_join")
		return
	}

	c.Logger.WebsocketInfo("Minecraft player join message received from client: " + fmt.Sprintf("%v", minecraftPlayerJoinMessage))

	data, err := c.Database.SavePlayerJoin(minecraftPlayerJoinMessage)
	if err != nil {
		c.sendErrorMessage(message.Client_id, "Error saving minecraft player join message to database")
		return
	}

	player := types.Player{
		Username: minecraftPlayerJoinMessage.Username,
		Uuid:     minecraftPlayerJoinMessage.Uuid,
		Latency:  minecraftPlayerJoinMessage.Latency,
		Server:   minecraftPlayerJoinMessage.Server,
		Head_url: head_url + minecraftPlayerJoinMessage.Username + "/16",
	}

	c.addUserToPlayerList(minecraftPlayerJoinMessage.Server, player)

	switch data.Action {
	case "new_name":
		c.BroadcastMessageToClients(WebsocketEvent{
			Client_id: message.Client_id,
			Action:    "new_name",
			Data:      map[string]interface{}{"user": data.Data, "server": c.Clients[message.Client_id].Mc_server},
		})
	case "new_user":
		c.BroadcastMessageToClients(WebsocketEvent{
			Client_id: message.Client_id,
			Action:    "new_user",
			Data:      map[string]interface{}{"user": data.Data, "server": c.Clients[message.Client_id].Mc_server},
		})
	case "none":
		c.BroadcastMessageToClients(message)
	default:
		c.Logger.Error("Invalid action from database")
	}
}

/*
* Handing minecraft player leave messages from our websocket
 */
func handleMinecraftPlayerLeave(c *Controller, message WebsocketEvent) {
	var minecraftPlayerLeaveMessage types.MinecraftPlayerLeaveMessage
	if err := mapstructure.Decode(message.Data, &minecraftPlayerLeaveMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "minecraft_player_leave")
		return
	}

	c.Logger.WebsocketInfo("Minecraft player leave message received from client: " + fmt.Sprintf("%v", minecraftPlayerLeaveMessage))

	if err := c.Database.SavePlayerLeave(minecraftPlayerLeaveMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "Error saving minecraft player leave message to database")
		return
	}

	c.removeUserFromPlayerList(minecraftPlayerLeaveMessage.Server, minecraftPlayerLeaveMessage.Username)

	c.BroadcastMessageToClients(message)
}

/*
* Handling minecraft player deaths and kills
 */
func handleMinecraftPlayerDeath(c *Controller, message WebsocketEvent) {
	var minecraftPlayerDeathMessage types.MinecraftPlayerDeathMessage
	if err := mapstructure.Decode(message.Data, &minecraftPlayerDeathMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "minecraft_player_death")
		return
	}

	c.Logger.WebsocketInfo("Minecraft player death message received from client: " + fmt.Sprintf("%v", minecraftPlayerDeathMessage))

	if err := c.Database.InsertPlayerDeathOrKill(minecraftPlayerDeathMessage); err != nil {
		c.sendErrorMessage(message.Client_id, "Error saving minecraft player death message to database")
		return
	}

	c.BroadcastMessageToClients(message)
}

/*
* Handling minecraft server lists, and users playtime update
 */
func handleUpdatePlayerList(c *Controller, message WebsocketEvent) {
	dataMap, ok := message.Data.(map[string]interface{})
	if !ok {
		c.sendErrorMessage(message.Client_id, "send_update_player_list")
		return
	}

	// Extract the "players" array from the map
	playersArray, ok := dataMap["players"].([]interface{})
	if !ok {
		c.Logger.Error("Expected 'players' field to be a []interface{}")
		return
	}

	// Directly decode []interface{} into []types.Player
	var minecraftPlayerListArray []types.Player

	if err := mapstructure.Decode(playersArray, &minecraftPlayerListArray); err != nil {
		c.sendErrorMessage(message.Client_id, "send_update_player_list")
		return
	}

	// Update player playtime and add to player list
	for _, player := range minecraftPlayerListArray {
		if err := c.Database.UpdatePlayerPlaytime(player.Uuid, player.Server); err != nil {
			c.sendErrorMessage(message.Client_id, "Error updating player playtime in database")
			continue
		}

		player.Head_url = head_url + player.Username + "/16"

		c.addUserToPlayerList(player.Server, player)
	}
}
