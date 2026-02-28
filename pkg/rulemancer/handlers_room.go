/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) roomRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Post("/create", e.apiCreateRoom)
		r.Get("/list", e.apiListRooms)
		r.Get("/{id}", e.apiGetRoom)
		r.Delete("/{id}", e.apiDeleteRoom)
		r.Route("/{id}/", e.roomSubRoutes)
	})
}

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GameRef     string `json:"game_ref"`
}

func (e *Engine) apiCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiCreateRoom]")+" ", 0)
			l.Printf("Invalid JSON: %v", err)
		}
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	if room, err := e.newRoom(req.Name, req.Description, req.GameRef); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiCreateRoom]")+" ", 0)
			l.Printf("Failed to create room: %v", err)
		}
		Error(w, http.StatusInternalServerError, "failed to create room: "+err.Error())
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiCreateRoom]")+" ", 0)
			l.Printf("Room created: %v", room)
		}
		JSON(w, http.StatusCreated, map[string]string{
			"id": room.id,
		})
	}
}

func (e *Engine) apiGetRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetRoom]")+" ", 0)
			l.Printf("Unauthorized get room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetRoom]")+" ", 0)
			l.Printf("Unauthorized get room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if room, err := e.searchRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetRoom]")+" ", 0)
			l.Printf("Room not found: %v", err)
		}
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetRoom]")+" ", 0)
			l.Printf("Room info provided to client: %v", room)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":                room.id,
			"name":              room.name,
			"description":       room.description,
			"clips_instance":    room.clipsInstance.Info(),
			"running_game":      room.game.name,
			"num_clients":       room.maxClients,
			"playing_clients":   room.clients,
			"watching_clients":  room.watchers,
			"connected_sockets": room.socketsInfo(),
		})
	}
}

func (e *Engine) apiDeleteRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteRoom]")+" ", 0)
			l.Printf("Unauthorized delete room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteRoom]")+" ", 0)
			l.Printf("Unauthorized delete room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if _, err := e.removeRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteRoom]")+" ", 0)
			l.Printf("Room not found: %v", err)
		}
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiDeleteRoom]")+" ", 0)
			l.Printf("Room deleted: %v", id)
		}
		JSON(w, http.StatusOK, map[string]string{
			"status": "deleted",
		})
	}
}

func (e *Engine) apiListRooms(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListRooms]")+" ", 0)
			l.Printf("Unauthorized list rooms attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListRooms]")+" ", 0)
			l.Printf("Unauthorized list rooms attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiListRooms]")+" ", 0)
		l.Printf("Rooms list provided to client: %v", e.listRooms())
	}

	roomsList := e.listRooms()

	JSON(w, http.StatusOK, map[string]any{
		"rooms": roomsList,
	})
}
