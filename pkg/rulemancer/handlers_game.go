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

func (e *Engine) gameRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Get("/list", e.apiListGames)
		r.Get("/{id}", e.apiGetGame)
	})
}

func (e *Engine) apiGetGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetGame]")+" ", 0)
			l.Printf("Unauthorized get game attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetGame]")+" ", 0)
			l.Printf("Unauthorized get game attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if game, err := e.searchGame(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetGame]")+" ", 0)
			l.Printf("Game not found: %v", err)
		}
		Error(w, http.StatusNotFound, "game not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetGame]")+" ", 0)
			l.Printf("Game %s info provided to client admin", id)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":           game.id,
			"name":         game.name,
			"description":  game.description,
			"rules":        game.rulesLocation,
			"assertable":   game.assertable,
			"responses":    game.responses,
			"queryable":    game.queryable,
			"playingRooms": game.runningRooms,
			"waitingRooms": game.partialRooms,
		})
	}
}

func (e *Engine) apiListGames(w http.ResponseWriter, r *http.Request) {

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListGames]")+" ", 0)
			l.Printf("Unauthorized list games attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if _, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListGames]")+" ", 0)
			l.Printf("Unauthorized list games attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiListGames]")+" ", 0)
		l.Printf("Listing all games")
	}

	gamesList := e.listGames()

	JSON(w, http.StatusOK, map[string]any{
		"games": gamesList,
	})
}
