package autocomplete

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/quickwit"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func DefineAutocompleteRoutes(autocomplete_subrouter *mux.Router) {
	autocomplete_subrouter.HandleFunc(
		"/files-basic",
		AutocompleteFileHandler,
	).Methods(http.MethodGet)
}

func AutocompleteFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("query")
	autocomplete_hits, err := AutoCompleteFileGetResults(query, ctx)
	if err != nil {
		log.Error("Error getting autocomplete results", "err", err)
		http.Error(w, fmt.Sprintf("Error getting autocomplete results: %v", err), http.StatusInternalServerError)
		return
	}
	return_bytes, err := json.Marshal(autocomplete_hits)
	if err != nil {
		log.Error("Error marshaling autocomplete results", "err", err)
		http.Error(w, fmt.Sprintf("Error marshaling autocomplete results: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(return_bytes)
}

type AutoCompleteHit struct {
	ID   uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

func AutoCompleteFileGetResults(query string, ctx context.Context) ([]AutoCompleteHit, error) {
	results_each := 10
	type AsyncResult struct {
		Results []AutoCompleteHit
		Err     error
	}

	convoChan := make(chan AsyncResult, 1)
	orgChan := make(chan AsyncResult, 1)
	go func() {
		search_convo_results := quickwit.ConvoSearchRequestData{
			Search: quickwit.ConversationSearchSchema{
				Query: query,
			},
			Limit:  results_each,
			Offset: 0,
		}
		results, err := quickwit.SearchConversations(search_convo_results, ctx)
		if err != nil {
			log.Error("Encountered Error while getting quickwit autocomplete", "err", err)
			convoChan <- AsyncResult{Err: err}
			return
		}
		log.Info("Creating Autocomplete Hits")
		autocomplete_hits := make([]AutoCompleteHit, len(results))
		for index, result := range results {
			autocomplete_hits[index] = AutoCompleteHit{
				ID:   result.ID,
				Name: result.Name,
				Type: "conversation",
			}
		}
		convoChan <- AsyncResult{Results: autocomplete_hits}
		return
	}()

	go func() {
		search_convo_results := quickwit.OrgSearchRequestData{
			Search: quickwit.OrganizationSearchSchema{
				Query: query,
			},
			Limit:  results_each,
			Offset: 0,
		}
		results, err := quickwit.SearchOrganizations(search_convo_results, ctx)
		if err != nil {
			log.Error("Encountered Error while getting quickwit autocomplete", "err", err)
			orgChan <- AsyncResult{Err: err}
			return
		}
		autocomplete_hits := make([]AutoCompleteHit, len(results))
		log.Info("Creating Autocomplete Hits")
		for index, result := range results {
			autocomplete_hits[index] = AutoCompleteHit{
				ID:   result.ID,
				Name: result.Name,
				Type: "organization",
			}
		}
		orgChan <- AsyncResult{Results: autocomplete_hits}
		return
	}()
	orgResults := <-orgChan
	if orgResults.Err != nil {
		return []AutoCompleteHit{}, nil
	}
	convoResults := <-convoChan
	if convoResults.Err != nil {
		return []AutoCompleteHit{}, nil
	}
	log.Info("Got autocomplete results!")
	return append(orgResults.Results, convoResults.Results...), nil
}
