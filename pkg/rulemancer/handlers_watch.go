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

func (e *Engine) watchRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Post("/room/{roomId}", e.watchRoom)   // Watch a specific room by ID (read-only)
		r.Post("/stop/{roomId}", e.unwatchRoom) // Unwatch a specific room by ID
	})
}

func (e *Engine) watchRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
			l.Printf("Unauthorized watch room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
			l.Printf("Unauthorized watch room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		if room, err := e.searchRoom(roomId); err != nil {
			// Room existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
				l.Printf("Room not found: %s", roomId)
			}
			Error(w, http.StatusNotFound, "room not found")
			return
		} else if client, err := e.searchClient(clientID); err != nil {
			// Client existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
				l.Printf("Client not found: %s", clientID)
			}
			Error(w, http.StatusNotFound, "client not found")
			return
		} else {
			// Start locking the room
			room.clientsMutex.RLock()

			if _, exists := room.clients[clientID]; exists {
				room.clientsMutex.RUnlock()
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
					l.Printf("Client already playing in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already playing in room")
				return
			}

			room.clientsMutex.RUnlock()

			room.watchersMutex.Lock()
			defer room.watchersMutex.Unlock()

			if _, exists := room.watchers[clientID]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
					l.Printf("Client already watching in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already watching in room")
				return
			}

			// Ok the room! now lock the client
			client.watchersMutex.Lock()
			defer client.watchersMutex.Unlock()

			if _, exists := client.watchingRooms[roomId]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/watchRoom]")+" ", 0)
					l.Printf("Client already watching in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already watching in room")
				return
			}

			// Apply the join to both the room and the client
			room.watchers[clientID] = client
			client.watchingRooms[roomId] = room
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/watchRoom]")+" ", 0)
				l.Printf("Client started watching room: %s", roomId)
			}
			JSON(w, http.StatusOK, map[string]string{"status": "watching"})
			return
		}
	}
}

func (e *Engine) unwatchRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/unwatchRoom]")+" ", 0)
			l.Printf("Unauthorized watch room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/unwatchRoom]")+" ", 0)
			l.Printf("Unauthorized unwatch room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		if room, err := e.searchRoom(roomId); err != nil {
			// Room existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/unwatchRoom]")+" ", 0)
				l.Printf("Room not found: %s", roomId)
			}
			Error(w, http.StatusNotFound, "room not found")
			return
		} else if client, err := e.searchClient(clientID); err != nil {
			// Client existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/unwatchRoom]")+" ", 0)
				l.Printf("Client not found: %s", clientID)
			}
			Error(w, http.StatusNotFound, "client not found")
			return
		} else {

			room.watchersMutex.Lock()
			defer room.watchersMutex.Unlock()

			if _, exists := room.watchers[clientID]; !exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/unwatchRoom]")+" ", 0)
					l.Printf("Client not watching in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client not watching in room")
				return
			}

			// Ok the room! now lock the client
			client.watchersMutex.Lock()
			defer client.watchersMutex.Unlock()

			if _, exists := client.watchingRooms[roomId]; !exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/unwatchRoom]")+" ", 0)
					l.Printf("Client not watching in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client not watching in room")
				return
			}

			// Apply the join to both the room and the client
			delete(room.watchers, clientID)
			delete(client.watchingRooms, roomId)
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/unwatchRoom]")+" ", 0)
				l.Printf("Client stopped watching room: %s", roomId)
			}
			JSON(w, http.StatusOK, map[string]string{"status": "not watching"})
			return
		}
	}
}
