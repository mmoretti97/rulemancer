# Rulemancer API Endpoints

The server exposes API routes under `/api/v1`.

## Authentication

All authenticated endpoints require:

```text
Authorization: Bearer <jwt_token>
```

- Admin token: printed to stdout at server startup.
- Client token: returned by `POST /api/v1/new/client`.

## Unauthenticated Routes

- `POST /api/v1/new/client` - Create a client and receive a JWT token
  - Request body: `{"name": "string", "description": "string"}`
  - Response: `{"id": "string", "api_token": "string"}`

## System Routes

Admin only:

- `GET /api/v1/system/health` - Health check
  - Response: `{"status": "OK"}`
- `POST /api/v1/system/quit` - Graceful shutdown
  - Request body: `{"graceful": true}`
  - Response: `{"status": "shutting down"}`
- `WS /api/v1/system/ws` - System monitoring websocket

## Client Routes

- `POST /api/v1/client/create` - Create client (deprecated, use `/new/client`)
- `GET /api/v1/client/list` - List clients
  - Response: `{"clients": ["id1", "id2", ...]}`
- `GET /api/v1/client/current` - Get current client from JWT
  - Response: `{"id": "string", "name": "string", "description": "string"}`
- `GET /api/v1/client/{id}` - Get client details
  - Response: `{"id": "string", "name": "string", "description": "string"}`
- `DELETE /api/v1/client/{id}` - Delete client
  - Response: `{"status": "deleted"}`

## Game Mode API

This mode uses the `games` configuration and room/join/watch flows.

### Game Routes

- `GET /api/v1/game/list` - List available games
  - Response: `{"games": ["game1", "game2", ...]}`
- `GET /api/v1/game/{id}` - Get game details
  - Response: `{"id": "string", "name": "string", "description": "string", "rules": "string", "assertable": {...}, "responses": {...}, "queryable": {...}}`

### Room Routes

- `POST /api/v1/room/create` - Create game room
  - Request body: `{"name": "string", "description": "string", "game_ref": "string"}`
  - Response: `{"id": "string"}`
- `GET /api/v1/room/list` - List active rooms
  - Response: `{"rooms": ["room1", "room2", ...]}`
- `GET /api/v1/room/{id}` - Get room details
  - Response: `{"id": "string", "name": "string", "description": "string", "clips_instance": {...}, "running_game": {...}}`
- `DELETE /api/v1/room/{id}` - Delete room
  - Response: `{"status": "deleted"}`

### Room Sub-Routes

- `POST /api/v1/room/{id}/assert/{assertion}` - Assert facts to room
  - Request body: JSON object with relation names as keys
  - Response: `{"status": "asserted", "response": {...}}`
  - Side effect: broadcasts websocket notification to room clients/watchers
- `POST /api/v1/room/{id}/query/{query}` - Query room facts
  - Response: `{"response": {...}}`
- `GET /api/v1/room/{id}/facts` - Get all room facts (debug mode only, admin only)
  - Response: `{"facts": [...]}`
- `WS /api/v1/room/{id}/ws` - Room websocket (players/watchers)

### Join Routes

- `POST /api/v1/join/available/{gameRef}` - Join first available room or create one
  - Response: `{"room_id": "string", "message": "joined room" | "created and joined new room"}`
- `POST /api/v1/join/room/{roomID}` - Join specific room
  - Response: `{"room_id": "string", "message": "joined room"}`
- `POST /api/v1/join/new/{gameRef}` - Create room and join
  - Response: `{"room_id": "string", "message": "created and joined new room"}`

### Watch Routes

- `POST /api/v1/watch/room/{roomId}` - Start watching room
  - Response: `{"room_id": "string", "message": "watching room"}`
- `POST /api/v1/watch/stop/{roomId}` - Stop watching room
  - Response: `{"room_id": "string", "message": "stopped watching room"}`

## Bridge Mode API

This mode uses the `bridges` configuration and direct JSON-to-CLIPS requests.

### Bridge Routes

- `GET /api/v1/bridge/list` - List loaded bridge IDs
  - Auth: any authenticated token
  - Response: `{"bridges": ["bridgeId1", "bridgeId2", ...]}`
- `GET /api/v1/bridge/{id}` - Get bridge details (by ID or bridge name)
  - Auth: admin only
  - Response: `{"id": "string", "name": "string", "rules": "string"}`

### Bridge Room Routes

- `POST /api/v1/brroom/create` - Create bridge room
  - Auth: any authenticated token
  - Request body: `{"name": "string", "bridge_ref": "bridge_id_or_name"}`
  - Response: `{"id": "string"}`
- `GET /api/v1/brroom/list` - List bridge rooms
  - Auth: admin only
  - Response: `{"brrooms": ["room1", "room2", ...]}`
- `GET /api/v1/brroom/{id}` - Get bridge room details
  - Auth: admin only
  - Response: `{"id": "string", "clips_instance": {...}}`
- `DELETE /api/v1/brroom/{id}` - Delete bridge room
  - Auth: admin only
  - Response: `{"status": "deleted"}`
- `POST /api/v1/brroom/{id}/request` - Assert facts and query relations in one call
  - Auth: any authenticated token
  - Request body:
    - `facts` (optional): array of relation assertions
    - `queries` (optional): array of relation names to query
  - Response: `{"asserted": ["(relation ...)", ...], "response": {"relation": [{...}]}}`
- `GET /api/v1/brroom/{id}/facts` - Get all bridge room facts (debug mode only, admin only)
  - Response: `{"facts": [...]}`

### Bridge Request Payload Format

`POST /api/v1/brroom/{id}/request` envelope:

```json
{
  "facts": [
    {
      "first": {
        "x": ["a"],
        "y": ["v"]
      }
    }
  ],
  "queries": ["first"]
}
```

`facts` supports both per-relation encodings:

- Object encoding (single fact):

```json
{"first": {"x": ["a"], "y": ["v"]}}
```

- Array encoding (multiple grouped slots):

```json
{"first": [{"x": ["a"]}, {"y": ["v"]}]}
```

Example response:

```json
{
  "asserted": ["(first ...)", "..."],
  "response": {
    "first": [
      {"x": "a", "y": "v"}
    ]
  }
}
```

Notes:

- If `facts` is omitted, assertions are skipped.
- If `queries` is omitted, `response` is empty.

## Error Responses

Standard error envelope:

```json
{
  "error": "error message description"
}
```

Common statuses:

- `200 OK`
- `201 Created`
- `400 Bad Request`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`
- `409 Conflict`
- `500 Internal Server Error`
