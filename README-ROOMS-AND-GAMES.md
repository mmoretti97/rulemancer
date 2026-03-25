# Rulemancer Rooms and Games Management

This guide explains how to manage game rooms, join games, and interact with the CLIPS-based game logic in Rulemancer.

## Overview

Rulemancer uses a room-based architecture where:

- **Games** are CLIPS-based rule definitions that define game logic
- **Rooms** are isolated game instances where clients can play
- **Clients** are authenticated users who can join rooms and interact with games
- **Watchers** are clients who observe rooms in read-only mode

This document is split into two independent parts:

- **Game Mode**: rooms, join/watch workflows, assert/query endpoints
- **Bridge Mode**: direct JSON<->CLIPS sessions via bridge rooms

## Core Concepts

### Games

Games are defined using CLIPS files and loaded at server startup. Each game specifies:

- **Assertables**: Facts that clients can assert (e.g., player moves)
- **Responses**: Facts returned after assertions (e.g., move validation results)
- **Queryables**: Facts that clients can query (e.g., game state, winner)

List available games:

```bash
GET /api/v1/game/list
```

Get game details:

```bash
GET /api/v1/game/{game_id}
```

### Rooms

Rooms are instances of games. Each room has:

- A unique ID
- A reference to a game
- A CLIPS instance for that specific game session
- A list of clients (players)
- A list of watchers (spectators)
- Maximum client capacity (defined by the game)

### Clients

Clients are authenticated entities with JWT tokens. They can:

- Create and join rooms
- Assert facts to rooms they've joined
- Query room state
- Watch rooms as spectators

## Getting Started

### 1. Create a Client

First, create a client to get a JWT token (no authentication required):

```bash
curl -k -X POST https://localhost:3000/api/v1/new/client \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "description": "Player 1"}'
```

Response:

```json
{
  "id": "client-123abc",
  "api_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Important**: Save the `api_token` for all subsequent requests:

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 2. Join a Game

There are three ways to join a game:

#### Option A: Join Any Available Room

Join the first available room for a game, or create a new one if none exists:

```bash
curl -k -X POST https://localhost:3000/api/v1/join/available/tictactoe \
  -H "Authorization: Bearer $API_TOKEN"
```

#### Option B: Create a New Room and Join

Create a new room and join it immediately:

```bash
curl -k -X POST https://localhost:3000/api/v1/join/new/tictactoe \
  -H "Authorization: Bearer $API_TOKEN"
```

#### Option C: Join a Specific Room

Join a specific room by its ID:

```bash
curl -k -X POST https://localhost:3000/api/v1/join/room/{room_id} \
  -H "Authorization: Bearer $API_TOKEN"
```

All join endpoints return:

```json
{
  "room_id": "room-456def",
  "message": "joined room"
}
```

### 3. Interact with the Game

#### Assert Facts (Make Moves)

Once in a room, assert facts to interact with the game. The exact structure depends on the game's assertables:

```bash
curl -k -X POST https://localhost:3000/api/v1/room/{room_id}/assert/move \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "move": [{
      "x": ["1"],
      "y": ["1"],
      "player": ["x"]
    }]
  }'
```

Response contains the result facts as defined by the game:
```json
{
  "response": {
    "last-move": [{
      "valid": "yes",
      "reason": "Move accepted"
    }]
  }
}
```

#### Query Game State

Query the current state of the game:

```bash
curl -k -X POST https://localhost:3000/api/v1/room/{room_id}/query/cell \
  -H "Authorization: Bearer $API_TOKEN"
```

Response:

```json
{
  "response": {
    "cell": [
      {"x": "1", "y": "1", "value": "x"},
      {"x": "1", "y": "2", "value": "empty"},
      ...
    ]
  }
}
```

Query for winner:

```bash
curl -k -X POST https://localhost:3000/api/v1/room/{room_id}/query/winner \
  -H "Authorization: Bearer $API_TOKEN"
