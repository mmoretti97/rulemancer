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

func (e *Engine) querySubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/{query}", e.apiQuery)
	})
}

func (e *Engine) apiQuery(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	query := chi.URLParam(r, "query")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {

		requester := ""
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Unauthorized query attempt: %v", err)
			}
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		} else if clientID, ok := claims["id"].(string); !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Unauthorized query attempt with invalid token: %v", claims)
			}
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		} else {
			requester = clientID
		}

		canQuery := false
		room.clientsMutex.RLock()
		room.watchersMutex.RLock()
		if _, ok := room.clients[requester]; ok {
			canQuery = true
		}
		if _, ok := room.watchers[requester]; ok {
			canQuery = true
		}
		room.watchersMutex.RUnlock()
		room.clientsMutex.RUnlock()

		if !canQuery {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Forbidden query attempt in room %s by %s", id, requester)
			}
			Error(w, http.StatusForbidden, "forbidden")
			return
		}

		ci := room.clipsInstance
		if relList, ok := room.game.queryable[query]; !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Query not found for room %s: %s", id, query)
			}
			Error(w, http.StatusNotFound, "query not found")
			return
		} else if len(relList) == 0 {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("No relations for query in room %s: %s", id, query)
			}
			Error(w, http.StatusNotFound, "no relations for query")
			return
		} else {

			// Aggregate all facts from all relations, the loop is split to limit the lock time
			allFacts := make([]string, len(relList))
			ci.Lock()
			for i, rel := range relList {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiQuery]")+" ", 0)
					l.Printf("Processing relation for query in room %s: %s", id, rel)
				}

				if factList, err := room.clipsInstance.QueryFactsAtomic(rel); err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
						l.Printf("Error querying status in room %s - %s: %v", id, rel, err)
					}
					ci.Unlock()
					Error(w, http.StatusInternalServerError, "failed to query status")
					return
				} else {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiQuery]")+" ", 0)
						l.Printf("Status in room %s - %s: %+v", id, rel, factList)
					}
					allFacts[i] = factList
				}
			}
			ci.Unlock()

			response := make(map[string][]map[string]string)
			for i, factList := range allFacts {

				if factMap, err := genericFactToMap(e.Config, relList[i], factList); err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
						l.Printf("Error converting fact to struct in room %s - %s: %v", id, relList[i], err)
					}
					Error(w, http.StatusInternalServerError, "failed to convert fact to struct")
					return
				} else {
					response[relList[i]] = factMap
				}
			}

			JSON(w, http.StatusOK, map[string]any{
				"response": response,
			})
		}
	}
}
