package search

import (
	"context"
	"encoding/json"
	"fmt"
	filterconfig "kessler/internal/filters/config"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("search-endpoint")

// Frontend request types
type SearchRequest struct {
	Query string `json:"query"`
}

// Frontend response types
type SearchResponse struct {
	Data []SearchResultItem `json:"data"`
}

// SearchResultItem represents a single search result for the frontend
type SearchResultItem struct {
	ID          string                 `json:"id"`
	Score       float32                `json:"score"`
	Metadata    map[string]interface{} `json:"metadata"`
	Facet       []string               `json:"facet,omitempty"`
	CaseNumber  string                 `json:"case_number,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	Description string                 `json:"description,omitempty"`
	FileName    string                 `json:"file_name,omitempty"`
	FiledDate   string                 `json:"filed_date,omitempty"`
	FilingType  string                 `json:"filing_type,omitempty"`
	PartyName   string                 `json:"party_name,omitempty"`
	// Text        string                 `json:"text,omitempty"`
}

// Fugu search request types
type FuguPagination struct {
	Page    *int `json:"page,omitempty"`
	PerPage *int `json:"per_page,omitempty"`
}

type FuguSearchQuery struct {
	Query   string          `json:"query"`
	Filters []string        `json:"filters,omitempty"`
	Page    *FuguPagination `json:"page,omitempty"`
}

type FuguSearchResult struct {
	Facet    []string                 `json:"facet"`
	ID       []string                 `json:"id"`
	Metadata []map[string]interface{} `json:"metadata"`
	Score    float32                  `json:"score"` // Changed to float64
	// Text     []string                 `json:"text"`
}

type FuguSearchResponse struct {
	Result []FuguSearchResult `json:"result"`
}

// Pagination extracted from query parameters
type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// SearchInfo represents search service information and capabilities
type SearchInfo struct {
	Status       string             `json:"status"`
	Version      string             `json:"version"`
	Capabilities SearchCapabilities `json:"capabilities"`
	Statistics   SearchStatistics   `json:"statistics"`
	LastUpdated  string             `json:"last_updated"`
}

// SearchCapabilities represents what the search service can do
type SearchCapabilities struct {
	FilterSupport     bool     `json:"filter_support"`
	PaginationSupport bool     `json:"pagination_support"`
	SortingSupport    bool     `json:"sorting_support"`
	HighlightSupport  bool     `json:"highlight_support"`
	FacetSupport      bool     `json:"facet_support"`
	SupportedQueries  []string `json:"supported_queries"`
	MaxQueryLength    int      `json:"max_query_length"`
	MaxResultsPerPage int      `json:"max_results_per_page"`
}

// SearchStatistics represents search service statistics
type SearchStatistics struct {
	TotalDocuments   int64    `json:"total_documents"`
	IndexedFields    []string `json:"indexed_fields"`
	AvailableFilters []string `json:"available_filters"`
	BackendStatus    string   `json:"backend_status"`
	LastIndexUpdate  string   `json:"last_index_update,omitempty"`
}

// SearchService handles the business logic for search operations
type SearchService struct {
	fuguServerURL   string
	filterConfigSvc *filterconfig.FilterConfigService
}

// NewSearchService creates a new search service
func NewSearchService(fuguServerURL string, filterConfigSvc *filterconfig.FilterConfigService) *SearchService {
	return &SearchService{
		fuguServerURL:   fuguServerURL,
		filterConfigSvc: filterConfigSvc,
	}
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

// convertFiltersToBackend converts frontend filters to backend format
func (s *SearchService) convertFiltersToBackend(ctx context.Context, filters map[string]string) ([]string, error) {
	if len(filters) == 0 {
		return []string{}, nil
	}

	logger.Info(ctx, "converting filters to backend format",
		zap.Int("filter_count", len(filters)),
		zap.Any("filters", filters))

	// If no filter config service, just pass through as simple key:value pairs
	if s.filterConfigSvc == nil {
		logger.Warn(ctx, "no filter config service available, using simple conversion")
		var filterStrings []string
		for key, value := range filters {
			if value != "" {
				filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", key, value))
			}
		}
		return filterStrings, nil
	}

	// Try to use filter configuration service
	backendFilterMap, err := s.filterConfigSvc.ConvertFiltersToBackend(ctx, filters)
	if err != nil {
		logger.Warn(ctx, "filter config service conversion failed, falling back to simple conversion", zap.Error(err))
		// Fallback to simple conversion
		var filterStrings []string
		for key, value := range filters {
			if value != "" {
				filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", key, value))
			}
		}
		return filterStrings, nil
	}

	// Convert map to string array format expected by fugu
	var filterStrings []string
	for key, value := range backendFilterMap {
		if valueStr, ok := value.(string); ok && valueStr != "" {
			// Format: "field:value"
			filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", key, valueStr))
		}
	}

	logger.Info(ctx, "filters converted successfully",
		zap.Int("backend_filter_count", len(filterStrings)),
		zap.Strings("backend_filters", filterStrings))

	return filterStrings, nil
}

// Update the createFuguSearchQuery method to use SDK types
func (s *SearchService) createFuguSearchQuery(query string, filters []string, pagination PaginationParams) fugusdk.FuguSearchQuery {
	// Convert our internal pagination to SDK pagination
	var fuguPagination *fugusdk.Pagination
	if pagination.Page > 0 || pagination.Limit > 0 {
		fuguPagination = &fugusdk.Pagination{
			Page:    &pagination.Page,
			PerPage: &pagination.Limit,
		}
	}

	// Convert filters slice to pointer to slice
	var filtersPtr *[]string
	if len(filters) > 0 {
		filtersPtr = &filters
	}

	return fugusdk.FuguSearchQuery{
		Query:   query,
		Filters: filtersPtr,
		Page:    fuguPagination,
	}
}

// In your search handler file, make sure these types match:

// FuguSearchResult for the internal conversion (the old one with arrays)

// Update the executeSearch function with proper type conversion:
func (s *SearchService) executeSearch(ctx context.Context, client *fugusdk.Client, query fugusdk.FuguSearchQuery) (*FuguSearchResponse, error) {
	ctx, span := tracer.Start(ctx, "search-service:execute-search")
	defer span.End()

	logger.Info(ctx, "executing search on fugu backend")

	// Health check
	healthCtx, healthCancel := context.WithTimeout(ctx, 5*time.Second)
	defer healthCancel()

	if err := client.Health(healthCtx); err != nil {
		logger.Error(ctx, "fugu health check failed", zap.Error(err))
		return nil, fmt.Errorf("fugu backend unhealthy: %w", err)
	}

	logger.Info(ctx, "fugu health check passed")

	logger.Info(ctx, "sending search query to fugu",
		zap.String("query", query.Query),
		zap.Any("filters", query.Filters),
		zap.Any("page", query.Page))

	// Make the search request
	response, err := client.Search(ctx, query)
	if err != nil {
		logger.Error(ctx, "fugu search failed", zap.Error(err))
		return nil, fmt.Errorf("fugu search failed: %w", err)
	}

	logger.Info(ctx, "received response from fugu",
		zap.Int("result_count", len(response.Results)),
		zap.Int("total", response.Total),
		zap.String("message", response.Message))

	// Convert Fugu's response format to our internal format
	var fuguResponse FuguSearchResponse

	if len(response.Results) > 0 {
		logger.Info(ctx, "converting fugu results to internal format",
			zap.Int("result_count", len(response.Results)))

		fuguResponse = FuguSearchResponse{
			Result: make([]FuguSearchResult, len(response.Results)),
		}

		for i, result := range response.Results {
			fuguResponse.Result[i] = FuguSearchResult{
				ID: []string{result.ID},
				// Text:     []string{result.Text},
				Score:    result.Score, // Now both are float64
				Metadata: []map[string]interface{}{},
				Facet:    result.Facets,
			}

			// Handle metadata
			if result.Metadata != nil {
				fuguResponse.Result[i].Metadata = []map[string]interface{}{result.Metadata}
			}
		}

		logger.Info(ctx, "successfully converted fugu response",
			zap.Int("converted_result_count", len(fuguResponse.Result)))
	} else {
		logger.Info(ctx, "no results in fugu response, returning empty results")
		fuguResponse = FuguSearchResponse{Result: []FuguSearchResult{}}
	}

	return &fuguResponse, nil
}

// Update the ProcessSearch method to use SDK types
func (s *SearchService) ProcessSearch(ctx context.Context, query string, filters map[string]string, pagination PaginationParams) (*SearchResponse, error) {
	ctx, span := tracer.Start(ctx, "search-service:process-search")
	defer span.End()

	logger.Info(ctx, "starting search processing",
		zap.String("query", query),
		zap.String("fugu_url", s.fuguServerURL))

	// Create fugu client with timeout
	clientCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := fugusdk.NewClient(clientCtx, s.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	logger.Info(ctx, "fugu client created successfully")

	// Convert filters to backend format
	backendFilters, err := s.convertFiltersToBackend(ctx, filters)
	if err != nil {
		logger.Warn(ctx, "failed to convert filters, proceeding without filters", zap.Error(err))
		backendFilters = []string{}
	}

	logger.Info(ctx, "filters converted",
		zap.Int("original_filter_count", len(filters)),
		zap.Int("backend_filter_count", len(backendFilters)))

	// Create fugu search query using SDK types
	fuguQuery := s.createFuguSearchQuery(query, backendFilters, pagination)

	logger.Info(ctx, "created fugu query",
		zap.String("query", fuguQuery.Query),
		zap.Any("filters", fuguQuery.Filters),
		zap.Any("page", fuguQuery.Page))

	// Execute search on fugu with timeout
	searchCtx, searchCancel := context.WithTimeout(ctx, 15*time.Second)
	defer searchCancel()

	fuguResponse, err := s.executeSearch(searchCtx, client, fuguQuery)
	if err != nil {
		logger.Error(ctx, "fugu search execution failed", zap.Error(err))
		return nil, fmt.Errorf("fugu search failed: %w", err)
	}

	logger.Info(ctx, "fugu search completed",
		zap.Int("result_count", len(fuguResponse.Result)))

	// Transform fugu response to frontend format
	frontendResponse := s.transformSearchResponse(fuguResponse)

	logger.Info(ctx, "search processing completed successfully",
		zap.Int("final_result_count", len(frontendResponse.Data)))

	return frontendResponse, nil
}

// Add a fallback search handler for testing
func (h *SearchServiceHandler) SearchWithFallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:search-with-fallback")
	defer span.End()

	logger.Info(ctx, "search request received (with fallback)")

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

	// Extract pagination and filters
	pagination := h.extractPagination(r)
	filters := h.extractFilters(r)

	logger.Info(ctx, "processing search request",
		zap.String("query", searchReq.Query),
		zap.Int("page", pagination.Page),
		zap.Int("limit", pagination.Limit),
		zap.Int("filter_count", len(filters)))

	// Try the real search first
	response, err := h.service.ProcessSearch(ctx, searchReq.Query, filters, pagination)
	if err != nil {
		logger.Warn(ctx, "real search failed, using fallback", zap.Error(err))

		// Fallback to mock response
		response = &SearchResponse{
			Data: []SearchResultItem{
				{
					ID:    "mock-1",
					Score: 0.95,
					// Text:        "This is a mock search result for query: " + searchReq.Query,
					Metadata:    map[string]interface{}{"source": "mock", "type": "fallback"},
					CaseNumber:  "MOCK-001",
					FileName:    "mock_document.pdf",
					FilingType:  "Mock Filing",
					Description: "Mock search result returned because backend is unavailable",
				},
				{
					ID:    "mock-2",
					Score: 0.87,
					// Text:        "Another mock result to test search functionality",
					Metadata:    map[string]interface{}{"source": "mock", "type": "fallback"},
					CaseNumber:  "MOCK-002",
					FileName:    "another_mock.pdf",
					FilingType:  "Test Filing",
					Description: "Second mock result for testing",
				},
			},
		}

		logger.Info(ctx, "using fallback response",
			zap.Int("mock_result_count", len(response.Data)))
	}

	logger.Info(ctx, "search processing completed successfully",
		zap.Int("result_count", len(response.Data)))

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	logger.Info(ctx, "about to encode response")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode search response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "response encoded and sent successfully")
}

// Update the route registration to use the fallback handler temporarily
func RegisterSearchRoutes(router *mux.Router) error {
	fuguServerURL := "http://fugudb:3301" // Make sure this matches your docker-compose
	filterConfigSvc := filterconfig.NewFilterConfigService(fuguServerURL)
	service := NewSearchService(fuguServerURL, filterConfigSvc)
	handler := NewSearchHandler(service)

	fmt.Println("🔧 Registering search routes with fallback...")

	// GET endpoint for search info/status
	router.HandleFunc("/", handler.GetSearchInfo).Methods(http.MethodGet)
	router.HandleFunc("", handler.GetSearchInfo).Methods(http.MethodGet)

	// Main search endpoint with fallback
	router.HandleFunc("/", handler.SearchWithFallback).Methods(http.MethodPost)
	router.HandleFunc("", handler.SearchWithFallback).Methods(http.MethodPost)

	// File search endpoint (for compatibility)
	router.HandleFunc("/file", handler.SearchFiles).Methods(http.MethodPost)

	// Search health check
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	fmt.Println("✅ Search routes registered with fallback")
	return nil
}

// GetSearchInfo handles GET requests to /search endpoint
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

func (h *SearchServiceHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:search")
	defer span.End()

	logger.Info(ctx, "search request received")

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

	// Extract pagination from query parameters
	pagination := h.extractPagination(r)

	// Extract filters from query parameters if any
	filters := h.extractFilters(r)

	logger.Info(ctx, "processing search request",
		zap.String("query", searchReq.Query),
		zap.Int("page", pagination.Page),
		zap.Int("limit", pagination.Limit),
		zap.Int("filter_count", len(filters)))

	// Process the search
	response, err := h.service.ProcessSearch(ctx, searchReq.Query, filters, pagination)
	if err != nil {
		logger.Error(ctx, "search processing failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ADD THIS DEBUG LOGGING
	logger.Info(ctx, "search processing completed successfully",
		zap.Int("result_count", len(response.Data)))

	// Return response
	w.Header().Set("Content-Type", "application/json")

	// ADD THIS DEBUG LOGGING
	logger.Info(ctx, "about to encode response")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode search response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ADD THIS DEBUG LOGGING
	logger.Info(ctx, "response encoded and sent successfully")

	logger.Info(ctx, "search completed successfully",
		zap.Int("result_count", len(response.Data)))
}

// SearchFiles handles file-specific search requests (for frontend compatibility)
func (h *SearchServiceHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
	// For now, this just delegates to the main search handler
	// Could be specialized later for file-specific logic
	h.Search(w, r)
}

// extractPagination extracts pagination parameters from query string
func (h *SearchServiceHandler) extractPagination(r *http.Request) PaginationParams {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

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
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// extractFilters extracts filter parameters from query string
func (h *SearchServiceHandler) extractFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)

	// Extract filter parameters that start with "filter_"
	for key, values := range r.URL.Query() {
		if len(values) > 0 && values[0] != "" {
			// All query parameters are potential filters
			filters[key] = values[0]
		}
	}

	// Remove pagination parameters
	delete(filters, "page")
	delete(filters, "limit")

	return filters
}

// GetSearchInfo provides information about search capabilities and status
func (s *SearchService) GetSearchInfo(ctx context.Context) (*SearchInfo, error) {
	ctx, span := tracer.Start(ctx, "search-service:get-search-info")
	defer span.End()

	// Create fugu client to get backend status
	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Check backend health
	backendStatus := "healthy"
	if err := client.Health(ctx); err != nil {
		backendStatus = "unhealthy"
		logger.Warn(ctx, "fugu backend health check failed", zap.Error(err))
	}

	// Get available filters from filter config service
	config, err := s.filterConfigSvc.BuildFilterConfiguration(ctx)
	var availableFilters []string
	var indexedFields []string

	if err != nil {
		logger.Warn(ctx, "failed to get filter configuration", zap.Error(err))
		availableFilters = []string{}
		indexedFields = []string{}
	} else {
		// Extract filter field names
		for _, field := range config.Fields {
			if field.Enabled {
				availableFilters = append(availableFilters, field.ID)
				indexedFields = append(indexedFields, field.BackendKey)
			}
		}
	}

	// Build search info response
	info := &SearchInfo{
		Status:      "operational",
		Version:     "1.0.0",
		LastUpdated: time.Now().Format(time.RFC3339),
		Capabilities: SearchCapabilities{
			FilterSupport:     true,
			PaginationSupport: true,
			SortingSupport:    false, // Not implemented yet
			HighlightSupport:  false, // Not implemented yet
			FacetSupport:      true,
			SupportedQueries: []string{
				"simple_text",
				"boolean_operators",
				"phrase_search",
				"field_targeting",
				"range_queries",
				"wildcard_search",
			},
			MaxQueryLength:    10000,
			MaxResultsPerPage: 100,
		},
		Statistics: SearchStatistics{
			TotalDocuments:   0, // Would need to query fugu for this
			IndexedFields:    indexedFields,
			AvailableFilters: availableFilters,
			BackendStatus:    backendStatus,
		},
	}

	// Try to get document count from fugu (optional)
	if backendStatus == "healthy" {
		// This would require a specific endpoint on fugu to get document count
		// For now, we'll leave it as 0 or could make a simple search to estimate
	}

	return info, nil
}

// transformSearchResponse transforms fugu response to frontend format
func (s *SearchService) transformSearchResponse(fuguResponse *FuguSearchResponse) *SearchResponse {
	if fuguResponse == nil || len(fuguResponse.Result) == 0 {
		return &SearchResponse{Data: []SearchResultItem{}}
	}

	var frontendResults []SearchResultItem

	for _, result := range fuguResponse.Result {
		// Extract the primary ID (first ID if multiple)
		var primaryID string
		if len(result.ID) > 0 {
			primaryID = result.ID[0]
		}

		// Extract text (join if multiple text segments)
		// var text string
		// if len(result.Text) > 0 {
		// 	// For now, take the first text segment
		// 	// Could be modified to join all segments if needed
		// 	text = result.Text[0]
		// }

		// Extract metadata (use first metadata object if multiple)
		var metadata map[string]interface{}
		if len(result.Metadata) > 0 {
			metadata = result.Metadata[0]
		} else {
			metadata = make(map[string]interface{})
		}

		// Create frontend result item
		item := SearchResultItem{
			ID:    primaryID,
			Score: result.Score,
			// Text:     text,
			Metadata: metadata,
			Facet:    result.Facet,
		}

		// Extract commonly used metadata fields for easier frontend access
		if metadata != nil {
			if caseNumber, ok := metadata["case_number"].(string); ok {
				item.CaseNumber = caseNumber
			}
			if createdAt, ok := metadata["created_at"].(string); ok {
				item.CreatedAt = createdAt
			}
			if description, ok := metadata["description"].(string); ok {
				item.Description = description
			}
			if fileName, ok := metadata["file_name"].(string); ok {
				item.FileName = fileName
			}
			if filedDate, ok := metadata["filed_date"].(string); ok {
				item.FiledDate = filedDate
			}
			if filingType, ok := metadata["filing_type"].(string); ok {
				item.FilingType = filingType
			}
			if partyName, ok := metadata["party_name"].(string); ok {
				item.PartyName = partyName
			}
		}

		frontendResults = append(frontendResults, item)
	}

	return &SearchResponse{
		Data: frontendResults,
	}
}

// Health check specifically for search functionality
func (h *SearchServiceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:health")
	defer span.End()
	logger.Info(ctx, "search health request received")

	// Create fugu client and test connection
	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client for health check", zap.Error(err))
		http.Error(w, "Search service unavailable", http.StatusServiceUnavailable)
		return
	}

	if err := client.Health(ctx); err != nil {
		logger.Error(ctx, "fugu server health check failed", zap.Error(err))
		http.Error(w, "Search backend unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "search",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode health response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
