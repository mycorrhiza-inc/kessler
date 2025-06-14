// package search

// import (
// 	"encoding/json"
// 	"fmt"
// 	"kessler/internal/objects/networking"
// 	"kessler/internal/quickwit"
// 	"net/http"
// 	"strings"

// 	"github.com/gorilla/mux"
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// type SearchRequest struct {
// 	Index         string                  `json:"index"`
// 	Query         string                  `json:"query"`
// 	SearchFilters networking.FilterFields `json:"filters"`
// 	SortBy        []string                `json:"sort_by"`
// 	MaxHits       int                     `json:"max_hits"`
// 	StartOffset   int                     `json:"start_offset"`
// 	GetText       bool                    `json:"get_text"`
// }

// var ExampleSearchRequest = SearchRequest{
// 	Index:         "cases",
// 	Query:         "test query",
// 	SearchFilters: networking.FilterFields{},
// 	SortBy:        []string{"timestamp"},
// 	MaxHits:       10,
// 	StartOffset:   0,
// 	GetText:       true,
// }

// func (s SearchRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
// 	enc.AddString("index", s.Index)
// 	enc.AddString("query", s.Query)
// 	enc.AddString("filters", s.SearchFilters.String())
// 	enc.AddString("sort_by", strings.Join(s.SortBy, ", "))
// 	enc.AddInt("max_hits", s.MaxHits)
// 	enc.AddInt("start_offset", s.StartOffset)
// 	enc.AddBool("get_text", s.GetText)
// 	return nil
// }

// type RecentUpdatesRequest struct {
// 	Page int `json:"page"`
// }

// func DefineSearchRoutes(search_router *mux.Router) {
// 	search_router.HandleFunc("/file", HandleSearchRequest)
// 	search_router.HandleFunc("/file/recent_updates", HandleRecentUpdatesRequest)
// 	search_router.HandleFunc("/conversation", quickwit.HandleConvoSearch)
// 	search_router.HandleFunc("/organization", quickwit.HandleOrgSearch)
// }

// func HandleSearchRequest(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		fmt.Fprint(w, "Hi there!")
// 		return
// 	case http.MethodPost:
// 		log.Info("Received a search request")

// 		// Create an instance of RequestData
// 		var RequestData SearchRequest

// 		// Decode the JSON body into the struct
// 		err := json.NewDecoder(r.Body).Decode(&RequestData)
// 		if err != nil {
// 			log.Error("Error decoding JSON", zap.Error(err))
// 			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
// 			return
// 		}
// 		defer r.Body.Close() // Close the body when done

// 		pagination := networking.PaginationFromUrlParams(r)
// 		RequestData.MaxHits = int(pagination.Limit)
// 		RequestData.StartOffset = int(pagination.Offset)

// 		hydrated_data, err := SearchQuickwit(RequestData)
// 		if err != nil {
// 			log.Error("Error searching quickwit", zap.Error(err))
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		// q := database.GetTx()
// 		// ctx := r.Context()
// 		// hydrated_data, err := HydrateSearchResults(data, ctx, q)
// 		// if err != nil {
// 		// 	errorstring := fmt.Sprintf("Error hydrating search results: %v\n", err)
// 		// 	log.Info(errorstring)
// 		// 	http.Error(w, errorstring, http.StatusInternalServerError)
// 		// 	return
// 		// }
// 		// TODO : Reneable validation once other stuff is certainly working.
// 		_, err = ValidateHydratedAgainstFilters(hydrated_data, RequestData.SearchFilters)
// 		if err != nil {
// 			log.Error("Returned results did not match filters", zap.Error(err))
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		if len(hydrated_data) == 0 {
// 			log.Error("Quickwit returned 0 search results")
// 		}

// 		respBytes, err := json.Marshal(hydrated_data)
// 		if err != nil {
// 			log.Error("Error marshalling response data", zap.Error(err))
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		respString := string(respBytes)

// 		w.Header().Set("Content-Type", "application/json")

// 		log.Info("Returning search_response", zap.String("body", respString))
// 		fmt.Fprint(w, string(respString))

// 	case http.MethodPut:
// 		fmt.Fprintf(w, "PUT request")
// 	case http.MethodDelete:
// 		fmt.Fprintf(w, "DELETE request")
// 	default:
// 		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
// 	}
// }

// func HandleRecentUpdatesRequest(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		pagination := networking.PaginationFromUrlParams(r)
// 		maxHits := int(pagination.Limit)
// 		offset := int(pagination.Offset)

// 		hydrated_data, err := GetRecentCaseData(maxHits, offset)
// 		if err != nil {
// 			log.Error("Error getting recent case data", zap.Error(err))
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		// q := database.GetTx()
// 		// ctx := r.Context()
// 		// hydrated_data, err := HydrateSearchResults(data, ctx, q)
// 		// if err != nil {
// 		// 	errorstring := fmt.Sprintf("Error hydrating search results: %v\n", err)
// 		// 	log.Info(errorstring)
// 		// 	http.Error(w, errorstring, http.StatusInternalServerError)
// 		// 	return
// 		// }
// 		respString, err := json.Marshal(hydrated_data)
// 		if err != nil {
// 			log.Error("Error marshalling response data", zap.Error(err))
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		fmt.Fprint(w, string(respString))
// 	default:
// 		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)

// 	}
// }
