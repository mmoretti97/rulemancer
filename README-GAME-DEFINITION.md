# Rulemancer Game Setup Guide

This guide explains how to set up CLIPS rule files to enable a game to work with Rulemancer. This only covers game mode rooms.

## Overview

Rulemancer uses CLIPS (C Language Integrated Production System) to define game logic. To integrate a new game, you need to create a `.clp` file that defines:

1. **Game Configuration** - Basic metadata about your game
2. **Game Interface** - How the game interacts with external systems through assertables, results, and queryables
3. **Game Logic** - The actual rules and templates that define the game's behavior (not covered in this guide)

## Required Schema

Your game file must respect the schema defined in `rulepool/common.clp`. This schema defines three core templates:

### 1. `game-config`
Provides basic metadata about your game. In particular, it includes the game's name and description.

**Structure:**
```clips
(deftemplate game-config
  (slot game-name)
  (slot description)
  (slot num-players))
```

### 2. `assertable`
Defines facts that can be asserted into the CLIPS environment from external sources (e.g., player moves, game actions).
The main purpose of the engine is to serve these assertables as REST endpoints.

**Structure:**
```clips
(deftemplate assertable
  (slot name)
  (multislot relations))
```

### 3. `results`
Defines facts that are generated as results of game logic and should be returned to external systems. 
These are the outputs of the REST endpoints.

**Structure:**
```clips
(deftemplate results
  (slot name)
  (multislot relations))
```

### 4. `queryable`
Defines facts that can be queried from the CLIPS environment (e.g., game state, winner information).
The engine exposes these queryables as REST endpoints for retrieving game state information, independent of the main game interaction.

**Structure:**
```clips
(deftemplate queryable
  (slot name)
  (multislot relations))
```

## Step-by-Step Setup

This workflow is specifically for games loaded via the `games` config key.

### Step 1: Create Your Game File

Create a new `.clp` file in the `rulepool/` directory, e.g., `rulepool/yourgamemeta.clp`.

### Step 2: Define Game Configuration

Use `deffacts` to declare your game configuration:

```clips
(deffacts yourgame-config
  (game-config
    (game-name YourGame)
    (description "A description of your game and its rules.")
    (num-players 2))
```

**Fields:**

- `game-name`: Identifier for your game (no spaces, use lowercase)
- `description`: Human-readable description of the game
- `num-players`: The number of players required to play the game

### Step 3: Define Game Interface

Use `deffacts` to declare how your game interfaces with the outside world:

```clips
(deffacts yourgame-interface
  ; Define what can be asserted
  (assertable
    (name action-type-1)
    (relations relation-name-1))
  
  ; Define what results are produced
  (results 
    (name result-type-1)
    (relations result-relation-1))
  
  ; Define what can be queried
  (queryable
    (name query-type-1)
    (relations relation-name-1 relation-name-2))
)
```

#### Assertables

Assertables define **inputs** to your game - actions or information that can be pushed into the CLIPS environment:

- `name`: The type of fact that can be asserted
- `relations`: The relation name(s) used in the actual CLIPS facts

**Example:**

```clips
(assertable
  (name move)
  (relations move))
```

This means external systems can assert facts like `(move ...)` into CLIPS.

**Real-Time Notifications:** When facts are asserted via the API, all clients connected to the room's WebSocket (`/api/v1/room/{id}/ws`) receive real-time notifications with the asserted fact. This enables live game updates and spectator views.

#### Results

Results define **outputs** from your game logic - facts that should be returned after processing:

- `name`: The type of result fact
- `relations`: The relation name(s) that will be matched in results

**Example:**

```clips
(results 
  (name move)
  (relations last-move))
```

This means the system will look for facts like `(last-move ...)` to return as results.

#### Queryables

Queryables define what information can be **queried** from the current game state:

- `name`: The type of query
- `relations`: The relation name(s) that can be queried

**Example:**

