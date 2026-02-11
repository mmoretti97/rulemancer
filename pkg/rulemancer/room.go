/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"errors"
	"sync"
	"time"
)

type Room struct {
	name          string
	description   string
	id            string
	game          *Game
	clients       map[string]*Client
	maxClients    int
	clientsMutex  sync.RWMutex
	watchers      map[string]*Client
	watchersMutex sync.RWMutex
	clipsInstance *ClipsInstance
	lastActive    int64
}

func (e *Engine) newRoom(name, description, gameRef string) (*Room, error) {

	game, err := e.searchGame(gameRef)
	if err != nil {
		return nil, err
	}

	rulesLocation := game.rulesLocation

	// Ensure unique ID generation and locking on the rooms map
	var cli *ClipsInstance
	if !e.ClipsLessMode {
		cli = e.NewClipsInstance()
		if err := cli.InitClips(); err != nil {
			return nil, err
		}
		if err := cli.loadGame(rulesLocation); err != nil {
			cli.Dispose()
			return nil, err
		}
	}
	e.roomsMutex.Lock()
	defer e.roomsMutex.Unlock()
	room := &Room{
		name:          name,
		description:   description,
		id:            e.generateRoomUniqueID(),
		game:          game,
		clipsInstance: cli,
		maxClients:    game.numPlayers,
		clients:       make(map[string]*Client),
		clientsMutex:  sync.RWMutex{},
		watchers:      make(map[string]*Client),
		watchersMutex: sync.RWMutex{},
		lastActive:    time.Now().Unix(),
	}
	e.numRooms++
	e.rooms[room.id] = room

	game.roomsMutex.Lock()
	defer game.roomsMutex.Unlock()
	game.partialRooms[room.id] = room

	return room, nil
}

func (e *Engine) generateRoomUniqueID() string {
	for {
		newId := randStringBytes(16)
		if _, exists := e.rooms[newId]; !exists {
			return newId
		}
	}
}

func (e *Engine) searchRoom(id string) (*Room, error) {
	e.roomsMutex.RLock()
	defer e.roomsMutex.RUnlock()
	if room, exists := e.rooms[id]; exists {
		return room, nil
	}
	return nil, errors.New("room not found")
}

func (e *Engine) removeRoom(id string) (*Room, error) {
	e.roomsMutex.Lock()
	defer e.roomsMutex.Unlock()
	if room, exists := e.rooms[id]; exists {
		if !e.ClipsLessMode {
			room.clipsInstance.Dispose()
		}
		delete(e.rooms, id)
		e.numRooms--
		return room, nil
	}
	return nil, errors.New("room not found")
}

func (e *Engine) listRooms() []string {
	e.roomsMutex.RLock()
	defer e.roomsMutex.RUnlock()
	rooms := make([]string, 0, len(e.rooms))
	for _, room := range e.rooms {
		rooms = append(rooms, room.id)
	}
	return rooms
}
