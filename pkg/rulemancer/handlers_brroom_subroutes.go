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

func (e *Engine) brRoomSubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/request", e.apiBridgeRequest)

		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/brRoomSubRoutes]")+" ", 0)
			l.Printf("Debug mode enabled: adding /facts endpoints")
			r.Get("/facts", e.apiGetFacts)
		}
	})
}

func (e *Engine) apiBridgeRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if brRoom, err := e.searchBrRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
			l.Printf("Bridge room not found: %s", id)
		}
		Error(w, http.StatusNotFound, "bridge room not found")
		return
	} else {

		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
				l.Printf("Unauthorized request attempt: %v", err)
			}
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		} else if _, ok := claims["id"].(string); !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
				l.Printf("Unauthorized request attempt with invalid token: %v", claims)
			}
			Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Read raw JSON body into a map
		var raw map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
				l.Printf("Error decoding JSON body for assertion in room %s: %v", id, err)
			}
			Error(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		ci := brRoom.clipsInstance
		ci.Lock()

		// The request is composed of multiple facts asserted together, the "facts" field specifies a list of
		// relations, each relation is a map of variable names to values, the values can be a single value or a
		// list of values, for example:

		// "facts": [
		//         {"move": {"player": ["Alice"], "from": ["A2","to","A3"]}},
		//         {"move": {"player": ["Bob"], "from": ["B2","to","B3"]}}
		//]

		// The request can also specify relations to be included in the response, for example:

		// "response": ["player_status", "game_status"]

		// Create the facts list
		facts := make([]string, 0)

		if factsListRaw, ok := raw["facts"]; !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiBridgeRequest]")+" ", 0)
				l.Printf("No facts field in request body for assertion in room %s", id)
			}
		} else {

			var factList []map[string]json.RawMessage

			if err := json.Unmarshal(factsListRaw, &factList); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
					l.Printf("Error decoding facts list for assertion in room %s: %v", id, err)
				}
				Error(w, http.StatusBadRequest, "invalid facts format")
				return

			} else {

				for _, factRaw := range factList {
					for rel, factProp := range factRaw {

						if newFacts, err := jsonGenericDecoder(e.Config, factProp); err != nil {
							if e.Debug {
								l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
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
			}

			for _, fact := range facts {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiBridgeRequest]")+" ", 0)
					l.Printf("Asserting fact in room %s: %s", id, fact)
				}
				if err := ci.AssertFactAtomic(fact); err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
						l.Printf("Error asserting fact in room %s - %s: %v", id, fact, err)
					}
					ci.Unlock()
					Error(w, http.StatusInternalServerError, "failed to assert")
					return
				} else {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiBridgeRequest]")+" ", 0)
						l.Printf("Successfully asserted fact in room %s: %s", id, fact)
					}
				}
			}

			if err := ci.RunAtomic(); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
					l.Printf("Error running CLIPS in room %s: %v", id, err)
				}
				ci.Unlock()
				Error(w, http.StatusInternalServerError, "failed to run")
				return
			} else {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiBridgeRequest]")+" ", 0)
					l.Printf("Successfully ran CLIPS in room %s", id)
				}
			}
		}

		// Prepare the response
		response := make(map[string][]map[string]string)

		if queries, ok := raw["queries"]; !ok {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiBridgeRequest]")+" ", 0)
				l.Printf("No queries field in request body for assertion in room %s", id)
			}
		} else {

			var queryList []string
			if err := json.Unmarshal(queries, &queryList); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
					l.Printf("Error decoding queries list for assertion in room %s: %v", id, err)
				}
				Error(w, http.StatusBadRequest, "invalid queries format")
				return
			} else {
				// The response is a list of relations to query after the run
				allFacts := make([]string, len(queryList))

				for i, rel := range queryList {

					if factList, err := ci.QueryFactsAtomic(rel); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
							l.Printf("Error querying status in room %s - %s: %v", id, rel, err)
						}
						ci.Unlock()
						Error(w, http.StatusInternalServerError, "failed to query status")
						return
					} else {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiBridgeRequest]")+" ", 0)
							l.Printf("Status in room %s - %s: %+v", id, rel, factList)
						}
						allFacts[i] = factList
					}
				}

				for i, factList := range allFacts {

					if factMap, err := genericFactToMap(e.Config, queryList[i], factList); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiBridgeRequest]")+" ", 0)
							l.Printf("Error converting fact to struct in room %s - %s: %v", id, queryList[i], err)
						}
						ci.Unlock()
						Error(w, http.StatusInternalServerError, "failed to convert fact to struct")
						return
					} else {
						response[queryList[i]] = factMap
					}
				}
			}
		}

		ci.Unlock()

		JSON(w, http.StatusOK, map[string]any{
			"asserted": facts,
			"response": response,
		})

	}
}
