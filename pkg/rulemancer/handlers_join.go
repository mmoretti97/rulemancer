/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) joinRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Post("/available/{gameRef}", e.availableRoom) // Join the first available room for the specified game
		r.Post("/room/{roomID}", e.joinRoom)            // Join a specific room by ID
		r.Post("/new/{gameRef}", e.newGameRoom)         // Create a new room for the specified game and join it
	})
}

func (e *Engine) availableRoom(w http.ResponseWriter, r *http.Request) {
	gameRef := chi.URLParam(r, "gameRef")
	_, claims, err := jwtauth.FromContext(r.Context())
	var client *Client
	var clientID string
	var game *Game

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("Unauthorized available room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if id, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("Unauthorized available room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		clientID = id
		if c, err := e.searchClient(clientID); err != nil {
			// Client existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
				l.Printf("Client not found: %s", clientID)
			}
			Error(w, http.StatusNotFound, "client not found")
			return
		} else {
			client = c
			if g, err := e.searchGame(gameRef); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
					l.Printf("Game not found: %s", gameRef)
				}
				Error(w, http.StatusNotFound, "game not found")
				return
			} else {
				game = g
			}
		}
	}

	game.roomsMutex.Lock()
	if len(game.partialRooms) == 0 {
		game.roomsMutex.Unlock()
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("No available rooms for game: %s, creating a new one", gameRef)
		}
		e.newGameRoom(w, r)
		return
	}

	defer game.roomsMutex.Unlock()

	var room *Room
	var roomId string

	for id, r := range game.partialRooms {
		room = r
		roomId = id
		break
	}

	// Start locking the room
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	if len(room.clients) >= room.maxClients {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("Room is full: %s", roomId)
		}
		Error(w, http.StatusForbidden, "room is full")
		return
	}
	if _, exists := room.clients[clientID]; exists {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("Client already in room: %s", roomId)
		}
		Error(w, http.StatusConflict, "client already in room")
		return
	}

	// Ok the room! now lock the client
	client.roomsMutex.Lock()
	defer client.roomsMutex.Unlock()

	if _, exists := client.playingRooms[roomId]; exists {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("Client already playing in room: %s", roomId)
		}
		Error(w, http.StatusConflict, "client already playing in room")
		return
	}

	if len(room.clients) == room.maxClients-1 {
		// Room is about to be full, remove it from partial rooms and place it
		// on full rooms
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/availableRoom]")+" ", 0)
			l.Printf("Room is full, placing it in running rooms: %s", roomId)
		}
		delete(game.partialRooms, roomId)
		game.runningRooms[roomId] = room
	}

	// Remove from watching if any
	room.watchersMutex.Lock()
	delete(room.watchers, clientID)
	room.watchersMutex.Unlock()

	client.watchersMutex.Lock()
	delete(client.watchingRooms, roomId)
	client.watchersMutex.Unlock()

	// Apply the join to both the room and the client
	room.clients[clientID] = client
	client.playingRooms[roomId] = room
	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/availableRoom]")+" ", 0)
		l.Printf("Client %s joined room: %s", clientID, roomId)
	}
	JSON(w, http.StatusOK, map[string]string{"status": "room found and joined"})

}

