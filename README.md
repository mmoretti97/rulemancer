# Rulemancer

<p align="center">
  <img src="logo.png" alt="Rulemancer Logo"/>
</p>

A Go application that embeds the CLIPS expert system engine to power rules-based games and direct JSON-to-CLIPS bridges. Define logic using CLIPS (an expressive rule/fact inference engine), then interact with it via HTTP or CLI.

## Features

- **CLIPS Integration**: Leverage CLIPS for complex rule-based inference and fact management
- **Multi-Game Support**: Host multiple game types simultaneously with dynamic game loading
- **Room-Based Multiplayer**: Create isolated game rooms for concurrent sessions
- **Direct Bridge Mode**: Create bridge rooms that accept raw JSON facts and query arbitrary CLIPS relations
- **HTTP API**: Comprehensive REST endpoints for system, game, and room management with TLS support
- **Real-Time WebSocket Support**: Subscribe to room events and receive live notifications when actions occur
- **Flexible Configuration**: JSON-based configuration with support for multiple game definitions
- **CLI Tools**: Commands for testing, building, and serving games
- **Client Management**: Track and manage connected clients per room

## Quick Start

### Prerequisites

- Git
- Go 1.25+
- C compiler (for CLIPS 6.4 compilation)
- re2c with go-bindings (for the build subcommand)
- JWT secret (set via `RULEMANCER_JWT_SECRET` environment variable or `--secret` flag)

### Clone the Repository

```bash
git clone https://github.com/mmirko/rulemancer.git
cd rulemancer
```

### Project Structure

Inside the project directory, you'll find the following key folders:

- **`core/`** - CLIPS 6.4 C source files and headers (not included in the repo, install using `install-clips.sh`)
- **`cmd/`** - CLI commands (serve, test, build, root)
- **`pkg/rulemancer/`** - Core engine, CLIPS bindings, HTTP handlers, and game management
- **`rulepool/`** - CLIPS rule directories loaded by config
- **`rulepool/tictactoe/`, `rulepool/magic/`** - Game mode rules (multiplayer rooms)
- **`rulepool/bridge/`** - Bridge mode rules (direct JSON<->CLIPS)
- **`interface/`** - Client interface examples and utilities (builded via `rulemancer build`)
- **`testpool/`** - Test rule files for development (unit tests for Tic-Tac-Toe game logic)

### Installation

To install dependencies, compile CLIPS, and build Rulemancer, run from the project root:

```bash
./install-clips.sh
make
```

The above commands will compile CLIPS and build the Rulemancer binary placed in the project root: `./rulemancer`

### Basic Usage

Before starting the server, set the JWT secret as an environment variable:

```bash
export RULEMANCER_JWT_SECRET="your-secret-key-here"
```

Then use the following commands:

- `./rulemancer test` - Run test suite
- `./rulemancer build` - Build extra tools (generates game and bridge shell scripts from templates)
- `./rulemancer serve` - Start HTTPS server (listens on :3000 with TLS)

Once the server is running, it will print an admin JWT token to stdout. The API can be accessed at `https://localhost:3000/api/v1/`

#### Shell Script Templates

The `rulemancer build` command generates shell client interfaces from:

- `pkg/rulemancer/templates/gameshell/` for game rooms
- `pkg/rulemancer/templates/bridgeshell/` for bridge rooms

Scripts are created under `interface/<name>/` for each configured game and bridge, providing convenient command-line access to the API.

#### Documentation

