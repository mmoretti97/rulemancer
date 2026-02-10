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
)

func (e *Engine) newRoutes(r chi.Router) {
	// For reference: Verifier and Authenticator are usually here. For the creation of new entities it is not the case
	// because the creation has to be possible without pre-existing token
	r.Route("/", func(r chi.Router) {
		r.Post("/client", e.apiCreateClient) // Client creation
	})
}

type CreateClientRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (e *Engine) apiCreateClient(w http.ResponseWriter, r *http.Request) {
	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiCreateClient]")+" ", 0)
			l.Printf("Invalid JSON: %v", err)
		}
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	client := e.newClient(req.Name, req.Description)

	_, tokenString, _ := e.Encode(map[string]interface{}{"id": client.id})

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiCreateClient]")+" ", 0)
		l.Printf("Creating client: %s with ID: %s", req.Name, client.id)
	}

	JSON(w, http.StatusCreated, map[string]string{
		"id":        client.id,
		"api_token": tokenString,
	})
}
