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

func (e *Engine) roomSubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/assert/{assertion}", e.apiAssert)
		r.Route("/query", e.querySubRoutes)

		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomSubRoutes]")+" ", 0)
			l.Printf("Debug mode enabled: adding /facts endpoints")
			r.Get("/facts", e.apiGetFacts)
		}
	})
}

func (e *Engine) apiAssert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	assertion := chi.URLParam(r, "assertion")

	if room, err := e.searchRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
			l.Printf("Room not found: %s", id)
		}
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {

		requester := ""
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
				l.Printf("Unauthorized assert attempt: %v", err)
			}
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		} else if clientID, ok := claims["id"].(string); !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
				l.Printf("Unauthorized assert attempt with invalid token: %v", claims)
			}
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		} else {
			requester = clientID
		}

		canAssert := false
		room.clientsMutex.RLock()
		if _, ok := room.clients[requester]; ok {
			canAssert = true
		}
		room.clientsMutex.RUnlock()

		if !canAssert {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
				l.Printf("Forbidden assert attempt in room %s by %s", id, requester)
			}
			Error(w, http.StatusForbidden, "forbidden")
			return
		}

		if relList, ok := room.game.assertable[assertion]; !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
				l.Printf("Assertion not found for room %s: %s", id, assertion)
			}
			Error(w, http.StatusNotFound, "assertion not found")
			return
		} else {

			// Read raw JSON body into a map
			var raw map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Error decoding JSON body for assertion in room %s: %v", id, err)
				}
				Error(w, http.StatusBadRequest, "invalid JSON body")
				return
			}

			// Create the facts list
			facts := make([]string, 0)

			for _, rel := range relList {
				if _, exists := raw[rel]; !exists {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
						l.Printf("Missing required field for assertion in room %s: %s", id, rel)
					}
					Error(w, http.StatusBadRequest, "missing required field: "+rel)
					return
				} else {
					if newFacts, err := jsonGenericDecoder(e.Config, raw[rel]); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
							l.Printf("Error decoding field for assertion in room %s - %s: %v", id, rel, err)
						}
						Error(w, http.StatusBadRequest, "invalid field format: "+rel)
						return
					} else {
						// Append each fact wrapped in the relation
						for _, fact := range newFacts {
							fact := "(" + rel + " " + fact + ")"
							facts = append(facts, fact)
						}
					}
				}
			}

			ci := room.clipsInstance
			ci.Lock()

			for _, fact := range facts {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Asserting fact in room %s: %s", id, fact)
				}
				if err := ci.AssertFactAtomic(fact); err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
						l.Printf("Error asserting fact in room %s - %s: %v", id, fact, err)
					}
					ci.Unlock()
					Error(w, http.StatusInternalServerError, "failed to assert")
					return
				}
			}

			if err := ci.RunAtomic(); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Error running CLIPS in room %s: %v", id, err)
				}
				ci.Unlock()
				Error(w, http.StatusInternalServerError, "failed to run")
				return
			} else {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Successfully ran CLIPS in room %s", id)
				}
			}

			// Prepare the response
			response := make(map[string][]map[string]string)

			if relList, ok := room.game.responses[assertion]; !ok {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Assertion has no response relations in room %s: %s", id, assertion)
				}
			} else if len(relList) == 0 {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("No relations for assertion in room %s: %s", id, assertion)
				}
			} else {

				// Aggregate all facts from all relations, the loop is split to limit the lock time
				allFacts := make([]string, len(relList))
				for i, rel := range relList {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
						l.Printf("Processing relation for assertion in room %s: %s", id, rel)
					}

					if factList, err := room.clipsInstance.QueryFactsAtomic(rel); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
							l.Printf("Error querying status in room %s - %s: %v", id, rel, err)
						}
						ci.Unlock()
						Error(w, http.StatusInternalServerError, "failed to query status")
						return
					} else {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
							l.Printf("Status in room %s - %s: %+v", id, rel, factList)
						}
						allFacts[i] = factList
					}
				}

				for i, factList := range allFacts {

					if factMap, err := genericFactToMap(e.Config, relList[i], factList); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
							l.Printf("Error converting fact to struct in room %s - %s: %v", id, relList[i], err)
						}
						ci.Unlock()
						Error(w, http.StatusInternalServerError, "failed to convert fact to struct")
						return
					} else {
						response[relList[i]] = factMap
					}
				}
			}

			ci.Unlock()

			JSON(w, http.StatusOK, map[string]any{
				"status":   "asserted",
				"response": response,
			})
		}
	}
}

func (e *Engine) apiGetFacts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetFacts]")+" ", 0)
			l.Printf("Unauthorized get facts attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetFacts]")+" ", 0)
			l.Printf("Unauthorized get facts attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if room, err := e.searchRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetFacts]")+" ", 0)
			l.Printf("Room not found: %v", err)
		}
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		facts, err := room.clipsInstance.QueryFactsAllFacts()
		if err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetFacts]")+" ", 0)
				l.Printf("Failed to get facts in room %s: %v", id, err)
			}
			Error(w, http.StatusInternalServerError, "failed to get facts")
			return
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetFacts]")+" ", 0)
			l.Printf("Facts in room %s: %+v", id, facts)
		}
		JSON(w, http.StatusOK, map[string]any{
			"facts": facts,
		})
	}
}
