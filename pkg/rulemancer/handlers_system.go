/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) systemRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Get("/health", e.health)
		r.Post("/quit", e.quit)
	})
}

func (e *Engine) health(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/health]")+" ", 0)
			l.Printf("Unauthorized health check attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/health]")+" ", 0)
			l.Printf("Unauthorized health check attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/health]")+" ", 0)
		l.Printf("Health check successful for client ID: %s", claims["id"])
	}
	e.gamesMutex.RLock()
	numGames := e.numGames
	e.gamesMutex.RUnlock()

	e.roomsMutex.RLock()
	numRooms := e.numRooms
	e.roomsMutex.RUnlock()

	e.clientsMutex.RLock()
	numClients := e.numClients
	e.clientsMutex.RUnlock()
	JSON(w, http.StatusOK, map[string]string{"status": "OK",
		"served games":      strconv.Itoa(numGames),
		"running rooms":     strconv.Itoa(numRooms),
		"connected clients": strconv.Itoa(numClients)})
}

func (e *Engine) quit(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/quit]")+" ", 0)
			l.Printf("Unauthorized quit attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/quit]")+" ", 0)
			l.Printf("Unauthorized quit attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/quit]")+" ", 0)
		l.Printf("Shutdown initiated by client ID: %s", claims["id"])
	}
	e.stopChan <- syscall.SIGTERM
	JSON(w, http.StatusOK, map[string]string{"status": "shutting down"})
}
