// src: internal/search/handler.go
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/internal/search/filter"
	"kessler/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("search-service")

// RegisterSearchRoutes registers all search-related routes including filter configuration
func RegisterSearchRoutes(router *mux.Router, db dbstore.DBTX) error {
	fuguServerURL := "http://fugudb:3301"

	// Create filter service and handler
	filterService := filter.NewService(fuguServerURL)
	filterHandler := filter.NewHandler(filterService)

	// Create search service and handler with database
	service, err := NewSearchService(fuguServerURL, filterService, db)
	if err != nil {
		return fmt.Errorf("failed to create search service: %w", err)
	}

	handler := NewSearchHandler(service)

	fmt.Println("🔧 Registering search routes...")

	// Main search endpoints
	router.HandleFunc("/", handler.SearchTextGet).Methods(http.MethodGet)
	router.HandleFunc("/", handler.Search).Methods(http.MethodPost)

	// Namespace-specific search endpoints
	router.HandleFunc("/conversations", handler.SearchConversations).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/organizations", handler.SearchOrganizations).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/all", handler.SearchAll).Methods(http.MethodGet, http.MethodPost)

	// Search info and health
	router.HandleFunc("/info", handler.GetSearchInfo).Methods(http.MethodGet)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	// Filter endpoints - use the dedicated filter handler
	router.HandleFunc("/filters", handler.GetAvailableFilters).Methods(http.MethodGet)
	router.HandleFunc("/filters/all", filterHandler.GetAllFilters).Methods(http.MethodGet)
	router.HandleFunc("/filters/namespace/{namespace}", filterHandler.GetNamespaceFilters).Methods(http.MethodGet)
	router.HandleFunc("/filters/path/{filterPath:.*}", filterHandler.GetFilterValues).Methods(http.MethodGet)
	router.HandleFunc("/filters/invalidate", filterHandler.InvalidateCache).Methods(http.MethodPost)

	fmt.Println("✅ Search and filter routes registered successfully")
	return nil
}

// MarshalJSON for SearchResponse to handle the interface type
func (sr SearchResponse) MarshalJSON() ([]byte, error) {
	// Create a temporary struct with the same fields but data as []interface{}
	type Alias SearchResponse
	return json.Marshal(&struct {
		Data []interface{} `json:"data"`
		*Alias
	}{
		Data:  convertCardsToInterface(sr.Data),
		Alias: (*Alias)(&sr),
	})
}

// Helper function to convert []CardData to []interface{}
func convertCardsToInterface(cards []CardData) []interface{} {
	result := make([]interface{}, len(cards))
	for i, card := range cards {
		result[i] = card
	}
	return result
}

// SearchServiceHandler handles HTTP requests for search
type SearchServiceHandler struct {
	service *SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *SearchService) *SearchServiceHandler {
	return &SearchServiceHandler{
		service: service,
	}
}

// GetSearchInfo handles GET requests to /search/info endpoint
func (h *SearchServiceHandler) GetSearchInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:get-search-info")
	defer span.End()

	logger.Info(ctx, "search info request received")

	// Get search statistics and capabilities
	info, err := h.service.GetSearchInfo(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get search info", zap.Error(err))
		http.Error(w, "Failed to get search information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		logger.Error(ctx, "failed to encode search info response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "search info served successfully")
}

