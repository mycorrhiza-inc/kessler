package cards

import (
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/internal/search"
	"kessler/internal/search/filter"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterCardLookupRoutes registers endpoints for fetching card data by object UUID.
func RegisterCardLookupRoutes(r *mux.Router, db dbstore.DBTX) error {
	// Fugu search server URL
	fuguServerURL := "http://fugudb:3301"

	// Initialize filter service (required for search service)
	filterService := filter.NewService(fuguServerURL)

	// Initialize search service
	service, err := search.NewSearchService(fuguServerURL, filterService, db)
	if err != nil {
		return fmt.Errorf("failed to create search service: %w", err)
	}

	// Organization (Author) card endpoint
	r.HandleFunc("/org/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := mux.Vars(req)["id"]

		// Hydrate organization card
		cardData, err := service.HydrateOrganization(ctx, id, 0, 0)
		if err != nil {
			http.Error(w, "organization not found", http.StatusNotFound)
			return
		}
		respondJSON(w, cardData)
	}).Methods(http.MethodGet)

	// Conversation (Docket) card endpoint
	r.HandleFunc("/convo/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := mux.Vars(req)["id"]

		// Hydrate conversation card
		cardData, err := service.HydrateConversation(ctx, id, 0, 0)
		if err != nil {
			http.Error(w, "conversation not found", http.StatusNotFound)
			return
		}
		respondJSON(w, cardData)
	}).Methods(http.MethodGet)

	// File (Document) card endpoint
	r.HandleFunc("/file/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := mux.Vars(req)["id"]

		// Build minimal search result for hydration
		result := fugusdk.FuguSearchResult{ID: id, Text: "", Metadata: nil}
		cardData, err := service.HydrateDocument(ctx, result, 0, true)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		respondJSON(w, cardData)
	}).Methods(http.MethodGet)

	return nil
}

// respondJSON writes the CardData as a JSON response
func respondJSON(w http.ResponseWriter, data search.CardData) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

