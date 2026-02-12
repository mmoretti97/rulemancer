# Rulemancer API Endpoints

The server exposes the following API routes under `/api/v1`. All endpoints require JWT authentication unless otherwise noted.

## Authentication

All authenticated endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

- **Admin Token**: Printed to stdout at server startup. Used for system operations.
- **Client Tokens**: Obtained by creating a client via `/api/v1/new/client`.

## New Entity Creation Routes

These endpoints do not require authentication:

- `POST /api/v1/new/client` - Create a new client and receive a JWT token
  - **Request Body**: `{"name": "string", "description": "string"}`
  - **Response**: `{"id": "string", "api_token": "string"}`

## System Routes

Requires admin authentication:

- `GET /api/v1/system/health` - Get system health status
  - **Response**: `{"status": "OK"}`
- `POST /api/v1/system/quit` - Gracefully shutdown the server
  - **Request Body**: `{"graceful": true}`
  - **Response**: `{"status": "shutting down"}`
- `WS /api/v1/system/ws` - WebSocket connection for system-wide monitoring (admin only)
  - **Protocol**: WebSocket (wss://)
  - **Authentication**: JWT token required via Authorization header during handshake
  - **Purpose**: Real-time system monitoring and notifications

## Client Routes

- `POST /api/v1/client/create` - Create a new client (deprecated, use `/new/client` instead)
- `GET /api/v1/client/list` - List all clients
  - **Response**: `{"clients": ["id1", "id2", ...]}`
- `GET /api/v1/client/current` - Get current client information (from JWT token)
  - **Response**: `{"id": "string", "name": "string", "description": "string"}`
- `GET /api/v1/client/{id}` - Get client details by ID
  - **Response**: `{"id": "string", "name": "string", "description": "string"}`
- `DELETE /api/v1/client/{id}` - Delete a client
  - **Response**: `{"status": "deleted"}`

## Game Routes

- `GET /api/v1/game/list` - List all available games
  - **Response**: `{"games": ["game1", "game2", ...]}`
- `GET /api/v1/game/{id}` - Get game details including assertables, queryables, and responses
  - **Response**: `{"id": "string", "name": "string", "description": "string", "rules": "string", "assertable": {...}, "responses": {...}, "queryable": {...}}`

## Room Routes

- `POST /api/v1/room/create` - Create a new game room
  - **Request Body**: `{"name": "string", "description": "string", "game_ref": "string"}`
  - **Response**: `{"id": "string"}`
- `GET /api/v1/room/list` - List all active rooms
  - **Response**: `{"rooms": ["room1", "room2", ...]}`
- `GET /api/v1/room/{id}` - Get room details
  - **Response**: `{"id": "string", "name": "string", "description": "string", "clips_instance": {...}, "running_game": {...}}`
- `DELETE /api/v1/room/{id}` - Delete a room
  - **Response**: `{"status": "deleted"}`

### Room Sub-Routes

- `POST /api/v1/room/{id}/assert/{assertion}` - Assert facts to a room
  - **Request Body**: JSON object with relation names as keys and arrays of fact objects as values
  - **Response**: `{"response": {...}}` - Returns the result facts as defined in the game's response configuration
  - **Side Effect**: Broadcasts notification to all WebSocket connections on the room
- `POST /api/v1/room/{id}/query/{query}` - Query facts from a room
  - **Response**: `{"response": {...}}` - Returns the queried facts
- `GET /api/v1/room/{id}/facts` - Get all facts from a room (debug mode only)
  - **Response**: `{"facts": [...]}`
- `WS /api/v1/room/{id}/ws` - WebSocket connection for real-time room updates
  - **Protocol**: WebSocket (wss://)
  - **Authentication**: JWT token required via Authorization header during handshake
  - **Access**: Available to room clients (players) and watchers (spectators)
  - **Notifications**: Receives messages when facts are asserted in the room
  - **Message Format**: Text messages like `"asserted (move x 1 y 1 player x)"`

## Join Routes

These routes allow clients to join or create game rooms:

- `POST /api/v1/join/available/{gameRef}` - Join the first available room for the specified game, or create a new one if none available
  - **Response**: `{"room_id": "string", "message": "joined room" | "created and joined new room"}`
- `POST /api/v1/join/room/{roomID}` - Join a specific room by ID
  - **Response**: `{"room_id": "string", "message": "joined room"}`
- `POST /api/v1/join/new/{gameRef}` - Create a new room for the specified game and join it
  - **Response**: `{"room_id": "string", "message": "created and joined new room"}`

## Watch Routes

These routes allow clients to watch rooms as spectators (read-only access):

- `POST /api/v1/watch/room/{roomId}` - Start watching a specific room
  - **Response**: `{"room_id": "string", "message": "watching room"}`
- `POST /api/v1/watch/stop/{roomId}` - Stop watching a specific room
  - **Response**: `{"room_id": "string", "message": "stopped watching room"}`

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "error message description"
}
```

Common HTTP status codes:
- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request body or parameters
- `401 Unauthorized` - Missing or invalid JWT token
- `404 Not Found` - Requested resource not found
- `409 Conflict` - Resource conflict (e.g., already joined room)
- `500 Internal Server Error` - Server error