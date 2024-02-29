# Authentication And Authorization

Here you can learn how the server authenticates you and how to obtain your own key.
<br>
Keys are generated manually by contacting an admin.

## Generating API Keys

ForestBot provides two types of API keys, each serving different purposes:

- **Read/Write Keys (For Bot Clients):**
  - Intended for bot clients interacting with Minecraft servers.
  - Enables both reading and writing operations.
  - Contact an administrator of the project to obtain a Read/Write key.
  - You'll need to provide your intentions for connecting to the server and a valid email for contact.
  - Once generated, you own the key and are responsible for its secure management.

- **Read-Only Keys (For Regular Clients):**
  - Designed for regular clients interested in reading data from the server.
  - Allows read-only operations.
  - To acquire a Read-Only key, reach out to an administrator of the project.
  - Clearly state your purpose for connecting and provide a valid email address for communication.
  - Upon generation, the key becomes your responsibility, and ensuring its security is crucial.

Please contact project administrators to request the appropriate API key based on your use case.



## Authenticating

The server uses a key based authentication system. Keys are generated manually on request.

#### Websocket Authentication:
After recieving your `client_id` you will want to send your api key as the next message event with the action event name as `x-api-key` the data being a string (your key), if successful you should recieve the event: `key-accepted`

Example:
```go
WebsocketEvent{
		Client_id: "your id",
		Action: "x-api-key",
		Data: "your key",
	}
```

#### HTTP Example:
When using regular http endpoints that are protected. (ex: need api key) then you will need to send your api key inside the `x-api-key` header

Example:
```json
{
  "method": "GET",
  "path": "/api/v1/your-endpoint",
  "http_version": "HTTP/1.1",
  "headers": {
    "Host": "example.com",
    "x-api-key": "your-api-key"
  }
}
```
