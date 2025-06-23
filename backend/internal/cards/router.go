package cards

import (
	"encoding/json"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/search"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// RegisterCardLookupRoutes registers endpoints for fetching card data by object UUID.
func RegisterCardLookupRoutes(r *mux.Router, db dbstore.DBTX) error {
	// Initialize cache controller (can continue without cache if unavailable)
	cacheCtrl, _ := cache.NewCacheController()

	// Organization (Author card)
	r.HandleFunc("/org/{id}", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		id := vars["id"]

		// Try cache first
		cacheKey := cache.PrepareKey("search", "organization", id)
		if cacheCtrl.Client != nil {
			if data, err := cacheCtrl.Get(cacheKey); err == nil {
				var cached search.AuthorCardData
				if err := json.Unmarshal(data, &cached); err == nil && cached.Type == "author" {
					// Return cached card
					cached.Index = 0
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(cached)
					return
				}
			}
		}

		// Parse and fetch from DB
		orgID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "invalid UUID", http.StatusBadRequest)
			return
		}
		queries := dbstore.New(db)
		org, err := queries.OrganizationRead(ctx, orgID)
		if err != nil {
			http.Error(w, "organization not found", http.StatusNotFound)
			return
		}

		// Build card
		extraInfo := ""
		if org.IsPerson.Valid && org.IsPerson.Bool {
			extraInfo = "Individual contributor"
		} else {
			extraInfo = "Organization"
		}
		card := search.AuthorCardData{
			Name:        org.Name,
			Description: org.Description,
			Timestamp:   org.CreatedAt.Time,
			ExtraInfo:   extraInfo,
			Index:       0,
			Type:        "author",
			ObjectUUID:  org.ID,
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
		vars := mux.Vars(req)
		id := vars["id"]

		// Try cache first
		cacheKey := cache.PrepareKey("search", "conversation", id)
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

		// Parse UUID and fetch
		convID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "invalid UUID", http.StatusBadRequest)
			return
		}
		queries := dbstore.New(db)
		conv, err := queries.DocketConversationRead(ctx, convID)
		if err != nil {
			http.Error(w, "docket not found", http.StatusNotFound)
			return
		}

		// Build card
		card := search.DocketCardData{
			Name:        conv.Name,
			Description: conv.Description,
			Timestamp:   conv.CreatedAt.Time,
			Index:       0,
			Type:        "docket",
			ObjectUUID:  conv.ID,
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

	return nil
}
