/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) clientRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(e.JWTAuth))
		r.Use(jwtauth.Authenticator(e.JWTAuth))
		r.Get("/list", e.apiListClients)
		r.Get("/{id}", e.apiGetClient)
		r.Get("/current", e.apiGetCurrentClient)
		r.Delete("/{id}", e.apiDeleteClient)
	})
}

func (e *Engine) apiGetClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, claims, err := jwtauth.FromContext(r.Context())
	requester := ""
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Unauthorized get clientattempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || (clientID != "admin" && clientID != id) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Unauthorized get client attempt by %s with invalid token: %v", requester, claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		requester = clientID
	}

	if client, err := e.searchClient(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Client not found: %v", err)
		}
		Error(w, http.StatusNotFound, "client not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Client requested by %s found: %s", requester, client.id)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":          client.id,
			"name":        client.name,
			"description": client.description,
		})
	}
}

func (e *Engine) apiGetCurrentClient(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	requester := ""
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetCurrentClient]")+" ", 0)
			l.Printf("Unauthorized get currentclient attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || (clientID == "admin") {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetCurrentClient]")+" ", 0)
			l.Printf("Unauthorized get current client attempt by %s with invalid token: %v", requester, claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		requester = clientID
	}

	if client, err := e.searchClient(requester); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetCurrentClient]")+" ", 0)
			l.Printf("Client not found: %v", err)
		}
		Error(w, http.StatusNotFound, "client not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetCurrentClient]")+" ", 0)
			l.Printf("Client requested by %s found: %s", requester, client.id)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":          client.id,
			"name":        client.name,
			"description": client.description,
		})
	}
}

func (e *Engine) apiDeleteClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteClient]")+" ", 0)
			l.Printf("Unauthorized delete client attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteClient]")+" ", 0)
			l.Printf("Unauthorized delete client attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiDeleteClient]")+" ", 0)
		l.Printf("Client deletion initiated by client ID: %s", claims["id"])
	}

	if _, err := e.removeClient(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiDeleteClient]")+" ", 0)
			l.Printf("Client not found: %v", err)
		}
		Error(w, http.StatusNotFound, "client not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiDeleteClient]")+" ", 0)
			l.Printf("Client deleted: %s", id)
		}
		JSON(w, http.StatusOK, map[string]string{
			"status": "deleted",
		})
	}
}

func (e *Engine) apiListClients(w http.ResponseWriter, r *http.Request) {

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListClients]")+" ", 0)
			l.Printf("Unauthorized list clients attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListClients]")+" ", 0)
			l.Printf("Unauthorized list clients attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiListClients]")+" ", 0)
		l.Printf("List clients requested by client ID: %s", claims["id"])
	}

	clientsList := e.listClients()
	JSON(w, http.StatusOK, map[string]any{
		"clients": clientsList,
	})
}
