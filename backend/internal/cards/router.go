package cards

// import (
// 	"encoding/json"
// 	"fmt"
// 	"kessler/internal/dbstore"
// 	"kessler/internal/fugusdk"
// 	"kessler/internal/search"
// 	"kessler/internal/search/filter"
// 	"net/http"
//
// 	"github.com/gorilla/mux"
//
// 	"kessler/internal/cache"
// 	"kessler/internal/dbstore"
// 	"kessler/internal/search"
// )
//
// // RegisterCardLookupRoutes registers endpoints for fetching card data by object UUID.
// func RegisterCardLookupRoutes(r *mux.Router, db dbstore.DBTX) error {
// 	// Fugu search server URL
// 	fuguServerURL := "http://fugudb:3301"
//
// <<<<<<< HEAD
// 	// Organization (Author card)
// 	r.HandleFunc("/orgs/{id}", func(w http.ResponseWriter, req *http.Request) {
// =======
// 	// Initialize filter service (required for search service)
// 	filterService := filter.NewService(fuguServerURL)
//
// 	// Initialize search service
// 	service, err := search.NewSearchService(fuguServerURL, filterService, db)
// 	if err != nil {
// 		return fmt.Errorf("failed to create search service: %w", err)
// 	}
//
// 	// Organization (Author) card endpoint
// 	r.HandleFunc("/org/{id}", func(w http.ResponseWriter, req *http.Request) {
// >>>>>>> main
// 		ctx := req.Context()
// 		id := mux.Vars(req)["id"]
//
// 		// Hydrate organization card
// 		cardData, err := service.HydrateOrganization(ctx, id, 0, 0)
// 		if err != nil {
// 			http.Error(w, "organization not found", http.StatusNotFound)
// 			return
// 		}
// 		respondJSON(w, cardData)
// 	}).Methods(http.MethodGet)
//
// <<<<<<< HEAD
// 		// Build card
// 		extraInfo := ""
// 		if org.IsPerson.Valid && org.IsPerson.Bool {
// 			extraInfo = "Individual contributor"
// 		} else {
// 			extraInfo = "Organization"
// 		}
// 		card := search.AuthorCardData{
// 			Name:        org.Name,
// 			Description: org.Description,
// 			Timestamp:   org.CreatedAt.Time.String(),
// 			ExtraInfo:   extraInfo,
// 			Index:       0,
// 			Type:        "author",
// 			ObjectUUID:  org.ID.String(),
// 		}
//
// 		// Cache result
// 		if cacheCtrl.Client != nil {
// 			if payload, err := json.Marshal(card); err == nil {
// 				cacheCtrl.Set(cacheKey, payload, int32(cache.LongDataTTL))
// 			}
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(card)
// 	}).Methods("GET")
//
// 	// Docket (Conversation card)
// 	r.HandleFunc("/dockets/{id}", func(w http.ResponseWriter, req *http.Request) {
// =======
// 	// Conversation (Docket) card endpoint
// 	r.HandleFunc("/convo/{id}", func(w http.ResponseWriter, req *http.Request) {
// >>>>>>> main
// 		ctx := req.Context()
// 		id := mux.Vars(req)["id"]
//
// 		// Hydrate conversation card
// 		cardData, err := service.HydrateConversation(ctx, id, 0, 0)
// 		if err != nil {
// 			http.Error(w, "conversation not found", http.StatusNotFound)
// 			return
// 		}
// 		respondJSON(w, cardData)
// 	}).Methods(http.MethodGet)
//
// 	// File (Document) card endpoint
// 	r.HandleFunc("/file/{id}", func(w http.ResponseWriter, req *http.Request) {
// 		ctx := req.Context()
// 		id := mux.Vars(req)["id"]
//
// 		// Build minimal search result for hydration
// 		result := fugusdk.FuguSearchResult{ID: id, Text: "", Metadata: nil}
// 		cardData, err := service.HydrateDocument(ctx, result, 0, true)
// 		if err != nil {
// 			http.Error(w, "file not found", http.StatusNotFound)
// 			return
// 		}
// <<<<<<< HEAD
//
// 		// Build card
// 		card := search.DocketCardData{
// 			Name:        conv.Name,
// 			Description: conv.Description,
// 			Timestamp:   conv.CreatedAt.Time.String(),
// 			Index:       0,
// 			Type:        "docket",
// 			ObjectUUID:  conv.ID.String(),
// 		}
//
// 		// Cache result
// 		if cacheCtrl.Client != nil {
// 			if payload, err := json.Marshal(card); err == nil {
// 				cacheCtrl.Set(cacheKey, payload, int32(cache.LongDataTTL))
// 			}
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(card)
// 	}).Methods("GET")
// =======
// 		respondJSON(w, cardData)
// 	}).Methods(http.MethodGet)
// >>>>>>> main
//
// 	return nil
// }
//
// // respondJSON writes the CardData as a JSON response
// func respondJSON(w http.ResponseWriter, data search.CardData) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(data)
// }
//