func (e *Engine) joinRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomID")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
			l.Printf("Unauthorized join room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
			l.Printf("Unauthorized join room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		if room, err := e.searchRoom(roomId); err != nil {
			// Room existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Room not found: %s", roomId)
			}
			Error(w, http.StatusNotFound, "room not found")
			return
		} else if client, err := e.searchClient(clientID); err != nil {
			// Client existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Client not found: %s", clientID)
			}
			Error(w, http.StatusNotFound, "client not found")
			return
		} else {
			game := room.game

			// Start locking the room
			room.clientsMutex.Lock()
			defer room.clientsMutex.Unlock()

			if len(room.clients) >= room.maxClients {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Room is full: %s", roomId)
				}
				Error(w, http.StatusForbidden, "room is full")
				return
			}
			if _, exists := room.clients[clientID]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Client already in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already in room")
				return
			}

			// Ok the room! now lock the client
			client.roomsMutex.Lock()
			defer client.roomsMutex.Unlock()

			if _, exists := client.playingRooms[roomId]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Client already playing in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already playing in room")
				return
			}

			// We expect to have the room inserted into the game's partial rooms.
			// Lock the game rooms
			game.roomsMutex.Lock()
			defer game.roomsMutex.Unlock()

			if _, ok := game.partialRooms[roomId]; !ok {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Room not found in game's partial rooms: %s", roomId)
				}
				Error(w, http.StatusNotFound, "room not found in game's partial rooms")
				return
			} else {
				if len(room.clients) == room.maxClients-1 {
					// Room is about to be full, remove it from partial rooms and place it
					// on full rooms
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/joinRoom]")+" ", 0)
						l.Printf("Room is full, placing it in running rooms: %s", roomId)
					}
					delete(game.partialRooms, roomId)
					game.runningRooms[roomId] = room
				}
			}

			// Remove from watching if any
			room.watchersMutex.Lock()
			delete(room.watchers, clientID)
			room.watchersMutex.Unlock()

			client.watchersMutex.Lock()
			delete(client.watchingRooms, roomId)
			client.watchersMutex.Unlock()

			// Apply the join to both the room and the client
			room.clients[clientID] = client
			client.playingRooms[roomId] = room
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Client %s joined room: %s", clientID, roomId)
			}
			JSON(w, http.StatusOK, map[string]string{"status": "joined"})
			return
		}
	}
}

func (e *Engine) newGameRoom(w http.ResponseWriter, r *http.Request) {
	gameRef := chi.URLParam(r, "gameRef")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/newGameRoom]")+" ", 0)
			l.Printf("Unauthorized new game room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/newGameRoom]")+" ", 0)
			l.Printf("Unauthorized new game room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		if client, err := e.searchClient(clientID); err != nil {
			// Client existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/newGameRoom]")+" ", 0)
				l.Printf("Client not found: %s", clientID)
			}
			Error(w, http.StatusNotFound, "client not found")
			return
		} else {
			if _, err := e.searchGame(gameRef); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/newGameRoom]")+" ", 0)
					l.Printf("Game not found: %s", gameRef)
				}
				Error(w, http.StatusNotFound, "game not found")
				return
			}

			var room *Room

			if newRoom, err := e.newRoom(clientID+"room", clientID+"room", gameRef); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/newGameRoom]")+" ", 0)
					l.Printf("Failed to create new room: %v", err)
				}
				Error(w, http.StatusInternalServerError, "failed to create new room")
				return
			} else {
				room = newRoom
			}

			roomId := room.id
			game := room.game

			// Start locking the room
			room.clientsMutex.Lock()
			defer room.clientsMutex.Unlock()

			if len(room.clients) >= room.maxClients {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Room is full: %s", roomId)
				}
				Error(w, http.StatusForbidden, "room is full")
				return
			}
			if _, exists := room.clients[clientID]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Client already in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already in room")
				return
			}

			// Ok the room! now lock the client
			client.roomsMutex.Lock()
			defer client.roomsMutex.Unlock()

			if _, exists := client.playingRooms[roomId]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Client already playing in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already playing in room")
				return
			}

			// We expect to have the room inserted into the game's partial rooms.
			// Lock the game rooms
			game.roomsMutex.Lock()
			defer game.roomsMutex.Unlock()

			if _, ok := game.partialRooms[roomId]; !ok {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Room not found in game's partial rooms: %s", roomId)
				}
				Error(w, http.StatusNotFound, "room not found in game's partial rooms")
				return
			} else {
				if len(room.clients) == room.maxClients-1 {
					// Room is about to be full, remove it from partial rooms and place it
					// on full rooms
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/joinRoom]")+" ", 0)
						l.Printf("Room is full, placing it in running rooms: %s", roomId)
					}
					delete(game.partialRooms, roomId)
					game.runningRooms[roomId] = room
				}
			}

			// Remove from watching if any
			room.watchersMutex.Lock()
			delete(room.watchers, clientID)
			room.watchersMutex.Unlock()

			client.watchersMutex.Lock()
			delete(client.watchingRooms, roomId)
			client.watchersMutex.Unlock()

			// Apply the join to both the room and the client
			room.clients[clientID] = client
			client.playingRooms[roomId] = room
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Client %s joined room: %s", clientID, roomId)
			}
			JSON(w, http.StatusOK, map[string]string{"status": "room created and joined"})
			return
		}
	}
}