// Search handles POST search requests
func (h *SearchServiceHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:search")
	defer span.End()

	logger.Info(ctx, "POST search request received")

	// Parse request body
	var searchReq SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		logger.Error(ctx, "failed to decode search request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if searchReq.Query == "" {
		logger.Error(ctx, "empty search query provided")
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	// Set defaults for pagination
	if searchReq.Page < 0 {
		searchReq.Page = 0
	}
	if searchReq.PerPage <= 0 {
		searchReq.PerPage = 20
	}

	pagination := PaginationParams{
		Page:  searchReq.Page,
		Limit: searchReq.PerPage,
	}

	logger.Info(ctx, "processing POST search request",
		zap.String("query", searchReq.Query),
		zap.String("namespace", searchReq.Namespace),
		zap.Int("page", pagination.Page),
		zap.Int("limit", pagination.Limit),
		zap.Int("filter_count", len(searchReq.Filters)))

	// Process the search
	response, err := h.service.ProcessSearch(ctx, searchReq.Query, searchReq.Filters, pagination, searchReq.Namespace)
	if err != nil {
		logger.Error(ctx, "search processing failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode search response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "POST search completed successfully",
		zap.Int("result_count", len(response.Data)))
}

// SearchTextGet handles GET search requests with query parameters
func (h *SearchServiceHandler) SearchTextGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:search-text-get")
	defer span.End()

	logger.Info(ctx, "GET search request received")

	// Extract query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		logger.Error(ctx, "empty search query provided")
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Extract namespace
	namespace := r.URL.Query().Get("namespace")

	// Extract pagination from query parameters
	pagination := h.extractPagination(r)

	// Extract filters from query parameters
	filters := h.extractFilters(r)

	logger.Info(ctx, "processing GET search request",
		zap.String("query", query),
		zap.String("namespace", namespace),
		zap.Int("page", pagination.Page),
		zap.Int("limit", pagination.Limit),
		zap.Int("filter_count", len(filters)))

	// Process the search
	response, err := h.service.ProcessSearch(ctx, query, filters, pagination, namespace)
	if err != nil {
		logger.Error(ctx, "search processing failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode search response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "GET search completed successfully",
		zap.Int("result_count", len(response.Data)))
}

// SearchConversations handles search requests specifically for conversations namespace
func (h *SearchServiceHandler) SearchConversations(w http.ResponseWriter, r *http.Request) {
	h.handleNamespaceSearch(w, r, "conversations")
}

// SearchOrganizations handles search requests specifically for organizations namespace
func (h *SearchServiceHandler) SearchOrganizations(w http.ResponseWriter, r *http.Request) {
	h.handleNamespaceSearch(w, r, "organizations")
}

// SearchAll handles search requests across all namespaces
func (h *SearchServiceHandler) SearchAll(w http.ResponseWriter, r *http.Request) {
	h.handleNamespaceSearch(w, r, "")
}

// handleNamespaceSearch is a helper for namespace-specific search endpoints
func (h *SearchServiceHandler) handleNamespaceSearch(w http.ResponseWriter, r *http.Request, namespace string) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, fmt.Sprintf("search-api:search-%s", strings.ToLower(namespace)))
	defer span.End()

	var query string
	var filters map[string]string
	var pagination PaginationParams

	if r.Method == http.MethodPost {
		// Handle POST request
		var searchReq SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
			logger.Error(ctx, "failed to decode search request", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		query = searchReq.Query
		filters = searchReq.Filters
		pagination = PaginationParams{
			Page:  searchReq.Page,
			Limit: searchReq.PerPage,
		}
		if pagination.Limit <= 0 {
			pagination.Limit = 20
		}
	} else {
		// Handle GET request
		query = r.URL.Query().Get("q")
		pagination = h.extractPagination(r)
		filters = h.extractFilters(r)
	}

	if query == "" {
		logger.Error(ctx, "empty search query provided")
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "processing namespace search request",
		zap.String("query", query),
		zap.String("namespace", namespace),
		zap.Int("page", pagination.Page),
		zap.Int("limit", pagination.Limit))

	// Process the search with namespace
	response, err := h.service.ProcessSearch(ctx, query, filters, pagination, namespace)
	if err != nil {
		logger.Error(ctx, "namespace search processing failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode search response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "namespace search completed successfully",
		zap.String("namespace", namespace),
		zap.Int("result_count", len(response.Data)))
}

// GetAvailableFilters returns legacy filter format for backwards compatibility
func (h *SearchServiceHandler) GetAvailableFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:get-available-filters")
	defer span.End()

	logger.Info(ctx, "get available filters request received")

	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := client.ListFilters(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get filters from fugu", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "available filters served successfully")
}

// extractPagination extracts pagination parameters from query string
func (h *SearchServiceHandler) extractPagination(r *http.Request) PaginationParams {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	perPageStr := r.URL.Query().Get("per_page")

	page := 0
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			page = p
		}
	}

	limit := 20 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	} else if perPageStr != "" {
		if l, err := strconv.Atoi(perPageStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// extractFilters extracts filter parameters from query string
func (h *SearchServiceHandler) extractFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)

	// Extract filter parameters
	for key, values := range r.URL.Query() {
		if len(values) > 0 && values[0] != "" {
			// Skip pagination and query parameters
			if key != "page" && key != "limit" && key != "per_page" && key != "q" && key != "namespace" {
				filters[key] = values[0]
			}
		}
	}

	return filters
}

// HealthCheck performs a comprehensive health check for the search service
func (h *SearchServiceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:health")
	defer span.End()

	logger.Info(ctx, "search health request received")

	// Create fugu client and test connection
	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client for health check", zap.Error(err))
		h.respondHealthError(w, "Failed to create fugu client", err.Error())
		return
	}

	// Test fugu backend health
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Health(healthCtx); err != nil {
		logger.Error(ctx, "fugu server health check failed", zap.Error(err))
		h.respondHealthError(w, "Fugu backend unavailable", err.Error())
		return
	}

	// All checks passed
	h.respondHealthSuccess(w)
	logger.Info(ctx, "search health check completed successfully")
}

// respondHealthError responds with an unhealthy status
func (h *SearchServiceHandler) respondHealthError(w http.ResponseWriter, errorMsg, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "unhealthy",
		"error":     errorMsg,
		"details":   details,
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "search",
	})
}

// respondHealthSuccess responds with a healthy status and capabilities
func (h *SearchServiceHandler) respondHealthSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "search",
		"version":   "1.1.0",
		"backend": map[string]interface{}{
			"name":   "fugu",
			"url":    h.service.fuguServerURL,
			"status": "healthy",
		},
		"capabilities": map[string]bool{
			"search":          true,
			"filters":         true,
			"pagination":      true,
			"namespaces":      true,
			"facets":          true,
			"health_check":    true,
			"facet_filtering": true,
		},
		"endpoints": []string{
			"GET/POST /search",
			"GET/POST /search/conversations",
			"GET/POST /search/organizations",
			"GET/POST /search/all",
			"GET /search/info",
			"GET /search/health",
			"GET /search/filters",
			"GET /search/filters/all",
			"GET /search/filters/namespace/{namespace}",
			"GET /search/filters/path/{filterPath}",
			"POST /search/filters/invalidate",
		},
		"supported_filters": []string{
			"metadata field filters (e.g., docket_gov_id=EPA-123)",
			"namespace filters (conversations, organizations)",
			"boolean filters (is_person=true)",
			"numeric filters (total_documents=5)",
		},
	}

	json.NewEncoder(w).Encode(response)
}
