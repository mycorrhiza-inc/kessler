package autocomplete

import (
	"net/http"

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

func AutoCompleteFileGetResults(query string) []AutoCompleteHit {
}
