## API Endpoints

### Retrieve User by Name and Server

- Method: GET
- Path: `/userbyname/{name}/{server}`
- Handler Function: `GetUserByNameHandler`

This endpoint retrieves a user by their name and server. It expects the name and server to be provided as path parameters. The response will contain the user information if found, or an error message if not found.

### Retrieve User by UUID and Server

- Method: GET
- Path: `/player/{uuid}/{server}`
- Handler Function: `GetUserByUUIDHandler`

This endpoint retrieves a user by their UUID and server. It expects the UUID and server to be provided as path parameters. The response will contain the user information if found, or an error message if not found.

### WebSocket Authentication

- Method: GET
- Path: `/websocket/connect`
- Handler Function: `HandleWebSocketAuth`

This endpoint handles the WebSocket authentication. It establishes a WebSocket connection and performs the necessary authentication steps.

### Retrieve Advancements

- Method: GET
- Path: `/advancements`
- Handler Function: `GetAdvancementsHandler`

This endpoint retrieves a list of advancements achieved by players.

### Retrieve Messages

- Method: GET
- Path: `/messages`
- Handler Function: `GetMessagesHandler`

This endpoint retrieves a list of messages sent by users.

### Retrieve Random Quotes

- Method: GET
- Path: `/quote`
- Handler Function: `GetRandomQuotesHandler`

This endpoint retrieves a random quote from a collection of quotes.

### Retrieve Tablist

- Method: GET
- Path: `/tablist`
- Handler Function: `GetTablistHandler`

This endpoint retrieves the list of players currently displayed in the tablist.

### Retrieve Minecraft Deaths

- Method: GET
- Path: `/deaths`
- Handler Function: `GetMinecraftDeathsHandler`

This endpoint retrieves a list of deaths recorded in the Minecraft server.

### Retrieve Minecraft Kills

- Method: GET
- Path: `/kills`
- Handler Function: `GetMinecraftKillsHandler`

This endpoint retrieves a list of kills recorded in the Minecraft server.

### Check User Online Status

- Method: GET
- Path: `/online`
- Handler Function: `GetUserOnlineCheckHandler`

This endpoint checks if a user is online. It expects the username to be provided as a query parameter. The response will indicate whether the user is online or not.

### Search User Names

- Method: GET
- Path: `/namesearch`
- Handler Function: `GetNameSearchHandler`

This endpoint searches for user names. It expects the search query to be provided as a query parameter. The response will contain a list of user names matching the search query.