```

## Watching Rooms

Clients can watch rooms as spectators without joining as players:

### Start Watching

```bash
curl -k -X POST https://localhost:3000/api/v1/watch/room/{room_id} \
  -H "Authorization: Bearer $API_TOKEN"
```

### Stop Watching

```bash
curl -k -X POST https://localhost:3000/api/v1/watch/stop/{room_id} \
  -H "Authorization: Bearer $API_TOKEN"
```

Watchers can:

- Query room state
- View all game facts
- Connect to the room's WebSocket for real-time updates

Watchers cannot:

- Assert facts
- Make moves
- Affect game state

## Real-Time WebSocket Notifications

Rulemancer supports WebSocket connections for real-time monitoring of room activities. This enables live game updates, spectator views, and interactive UI development.

### Room WebSocket Endpoint

```
wss://localhost:3000/api/v1/room/{room_id}/ws
```

**Authentication**: JWT token must be provided in the Authorization header during the WebSocket handshake.

**Access Control**: Only clients who have joined the room (as players) or are watching the room (as spectators) can connect to the room's WebSocket.
**Notifications**: Whenever a fact is asserted in the room (e.g., a player makes a move), all connected WebSocket clients receive a notification with the asserted fact. This allows for real-time updates in game interfaces and live spectator views.

## Room Management

### Create a Room Manually

```bash
curl -k -X POST https://localhost:3000/api/v1/room/create \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My TicTacToe Game",
    "description": "A friendly match",
    "game_ref": "tictactoe"
  }'
```

### List All Rooms

```bash
curl -k -X GET https://localhost:3000/api/v1/room/list \
  -H "Authorization: Bearer $API_TOKEN"
```

### Get Room Details

```bash
curl -k -X GET https://localhost:3000/api/v1/room/{room_id} \
  -H "Authorization: Bearer $API_TOKEN"
```

### Delete a Room

```bash
curl -k -X DELETE https://localhost:3000/api/v1/room/{room_id} \
  -H "Authorization: Bearer $API_TOKEN"
```

## Bridge Mode (Direct JSON <-> CLIPS)

Use this mode when you need a direct bridge instead of player/join/watch room mechanics.

Bridge mode is separate from game rooms. A bridge room has:

- A unique ID
- A reference to a loaded bridge (`bridges` in `rulemancer.json`)
- A dedicated CLIPS instance
- A generic request endpoint: `POST /api/v1/brroom/{id}/request`

### Create a Bridge Room

```bash
curl -k -X POST https://localhost:3000/api/v1/brroom/create \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "bridge-session-1",
    "bridge_ref": "bridge"
  }'
```

### Send a Bridge Request

```bash
curl -k -X POST https://localhost:3000/api/v1/brroom/bridge-session-1/request \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "facts": [
      {
        "first": {
          "x": ["a"],
          "y": ["v"]
        }
      }
    ],
    "queries": ["first"]
  }'
```

Example response:

```json
{
  "asserted": ["(first ...)"],
  "response": {
    "first": [
      {"x": "a", "y": "v"}
    ]
  }
}
```

### Inspect and Delete Bridge Rooms

Admin-only endpoints:

- `GET /api/v1/brroom/list`
- `GET /api/v1/brroom/{id}`
- `DELETE /api/v1/brroom/{id}`

Bridge discovery endpoints:

- `GET /api/v1/bridge/list` (authenticated)
- `GET /api/v1/bridge/{id_or_name}` (admin only)

## Complete Workflow Example

Here's a complete example of two players playing Tic-Tac-Toe:

### Player 1 (X)

```bash
# 1. Create client
curl -k -X POST https://localhost:3000/api/v1/new/client \
  -H "Content-Type: application/json" \
  -d '{"name": "PlayerX", "description": "Player X"}' \
  > player1.json

export TOKEN_X=$(cat player1.json | jq -r '.api_token')

# 2. Join or create a game
curl -k -X POST https://localhost:3000/api/v1/join/new/tictactoe \
  -H "Authorization: Bearer $TOKEN_X" \
  > room.json

