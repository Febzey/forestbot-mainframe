# ForestBot WebSocket Documentation

## WebSocket Connection

When a user connects to the ForestBot WebSocket, they can provide the following parameters as query parameters in the URL:

- **server:** The Minecraft server identifier. (optional)
- **is-bot-client:** Boolean indicating if the client is a bot client. (optional)

### Example URL

ws://forestbot-server.com/api/v1/websocket/connect?server=simplyvanilla&x-api-key=your-api-key&is-bot-client=true


## Bot Client Considerations

- If `is-bot-client` is set to true, the `server` parameter is mandatory.
- Only one bot client (`is-bot-client="true"`) is allowed per Minecraft server to prevent redundancy in data gathering.
- Bot clients, which act as Minecraft bots for data gathering, use read-write API keys.

## API keys and Authentication
Read here for documenation on authentication and obtaining/using keys.
[Authentication and Keys Guide](/keyservice/auth.md)


## Self-Hosting Bot Client

Users can self-host the bot client while connecting to the main server and using the central database. Contact the administrator to generate a read-write API key for self-hosting.

## WebSocket Event Structure

All incoming and outbound WebSocket events (messages) follow this structure:

```go
type WebsocketEvent struct {
    // Client id given by the user, matching the UUID generated at the start of their session.
    Client_id string `json:"client_id"`

    // Action for the event
    Action string `json:"action"`

    // The data for the message.
    Data interface{} `json:"data"`
}
```
This structure will be followed when sending or recieving messages

# Central Processor

All inbound messages within the ForestBot WebSocket are directed to a central processor. This processor is responsible for handling and distributing messages, implementing a centralized approach to streamline event processing. This ensures uniformity in message handling across the entire WebSocket service.

## Supported Actions

The `Action` field within the `WebsocketEvent` structure defines the type of message being processed. Examples of supported actions include:

- `inbound_minecraft_chat` (directional)
- `minecraft_player_join` (directional)
- `minecraft_player_leave` (directional)
- `minecraft_player_kill` (directional)
- `minecraft_player_death` (directional)
- `inbound_discord_message` (directional)
- `error` (outbound)
- `new_name` (outbound)
- `new_user`(outbound)
- `key-accepted` (outbound)

Each action corresponds to specific data structures, enabling seamless integration and processing of diverse events.


## Example Use Cases

### Regular Client Connection

- A regular client connecting to the server only needs to provide the `x-api-key` in the URL.
- Live chat bridges
- Discord Bot
- Personal data project
- Filtering for specific server data can be done on the client side.

### Minecraft Bot Client

- Write key required
- Register a bot client for your minecraft server
- Self host ForestBot and start saving data on your server
