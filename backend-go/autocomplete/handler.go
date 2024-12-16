package autocomplete

import (
	"net/http"

	"github.com/gorilla/mux"
)

func DefineAutocompleteRoutes(autocomplete_subrouter *mux.Router) {
	autocomplete_subrouter.HandleFunc(
		"/conversation-autocomplete",
		AutocompleteConversationHandler,
	).Methods(http.MethodGet)
}

func AutocompleteConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