export ROOM_ID=$(cat room.json | jq -r '.room_id')

# 3. Make first move
curl -k -X POST https://localhost:3000/api/v1/room/$ROOM_ID/assert/move \
  -H "Authorization: Bearer $TOKEN_X" \
  -H "Content-Type: application/json" \
  -d '{"move": [{"x": ["1"], "y": ["1"], "player": ["x"]}]}'
```

### Player 2 (O)

```bash
# 1. Create client
curl -k -X POST https://localhost:3000/api/v1/new/client \
  -H "Content-Type: application/json" \
  -d '{"name": "PlayerO", "description": "Player O"}' \
  > player2.json

export TOKEN_O=$(cat player2.json | jq -r '.api_token')

# 2. Join the same room
curl -k -X POST https://localhost:3000/api/v1/join/room/$ROOM_ID \
  -H "Authorization: Bearer $TOKEN_O"

# 3. Make move
curl -k -X POST https://localhost:3000/api/v1/room/$ROOM_ID/assert/move \
  -H "Authorization: Bearer $TOKEN_O" \
  -H "Content-Type: application/json" \
  -d '{"move": [{"x": ["2"], "y": ["2"], "player": ["o"]}]}'

# 4. Check game state
curl -k -X POST https://localhost:3000/api/v1/room/$ROOM_ID/query/cell \
  -H "Authorization: Bearer $TOKEN_O" | jq .

# 5. Check for winner
curl -k -X POST https://localhost:3000/api/v1/room/$ROOM_ID/query/winner \
  -H "Authorization: Bearer $TOKEN_O" | jq .
```

## Using Shell Scripts

The `rulemancer build` command generates shell scripts in `interface/<game_name>/` for easier interaction:

```bash
# Build the shell scripts
./rulemancer build

# Navigate to the game interface
cd interface/tictactoe/

# Source the launch script to start server and get admin token
source launch.sh

# Create a client (token automatically exported)
source client-create.sh

# Join a new game
./join-new.sh tictactoe

# Make a move (scripts are game-specific)
./assert_move.sh $ROOM_ID 1 1 x

# Query state
./query.sh $ROOM_ID cell
```

## Debugging

### View All Facts in a Room

In debug mode, you can view all CLIPS facts in a room:

```bash
curl -k -X GET https://localhost:3000/api/v1/room/{room_id}/facts \
  -H "Authorization: Bearer $API_TOKEN"
```

This endpoint is only available when the server is running with `"debug": true` in the configuration.

## Best Practices

1. **Token Management**: Keep client tokens secure. Each client should have their own token.

2. **Room Lifecycle**: Delete rooms when games are finished to free resources.

3. **Error Handling**: Always check response status and error messages. Games may reject invalid assertions.

4. **Game Rules**: Understand the game's assertables, responses, and queryables by inspecting the game details via `/api/v1/game/{game_id}`.

5. **Concurrent Access**: Rooms use mutex locks to ensure thread-safe access. Multiple clients can safely interact with the same room.

## Troubleshooting

### "Room not found"

- The room may have been deleted
- Verify the room_id is correct
- List all rooms to see available options

### "Unauthorized"

- Check that your JWT token is valid
- Ensure the token is included in the Authorization header
- Verify the token hasn't expired

### "Assertion not found"

- The assertion name doesn't match the game's assertables
- Check the game details to see valid assertion names

### "Query not found"

- The query name doesn't match the game's queryables
- Check the game details to see valid query names

### "Invalid JSON body"

- Verify JSON syntax is correct
- Ensure relation names match the game's expected structure
- Check that all required fields are provided

## See Also

- [API Endpoints](https://github.com/mmirko/rulemancer/blob/master/README-API.md) - Complete API reference
- [Game Definition](https://github.com/mmirko/rulemancer/blob/master/README-GAME-DEFINITION.md) - How to create new games
- [Main README](https://github.com/mmirko/rulemancer/blob/master/README.md) - Installation and configuration