- [API Endpoints](https://github.com/mmirko/rulemancer/blob/master/README-API.md) - Complete API reference for all endpoints
- [Rooms and Games](https://github.com/mmirko/rulemancer/blob/master/README-ROOMS-AND-GAMES.md) - Guide for game mode rooms and bridge mode sessions
- [Game Definition](https://github.com/mmirko/rulemancer/blob/master/README-GAME-DEFINITION.md) - How to define new games using CLIPS

The `rulemancer.json` configuration file can be edited to customize server settings.

## Client Management

Rulemancer uses JWT-based authentication for API access:

- **Admin Token**: Automatically generated and printed to stdout at server startup. Required for system operations (health checks, shutdown). Has the ID `"admin"` in its JWT payload.
- **Client Tokens**: Create clients via the `/api/v1/new/client` endpoint to get individual JWT tokens. Each client receives a unique token for authenticated API access.

### Client Workflow

1. **Create a Client**: POST to `/api/v1/new/client` with name and description - no authentication required for creation
2. **Receive Token**: The response includes a unique JWT token (`api_token`) for that client
3. **Use Token**: Include the token in the `Authorization: Bearer <token>` header for all subsequent API calls
4. **Manage Clients**: Admin can list, view, and delete clients through `/api/v1/client` endpoints

Clients can join game rooms, watch games as spectators, and interact with game logic through the API. See the [Rooms and Games](https://github.com/mmirko/rulemancer/blob/master/README-ROOMS-AND-GAMES.md) guide for more details.

## Real-Time WebSocket Notifications

Rulemancer supports real-time WebSocket connections for monitoring room activities:

- **System Monitor**: Admin-only WebSocket at `/api/v1/system/ws` for system-wide monitoring
- **Room Monitor**: Room-specific WebSocket at `/api/v1/room/{id}/ws` for real-time game updates

Clients connected to a room's WebSocket receive instant notifications when facts are asserted (e.g., when players make moves). This enables real-time game interfaces and live spectator views. See the [Rooms and Games](https://github.com/mmirko/rulemancer/blob/master/README-ROOMS-AND-GAMES.md) guide for WebSocket usage examples.

## Configuration

Edit `rulemancer.json`:

```json
{
  "debug": true,
  "debug_level": 10,
  "tls_cert_file": "server.crt",
  "tls_key_file": "server.key",
  "clipsless_mode": false,
  "games": ["rulepool/tictactoe", "rulepool/magic"],
  "bridges": {"bridge": "rulepool/bridge"}
}
```

### Configuration Options

- **debug**: Enable debug logging
- **debug_level**: Verbosity level for debugging (0-10)
- **tls_cert_file**: Path to TLS certificate file
- **tls_key_file**: Path to TLS private key file
- **clipsless_mode**: Run without CLIPS for testing purposes
- **games**: Array of game directories to load
- **bridges**: Map of bridge name to CLIPS rules directory. Each entry creates a bridge definition loadable through `/api/v1/bridge/*` and spawnable as bridge rooms via `/api/v1/brroom/*`

## Game Mode

Game mode is the multiplayer flow based on `games` config and `/api/v1/room/*` endpoints (`join`, `watch`, websocket updates).

### Game Definition

Each game directory should contain CLIPS files with:

- **Game metadata** via `game-config` fact:

  ```clips
  (game-config 
    (game-name "TicTacToe")
    (description "Classic 3x3 grid game")
    (num-players 2))
  ```

- **Assertable facts**: Facts that can be asserted by clients
- **Queryable facts**: Facts that can be queried by clients  
- **Response facts**: Facts returned after assertions
- **Game rules**: CLIPS rules implementing game logic

For more details check [README-GAME-DEFINITION.md](https://github.com/mmirko/rulemancer/blob/master/README-GAME-DEFINITION.md).

### Example: Tic-Tac-Toe

See [rulepool/tictactoe.clp](https://github.com/mmirko/rulemancer/blob/master/rulepool/tictactoe.clp) and [rulepool/tictactoemeta.clp](https://github.com/mmirko/rulemancer/blob/master/rulepool/tictactoemeta.clp) for a complete game implementation using CLIPS rules and facts.

## Bridge Mode (Direct JSON <-> CLIPS)

Bridge mode is the direct integration flow based on `bridges` config and `/api/v1/brroom/*` endpoints.

1. Configure one or more bridges in `rulemancer.json` under `bridges`.
2. Create a bridge room with `POST /api/v1/brroom/create`.
3. Send a combined request to `POST /api/v1/brroom/{id}/request` with:
   - `facts`: list of relations to assert
   - `queries`: list of relations to read back

Example request body:

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

Response shape:

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

Bridge routes are documented in detail in [README-API.md](https://github.com/mmirko/rulemancer/blob/master/README-API.md).

## License

See LICENSE file
