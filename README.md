# Rulemancer

<p align="center">
  <img src="logo.png" alt="Rulemancer Logo"/>
</p>

A Go application that embeds the CLIPS expert system engine to power rules-based games. Define game logic using CLIPS (an expressive rule/fact inference engine), then interact with it via HTTP or CLI.

## Features

- **CLIPS Integration**: Leverage CLIPS for complex rule-based inference and fact management
- **Multi-Game Support**: Host multiple game types simultaneously with dynamic game loading
- **Room-Based Multiplayer**: Create isolated game rooms for concurrent sessions
- **HTTP API**: Comprehensive REST endpoints for system, game, and room management with TLS support
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
- **`rulepool/`** - CLIPS rule files (`.clp`) for game definitions (e.g., Tic-Tac-Toe example)
- **`interface/`** - Client interface examples and utilities (builded via `rulemancer build`)
- **`testpool/`** - Test rule files for development (unit tests for Tic-Tac-Toe game logic)

### Installation

To install dependencies, compile CLIPS, and build Rulemancer, run from the project root:

```bash
./install-clips.sh
make
```

The above commands will compile CLIPS and build the Rulemancer binary placed in the project root: `./rulemancer`

### 

### Basic Usage

Before starting the server, set the JWT secret as an environment variable:

```bash
export RULEMANCER_JWT_SECRET="your-secret-key-here"
```

Then use the following commands:

- `./rulemancer test` - Run test suite
- `./rulemancer build` - Build the extra tools (generates client shell scripts from templates)
- `./rulemancer serve` - Start HTTPS server (listens on :3000 with TLS)

Once the server is running, it will print an admin JWT token to stdout. The API can be accessed at `https://localhost:3000/api/v1/`

#### Shell Script Templates

The `rulemancer build` command generates shell client interfaces from templates located in `pkg/rulemancer/templates/shell/`. These scripts are created in the `interface/` folder for each available game, providing convenient command-line access to the API.

#### Documentation

- [API Endpoints](README-API.md) - Complete API reference for all endpoints
- [Rooms and Games](README-ROOMS-AND-GAMES.md) - Guide for creating and managing game rooms
- [Game Definition](README-GAME-DEFINITION.md) - How to define new games using CLIPS

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

Clients can join game rooms, watch games as spectators, and interact with game logic through the API. See the [Rooms and Games](README-ROOMS-AND-GAMES.md) guide for more details.

## Configuration

Edit `rulemancer.json`:

```json
{
  "debug": true,
  "debug_level": 10,
  "tls_cert_file": "server.crt",
  "tls_key_file": "server.key",
  "clipsless_mode": false,
  "games": ["rulepool"]
}
```

### Configuration Options

- **debug**: Enable debug logging
- **debug_level**: Verbosity level for debugging (0-10)
- **tls_cert_file**: Path to TLS certificate file
- **tls_key_file**: Path to TLS private key file
- **clipsless_mode**: Run without CLIPS for testing purposes
- **games**: Array of game directories to load

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

for more details check the [Game Definition](README-GAME-DEFINITION.md) document.

## Example: Tic-Tac-Toe

See [rulepool/tictactoe.clp](rulepool/tictactoe.clp) and [rulepool/tictactoemeta.clp](rulepool/tictactoemeta.clp) for a complete game implementation using CLIPS rules and facts.

## License

See LICENSE file
