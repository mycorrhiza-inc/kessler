package cards

import (
	"encoding/json"
	"net/http"

	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/search"

	"github.com/gorilla/mux"
)

// RegisterCardLookupRoutes registers endpoints for fetching card data by object UUID.
func RegisterCardLookupRoutes(r *mux.Router, db dbstore.DBTX) error {
	// Initialize cache controller (can continue without cache if unavailable)
	cacheCtrl, _ := cache.NewCacheController()

	// Organization (Author card)
	r.HandleFunc("/org/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := mux.Vars(req)["id"]
		cacheKey := cache.PrepareKey("search", "organization", id)

		// Try cache first
		if cacheCtrl.Client != nil {
			if data, err := cacheCtrl.Get(cacheKey); err == nil {
				var cached search.AuthorCardData
				if err := json.Unmarshal(data, &cached); err == nil && cached.Type == "author" {
					cached.Index = 0
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(cached)
					return
				}
			}
		}

		// Build card data
		card, err := search.BuildAuthorCard(ctx, db, id, 0)
		if err != nil {
			http.Error(w, "organization not found", http.StatusNotFound)
			return
		}

		// Cache result
		if cacheCtrl.Client != nil {
			if payload, err := json.Marshal(card); err == nil {
				cacheCtrl.Set(cacheKey, payload, int32(cache.LongDataTTL))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}).Methods("GET")

	// Docket (Conversation card)
	r.HandleFunc("/convo/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := mux.Vars(req)["id"]
		cacheKey := cache.PrepareKey("search", "conversation", id)

		// Try cache first
		if cacheCtrl.Client != nil {
			if data, err := cacheCtrl.Get(cacheKey); err == nil {
				var cached search.DocketCardData
				if err := json.Unmarshal(data, &cached); err == nil && cached.Type == "docket" {
					cached.Index = 0
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(cached)
					return
				}
			}
		}

		// Build card data
		card, err := search.BuildDocketCard(ctx, db, id, 0)
		if err != nil {
			http.Error(w, "conversation not found", http.StatusNotFound)
			return
		}

		// Cache result
		if cacheCtrl.Client != nil {
			if payload, err := json.Marshal(card); err == nil {
				cacheCtrl.Set(cacheKey, payload, int32(cache.LongDataTTL))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}).Methods("GET")

	
	// File (Document card)
	r.HandleFunc("/file/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := mux.Vars(req)["id"]
		// Build raw document card (full hydration)
		// Create a minimal FuguSearchResult wrapper
		res := search.FuguResultWrapper{ID: id, Text: "", Metadata: nil}
		card, err := search.BuildDocumentCard(ctx, db, res, 0, true)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}).Methods("GET")

	return nil
}