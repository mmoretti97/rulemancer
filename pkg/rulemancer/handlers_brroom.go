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

func (e *Engine) brRoomRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Post("/create", e.apiCreateBrRoom)
		r.Get("/list", e.apiListBrRooms)
		r.Get("/{id}", e.apiGetBrRoom)
		r.Delete("/{id}", e.apiDeleteBrRoom)
		r.Route("/{id}/", e.brRoomSubRoutes)
	})
}

type CreateBrRoomRequest struct {
	Name      string `json:"name"`
	BridgeRef string `json:"bridge_ref"`
}

func (e *Engine) apiCreateBrRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateBrRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiCreateBrRoom]")+" ", 0)
			l.Printf("Invalid JSON: %v", err)
		}
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	if brRoom, err := e.newBrRoom(req.Name, req.BridgeRef); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiCreateBrRoom]")+" ", 0)
			l.Printf("Failed to create bridge room: %v", err)
		}
		Error(w, http.StatusInternalServerError, "failed to create bridge room: "+err.Error())
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiCreateBrRoom]")+" ", 0)
			l.Printf("Bridge room created: %v", brRoom)
		}
		JSON(w, http.StatusCreated, map[string]string{
			"id": brRoom.id,
		})
	}
}

func (e *Engine) apiGetBrRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetBrRoom]")+" ", 0)
			l.Printf("Unauthorized get room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetBrRoom]")+" ", 0)
			l.Printf("Unauthorized get room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if brRoom, err := e.searchBrRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetBrRoom]")+" ", 0)
			l.Printf("Bridge room not found: %v", err)
		}
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetBrRoom]")+" ", 0)
			l.Printf("Bridge room info provided to client: %v", brRoom)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":             brRoom.id,
			"clips_instance": brRoom.clipsInstance.Info(),
		})
	}
}

func (e *Engine) apiDeleteBrRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteBrRoom]")+" ", 0)
			l.Printf("Unauthorized delete bridge room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteBrRoom]")+" ", 0)
			l.Printf("Unauthorized delete bridge room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if _, err := e.removeBrRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteBrRoom]")+" ", 0)
			l.Printf("Bridge room not found: %v", err)
		}
		Error(w, http.StatusNotFound, "bridge room not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiDeleteBrRoom]")+" ", 0)
			l.Printf("Bridge room deleted: %v", id)
		}
		JSON(w, http.StatusOK, map[string]string{
			"status": "deleted",
		})
	}
}

func (e *Engine) apiListBrRooms(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListBrRooms]")+" ", 0)
			l.Printf("Unauthorized list bridge rooms attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListBrRooms]")+" ", 0)
			l.Printf("Unauthorized list bridge rooms attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiListBrRooms]")+" ", 0)
		l.Printf("Bridge rooms list provided to client: %v", e.listBrRooms())
	}

	brRoomsList := e.listBrRooms()

	JSON(w, http.StatusOK, map[string]any{
		"brrooms": brRoomsList,
	})
}
