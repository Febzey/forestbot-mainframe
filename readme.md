# API Documentation

## Overview

Welcome to the ForestBot API documentation. ForestBot is a comprehensive system that facilitates real-time data exchange and interaction with Minecraft servers. The API provides endpoints to retrieve information about players, their achievements, playtime, and enables WebSocket connections for dynamic, real-time updates.

## WebSocket Integration

ForestBot's WebSocket integration empowers developers to establish real-time communication channels between their applications and Minecraft servers. This feature is particularly useful for applications requiring dynamic updates such as live chats, player events, and more.

For detailed information on how to integrate and utilize ForestBot's WebSocket functionality, please refer to the [WebSocket Integration Guide](/controllers/websocket.md).


## API keys and Authentication
Read here for documenation on authentication and obtaining/using keys.
[Authentication and Keys Guide](/keyservice/auth.md)

## HTTP Endpoints


### GET Requests

### Get User by Name
- **Endpoint:** `/api/v1/playername`
- **Description:** Gets a user by their name
- **Example URL:** `http://localhost:5000/api/v1/playername?name=febzey&server=simplyvanilla`
- **Queries:** 
  - `name`: The username of the player
  - `server`: The Minecraft server name

### Get User by UUID
- **Endpoint:** `/api/v1/playeruuid`
- **Description:** Gets a user by their UUID
- **Example URL:** `http://localhost:5000/api/v1/playeruuid?uuid=30303-addwdwd-222=3333&server=simplyvanilla`
- **Queries:** 
  - `uuid`: The UUID of the player
  - `server`: The Minecraft server name

### WebSocket Connect
- **Endpoint:** `/api/v1/websocket/connect`
- **Description:** WebSocket for real-time data exchange between server and client (playtime, chat, etc.)

### Get Advancements
- **Endpoint:** `/api/v1/advancements`
- **Description:** Gets the advancements of a player
- **Example URL:** `http://localhost:5000/api/v1/advancements?uuid=30303-addwdwd-222=3333&server=simplyvanilla&limit=100&order=DESC`
- **Queries:** 
  - `uuid`: The UUID of the player
  - `server`: The Minecraft server name
  - `limit`: (Optional) Limit the number of results
  - `order`: (Optional) Order of results (ASC or DESC)

### Get Messages
- **Endpoint:** `/api/v1/messages`
- **Description:** Gets the messages for a player
- **Example URL:** `http://localhost:5000/api/v1/messages?name=febzey&server=simplyvanilla&limit=100&order=DESC`
- **Queries:** 
  - `name`: The username of the player
  - `server`: The Minecraft server name
  - `limit`: (Optional) Limit the number of results
  - `order`: (Optional) Order of results (ASC or DESC)

### Get Random Quote
- **Endpoint:** `/api/v1/quote`
- **Description:** Get a random quote from a user on a server
- **Example URL:** `http://localhost:5000/api/v1/quote?name=febzey&server=simplyvanilla`
- **Queries:** 
  - `name`: The username of the player
  - `server`: The Minecraft server name

### Get Tablist
- **Endpoint:** `/api/v1/tablist`
- **Description:** Get the tablist for a server
- **Example URL:** `http://localhost:5000/api/v1/tablist?server=simplyvanilla`
- **Queries:** 
  - `server`: The Minecraft server name

### Get Deaths
- **Endpoint:** `/api/v1/deaths`
- **Description:** Gets the deaths of a player
- **Example URL:** `http://localhost:5000/api/v1/deaths?uuid=30303-addwdwd-222=3333&server=simplyvanilla&limit=100&order=DESC`
- **Queries:** 
  - `uuid`: The UUID of the player
  - `server`: The Minecraft server name
  - `limit`: (Optional) Limit the number of results
  - `order`: (Optional) Order of results (ASC or DESC)

### Get Kills
- **Endpoint:** `/api/v1/kills`
- **Description:** Gets the kills of a player
- **Example URL:** `http://localhost:5000/api/v1/kills?uuid=30303-addwdwd-222=3333&server=simplyvanilla&limit=100&order=DESC`
- **Queries:** 
  - `uuid`: The UUID of the player
  - `server`: The Minecraft server name
  - `limit`: (Optional) Limit the number of results
  - `order`: (Optional) Order of results (ASC or DESC)

### Get User Online Check
- **Endpoint:** `/api/v1/online`
- **Description:** Checks if a user is online and returns server and true or false
- **Queries:** 
  - `username`: The username of the player

### Get Name Search
- **Endpoint:** `/api/v1/namesearch`
- **Description:** Returns back 6 closest names to the username
- **Queries:** 
  - `username`: The username of the player
  - `server`: (Optional) The Minecraft server name

### Get WhoIs
- **Endpoint:** `/api/v1/whois`
- **Description:** Returns back the whois data for the username
- **Example URL:** `http://localhost:5000/api/v1/whois?username=febzey`
- **Queries:** 
  - `username`: The username of the player

### Get Discord Guilds
- **Endpoint:** `/api/v1/discord/guilds`
- **Description:** Get all the guilds the Discord bot is in

### Get Discord Live Chat Channels
- **Endpoint:** `/api/v1/discord/livechats`
- **Description:** Get all live chat channels for the Discord bot


## POST Requests

### Add Discord Guild
- **Endpoint:** `/api/v1/discord/addguild`
- **Description:** Adds a guild to the database
- **Example URL:** `http://localhost:5000/api/v1/discord/addguild`
- **Method:** `POST`
- **Handler Function:** `controller.PostDiscordGuild`

### Add Discord Live Chat
- **Endpoint:** `/api/v1/discord/addlivechat`
- **Description:** Adds a live chat channel to the database
- **Example URL:** `http://localhost:5000/api/v1/discord/addlivechat`
- **Method:** `POST`
- **Handler Function:** `controller.PostDiscordLiveChat`

## DELETE Requests

### Delete Discord Guild
- **Endpoint:** `/api/v1/discord/deleteguild`
- **Description:** Deletes a guild from the database
- **Example URL:** `http://localhost:5000/api/v1/discord/deleteguild?guild_id=123`
- **Method:** `DELETE`
- **Handler Function:** `controller.DeleteDiscordGuild`

### Delete Discord Live Chat
- **Endpoint:** `/api/v1/discord/deletelivechat`
- **Description:** Deletes a live chat channel from the database
- **Example URL:** `http://localhost:5000/api/v1/discord/deletelivechat?channel_id=123`
- **Method:** `DELETE`
- **Handler Function:** `controller.DeleteDiscordLiveChat`
