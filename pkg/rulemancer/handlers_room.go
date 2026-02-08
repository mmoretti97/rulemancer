/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (e *Engine) roomRoutes(r chi.Router) {
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
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	if room, err := e.newRoom(req.Name, req.Description, req.GameRef); err != nil {
		Error(w, http.StatusInternalServerError, "failed to create room: "+err.Error())
		return
	} else {
		JSON(w, http.StatusCreated, map[string]string{
			"id": room.id,
		})
	}
}

func (e *Engine) apiGetRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		JSON(w, http.StatusOK, map[string]any{
			"id":             room.id,
			"name":           room.name,
			"description":    room.description,
			"clips_instance": room.clipsInstance.Info(),
			"running_game":   room.game.Info(),
			"num_clients":    room.numClients,
		})
	}
}

func (e *Engine) apiDeleteRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := e.removeRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		JSON(w, http.StatusOK, map[string]string{
			"status": "deleted",
		})
	}
}

func (e *Engine) apiListRooms(w http.ResponseWriter, r *http.Request) {
	roomsList := e.listRooms()

	JSON(w, http.StatusOK, map[string]any{
		"rooms": roomsList,
	})
}
