package autocomplete

import (
	"context"
	"kessler/quickwit"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func DefineAutocompleteRoutes(autocomplete_subrouter *mux.Router) {
	autocomplete_subrouter.HandleFunc(
		"/files-autocomplete",
		AutocompleteFileHandler,
	).Methods(http.MethodGet)
}

func AutocompleteFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
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
		autocomplete_hits := make([]AutoCompleteHit, 0, len(results))
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
		autocomplete_hits := make([]AutoCompleteHit, 0, len(results))
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
	return append(orgResults.Results, convoResults.Results...), nil
}
