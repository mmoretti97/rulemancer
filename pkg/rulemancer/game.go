package rulemancer

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
)

type Game struct {
	name          string
	description   string
	id            string
	rulesLocation string
	numPlayers    int
	assertable    map[string][]string
	responses     map[string][]string
	queryable     map[string][]string
	runningRooms  map[string]*Room
	roomsMutex    sync.RWMutex
}

func (g *Game) Info() map[string]any {
	return map[string]any{
		"id":            g.id,
		"name":          g.name,
		"description":   g.description,
		"rulesLocation": g.rulesLocation,
		"assertable":    g.assertable,
		"responses":     g.responses,
		"queryable":     g.queryable,
		"runningRooms":  g.runningRooms,
	}
}

func (e *Engine) loadGames() {
	// Load games from the configured games list
	for _, gameLocation := range e.Games {
		if err := e.newGame(gameLocation); err != nil {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/loadGames]")+" ", 0)
			l.Printf("error loading game from %s: %v", gameLocation, err)
		} else {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/loadGames]")+" ", 0)
				l.Printf("successfully loaded game from %s", gameLocation)
			}
		}
	}
}

func (e *Engine) newGame(rulesLocation string) error {
	// Ensure unique ID generation and locking on the games map
	var cli *ClipsInstance
	cli = e.NewClipsInstance()
	defer cli.Dispose()

	if err := cli.InitClips(); err != nil {
		return err
	}
	// Load knowledge base from the specified game rules location
	if err := cli.loadGame(rulesLocation); err != nil {
		return err
	}

	// Retrieve game configuration facts
	var name string
	var description string
	var numPlayers int

	gc, err := cli.QueryFacts("game-config")
	if err != nil {
		return err
	}
	gcMap, err := genericFactToMap(e.Config, "game-config", gc)
	if err != nil {
		return err
	}
	switch len(gcMap) {
	case 0:
		return errors.New("no game-config found in the rules location")
	case 1:
		// All good
		if gameName, ok := gcMap[0]["game-name"]; ok {
			name = gameName
		} else {
			return errors.New("game-config missing game-name slot")
		}
		if desc, ok := gcMap[0]["description"]; ok {
			description = desc
		} else {
			return errors.New("game-config missing description slot")
		}
		if numPlayersStr, ok := gcMap[0]["num-players"]; ok {
			numPlayersInt, err := strconv.Atoi(numPlayersStr)
			if err != nil {
				return errors.New("game-config num-players slot must be an integer")
			}
			// The numPlayers value is currently not used, but it can be stored in the Game struct for future use
			numPlayers = numPlayersInt
		} else {
			return errors.New("game-config missing num-players slot")
		}
	default:
		return errors.New("multiple game-config facts found in the rules location")
	}

	// Get the assertable facts
	assertableFacts, err := cli.getGameConfig("assertable")
	if err != nil {
		return err
	}

	// Get the response facts
	results, err := cli.getGameConfig("results")
	if err != nil {
		return err
	}

	// Get the queryable facts
	queryableFacts, err := cli.getGameConfig("queryable")
	if err != nil {
		return err
	}

	// The game is successfully loaded, the CLIPS instance can be disposed by deferring

	e.gamesMutex.Lock()
	defer e.gamesMutex.Unlock()
	game := &Game{
		name:          name,
		description:   description,
		rulesLocation: rulesLocation,
		numPlayers:    numPlayers,
		assertable:    assertableFacts,
		responses:     results,
		queryable:     queryableFacts,
		id:            e.generateGameUniqueID(),
		runningRooms:  make(map[string]*Room),
		roomsMutex:    sync.RWMutex{},
	}
	e.numGames++
	e.games[game.id] = game

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/newGame]")+" ", 0)
		l.Printf("Loaded game %s (%s) with ID %s", game.name, game.description, game.id)
		l.Print(game)
	}
	return nil
}

func (e *Engine) generateGameUniqueID() string {
	for {
		newId := RandStringBytes(16)
		if _, exists := e.games[newId]; !exists {
			return newId
		}
	}
}

// Search for a game by its ID or name
func (e *Engine) searchGame(id string) (*Game, error) {
	e.gamesMutex.RLock()
	defer e.gamesMutex.RUnlock()

	// First search by ID
	if game, exists := e.games[id]; exists {
		return game, nil
	}

	// Then search by name
	for _, game := range e.games {
		if game.name == id {
			return game, nil
		}
	}

	return nil, errors.New("game not found")
}

func (e *Engine) listGames() []string {
	e.gamesMutex.RLock()
	defer e.gamesMutex.RUnlock()
	games := make([]string, 0, len(e.games))
	for _, game := range e.games {
		games = append(games, game.id)
	}
	return games
}