```clips
(queryable
  (name winner)
  (relations winner cell))
```

This means external systems can query for facts matching `(winner ...)` or `(cell ...)`.

## Step 4: Implement Game Logic

After creating your metadata file, the interface is defined. The next steps involve implementing the actual game logic:

1. **Create Game Logic Files**: Write CLIPS rules that implement your game logic, not necessarily all the game logic needs to be exposed via the interface. Actually, it is the opposite. Most of the game logic should be internal and only a few relevant facts should be exposed via the interface.
2. **Define Game Templates**: Create templates for your game-specific exportable facts, results, and queryables
3. **Configure Rulemancer**: Ensure your game file is included in the Rulemancer configuration so it gets loaded properly.
4. **Test**: Optionally, create test files to validate your game logic and interface and test them using the `rulemancer test` command.

The game is ready to be served via Rulemancer!

## Step 5 (Optional): Shell interface

You can also create a shell interface to interact with your game via command line (using `curl` commands). The `rulemancer build` command can help you set this up by generating the necessary shell scripts based on your game metadata. By default, the shell interface will be created in the `interfaces/gameshell/` directory.
The `rulemancer build` command will parse your game CLIPS file to get the assertables, results, and queryables templates and generate the corresponding shell scripts to interact with your game.

## Complete Example: Tic-Tac-Toe

Here's the complete metadata file for Tic-Tac-Toe (`rulepool/tictactoemeta.clp`):

```clips
(deffacts tictactoe-config
  (game-config
    (game-name tictactoe)
    (description "A simple Tic Tac Toe game between two players.")
    (num-players 2))

(deffacts tictactoe-interface
  (assertable
    (name move)
    (relations move))
  (results 
    (name move)
    (relations last-move))
  (queryable
    (name winner)
    (relations winner cell))
  (queryable
    (name cell)
    (relations cell)))
```

### What This Means:

1. **Game Config**: Identifies the game as "tictactoe"

2. **Assertable `move`**: External systems can assert move facts like:

Looking at the `move` template, external systems can assert facts like:

   ```clips
  (move (x 1) (y 1) (player x))  ; Player x places x at (1,1)
   ```

3. **Results `move`**: After processing, the system returns facts like:

   ```clips
   (last-move (valid no) (reason "Cell already occupied"))  ; Invalid move result
   (last-move (valid yes) (reason "Move accepted"))  ; Valid move result
   ```

4. **Queryable `winner`**: Can query for winner status:

   ```clips
   (winner (player x))  ; Player x has won
   ```

5. **Queryable `cell`**: Can query the board state:

   ```clips
   (cell (x 1) (y 1) (value x))  ; Cell (1,1) is occupied by x
   (cell (x 2) (y 2) (value o))  ; Cell (2,2) is occupied by o
   ...  ; Other cells
   ```



## Best Practices

1. **Naming Conventions**:
   - Use descriptive names for `game-name` (lowercase, no spaces)
   - Keep relation names short but meaningful
   - Be consistent with naming across your game files

2. **Separation of Concerns**:
   - Keep metadata (interface definitions) in a separate `*meta.clp` file
   - Keep game logic (rules, templates, functions) in separate files

3. **Documentation**:
   - Always provide a clear `description` in your game config
   - Comment your code to explain complex rules

4. **Testing**:
   - Test that assertables work by asserting facts and checking the results
   - Test that queryables return expected game state
   - Verify results are properly generated

5. **Real-Time Features**:
   - Design your assertables with real-time notifications in mind
   - Clients can subscribe to room WebSockets to receive live updates
   - Consider how spectators will experience the game through WebSocket broadcasts

## Additional Resources

- CLIPS Documentation: Learn more about CLIPS syntax and features
- Example Games: Check `rulepool/` directory for complete game implementations
- Rulemancer Documentation: See main README.md for system architecture

---

For questions or issues, please refer to the main project documentation or open an issue on the project repository.
