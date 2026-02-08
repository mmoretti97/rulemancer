/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"net/http"

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
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func (e *Engine) joinRoom(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func (e *Engine) newGameRoom(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})

}

// TODO
