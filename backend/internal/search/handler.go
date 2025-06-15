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
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("search-endpoint")

func RegisterSearchRoutes(router *mux.Router) error {
	fuguServerURL := "http://fugudb:3301" // Make sure this matches your docker-compose
	filterConfigSvc := filterconfig.NewFilterConfigService(fuguServerURL)
	filterConfigHandler := filterconfig.NewFilterConfigHandler(filterConfigSvc)
	service := NewSearchService(fuguServerURL, filterConfigSvc)
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

	// Filter configuration endpoints - using the actual filterconfig handlers
	router.HandleFunc("/filters", handler.GetAvailableFilters).Methods(http.MethodGet)
	router.HandleFunc("/filters/{namespace}", filterConfigHandler.GetNamespaceFilters).Methods(http.MethodGet)
	router.HandleFunc("/filters/configuration", filterConfigHandler.GetConfiguration).Methods(http.MethodGet)
	router.HandleFunc("/filters/convert", filterConfigHandler.ConvertFilters).Methods(http.MethodPost)
	router.HandleFunc("/filters/validate", filterConfigHandler.ValidateFilters).Methods(http.MethodPost)
	router.HandleFunc("/filters/options", filterConfigHandler.GetOptions).Methods(http.MethodPost)

	fmt.Println("✅ Search routes registered successfully")
	return nil
}

// Frontend request types
type SearchRequest struct {
	Query     string            `json:"query"`
	Filters   map[string]string `json:"filters,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Page      int               `json:"page,omitempty"`
	PerPage   int               `json:"per_page,omitempty"`
}

// Frontend response types
type SearchResponse struct {
	Data        []SearchResultItem `json:"data"`
	Total       int                `json:"total,omitempty"`
	Page        int                `json:"page,omitempty"`
	PerPage     int                `json:"per_page,omitempty"`
	Query       string             `json:"query,omitempty"`
	Namespace   string             `json:"namespace,omitempty"`
	ProcessTime string             `json:"process_time,omitempty"`
}

// SearchResultItem represents a single search result for the frontend
type SearchResultItem struct {
	ID          string                 `json:"id"`
	Score       float32                `json:"score"`
	Text        string                 `json:"text,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	Facet       []string               `json:"facet,omitempty"`
	Namespace   string                 `json:"namespace,omitempty"`
	CaseNumber  string                 `json:"case_number,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	Description string                 `json:"description,omitempty"`
	FileName    string                 `json:"file_name,omitempty"`
	FiledDate   string                 `json:"filed_date,omitempty"`
	FilingType  string                 `json:"filing_type,omitempty"`
	PartyName   string                 `json:"party_name,omitempty"`
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
	Score    float32                  `json:"score"`
	Text     []string                 `json:"text"`
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
	FilterSupport       bool     `json:"filter_support"`
	PaginationSupport   bool     `json:"pagination_support"`
	SortingSupport      bool     `json:"sorting_support"`
	HighlightSupport    bool     `json:"highlight_support"`
	FacetSupport        bool     `json:"facet_support"`
	NamespaceSupport    bool     `json:"namespace_support"`
	SupportedQueries    []string `json:"supported_queries"`
	MaxQueryLength      int      `json:"max_query_length"`
	MaxResultsPerPage   int      `json:"max_results_per_page"`
	SupportedNamespaces []string `json:"supported_namespaces"`
}

// SearchStatistics represents search service statistics
type SearchStatistics struct {
	TotalDocuments   int64            `json:"total_documents"`
	IndexedFields    []string         `json:"indexed_fields"`
	AvailableFilters []string         `json:"available_filters"`
	BackendStatus    string           `json:"backend_status"`
	LastIndexUpdate  string           `json:"last_index_update,omitempty"`
	NamespaceStats   map[string]int64 `json:"namespace_stats,omitempty"`
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

// convertFiltersToBackend converts frontend filters to backend format
func (s *SearchService) convertFiltersToBackend(ctx context.Context, filters map[string]string, namespace string) ([]string, error) {
	var filterStrings []string

	// Add namespace filter if specified
	if namespace != "" {
		filterStrings = append(filterStrings, namespace)
	}

	if len(filters) == 0 {
		return filterStrings, nil
	}

	logger.Info(ctx, "converting filters to backend format",
		zap.Int("filter_count", len(filters)),
		zap.Any("filters", filters),
		zap.String("namespace", namespace))

	// If no filter config service, just pass through as simple key:value pairs
	if s.filterConfigSvc == nil {
		logger.Warn(ctx, "no filter config service available, using simple conversion")
		for key, value := range filters {
			if value != "" && key != "q" && key != "page" && key != "per_page" && key != "limit" {
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
		for key, value := range filters {
			if value != "" && key != "q" && key != "page" && key != "per_page" && key != "limit" {
				filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", key, value))
			}
		}
		return filterStrings, nil
	}

	// Convert map to string array format expected by fugu
	for key, value := range backendFilterMap {
		if valueStr, ok := value.(string); ok && valueStr != "" {
			filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", key, valueStr))
		}
	}

	logger.Info(ctx, "filters converted successfully",
		zap.Int("backend_filter_count", len(filterStrings)),
		zap.Strings("backend_filters", filterStrings))

	return filterStrings, nil
}

// createFuguSearchQuery creates a Fugu search query from the request parameters
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

// executeSearch executes the search against Fugu backend
func (s *SearchService) executeSearch(ctx context.Context, client *fugusdk.Client, query fugusdk.FuguSearchQuery) (*fugusdk.SanitizedResponse, error) {
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

	// Make the search request using the SDK
	response, err := client.Search(ctx, query)
	if err != nil {
		logger.Error(ctx, "fugu search failed", zap.Error(err))
		return nil, fmt.Errorf("fugu search failed: %w", err)
	}

	logger.Info(ctx, "received response from fugu",
		zap.Int("result_count", len(response.Results)),
		zap.Int("total", response.Total),
		zap.String("message", response.Message))

	return response, nil
}

// ProcessSearch processes a search request with namespace support
func (s *SearchService) ProcessSearch(ctx context.Context, query string, filters map[string]string, pagination PaginationParams, namespace string) (*SearchResponse, error) {
	ctx, span := tracer.Start(ctx, "search-service:process-search")
	defer span.End()

	startTime := time.Now()

	logger.Info(ctx, "starting search processing",
		zap.String("query", query),
		zap.String("namespace", namespace),
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

	// Convert filters to backend format with namespace
	backendFilters, err := s.convertFiltersToBackend(ctx, filters, namespace)
	if err != nil {
		logger.Warn(ctx, "failed to convert filters, proceeding without filters", zap.Error(err))
		backendFilters = []string{}
		if namespace != "" {
			backendFilters = []string{namespace}
		}
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
		zap.Int("result_count", len(fuguResponse.Results)))

	// Transform fugu response to frontend format
	frontendResponse := s.transformSearchResponse(fuguResponse, query, namespace, pagination, time.Since(startTime))

	logger.Info(ctx, "search processing completed successfully",
		zap.Int("final_result_count", len(frontendResponse.Data)))

	return frontendResponse, nil
}
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
		Version:     "1.1.0",
		LastUpdated: time.Now().Format(time.RFC3339),
		Capabilities: SearchCapabilities{
			FilterSupport:     true,
			PaginationSupport: true,
			SortingSupport:    false,
			HighlightSupport:  false,
			FacetSupport:      true,
			NamespaceSupport:  true,
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
			SupportedNamespaces: []string{
				"conversations",
				"organizations",
			},
		},
		Statistics: SearchStatistics{
			TotalDocuments:   0, // Would need to query fugu for this
			IndexedFields:    indexedFields,
			AvailableFilters: availableFilters,
			BackendStatus:    backendStatus,
		},
	}

	return info, nil
}

// transformSearchResponse transforms fugu response to frontend format
func (s *SearchService) transformSearchResponse(fuguResponse *fugusdk.SanitizedResponse, query, namespace string, pagination PaginationParams, processTime time.Duration) *SearchResponse {
	if fuguResponse == nil || len(fuguResponse.Results) == 0 {
		return &SearchResponse{
			Data:        []SearchResultItem{},
			Total:       0,
			Page:        pagination.Page,
			PerPage:     pagination.Limit,
			Query:       query,
			Namespace:   namespace,
			ProcessTime: processTime.String(),
		}
	}

	var frontendResults []SearchResultItem

	for _, result := range fuguResponse.Results {
		// Create frontend result item
		item := SearchResultItem{
			ID:       result.ID,
			Score:    result.Score,
			Text:     result.Text,
			Metadata: result.Metadata,
			Facet:    result.Facets,
		}

		// Extract namespace from facets if available
		if len(result.Facets) > 0 {
			item.Namespace = result.Facets[0]
		}

		// Extract commonly used metadata fields for easier frontend access
		if result.Metadata != nil {
			if caseNumber, ok := result.Metadata["case_number"].(string); ok {
				item.CaseNumber = caseNumber
			}
			if createdAt, ok := result.Metadata["created_at"].(string); ok {
				item.CreatedAt = createdAt
			}
			if description, ok := result.Metadata["description"].(string); ok {
				item.Description = description
			}
			if fileName, ok := result.Metadata["file_name"].(string); ok {
				item.FileName = fileName
			}
			if filedDate, ok := result.Metadata["filed_date"].(string); ok {
				item.FiledDate = filedDate
			}
			if filingType, ok := result.Metadata["filing_type"].(string); ok {
				item.FilingType = filingType
			}
			if partyName, ok := result.Metadata["party_name"].(string); ok {
				item.PartyName = partyName
			}
			if docketGovID, ok := result.Metadata["docket_gov_id"].(string); ok {
				item.Metadata["docket_gov_id"] = docketGovID
			}
			if totalDocuments, ok := result.Metadata["total_documents"]; ok {
				item.Metadata["total_documents"] = totalDocuments
			}
			if totalDocumentsAuthored, ok := result.Metadata["total_documents_authored"]; ok {
				item.Metadata["total_documents_authored"] = totalDocumentsAuthored
			}
			if isPerson, ok := result.Metadata["is_person"]; ok {
				item.Metadata["is_person"] = isPerson
			}
		}

		frontendResults = append(frontendResults, item)
	}

	return &SearchResponse{
		Data:        frontendResults,
		Total:       fuguResponse.Total,
		Page:        pagination.Page,
		PerPage:     pagination.Limit,
		Query:       query,
		Namespace:   namespace,
		ProcessTime: processTime.String(),
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

// GetAvailableFilters returns all available filters
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

// GetNamespaceFilters returns filters for a specific namespace
func (h *SearchServiceHandler) GetNamespaceFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "search-api:get-namespace-filters")
	defer span.End()

	namespace := mux.Vars(r)["namespace"]
	logger.Info(ctx, "get namespace filters request received", zap.String("namespace", namespace))

	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := client.GetNamespaceFilters(ctx, namespace)
	if err != nil {
		logger.Error(ctx, "failed to get namespace filters from fugu", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode namespace filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "namespace filters served successfully", zap.String("namespace", namespace))
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

// GetSearchInfo provides information about search capabilities and status

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "unhealthy",
			"error":     "Failed to create fugu client",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "search",
		})
		return
	}

	// Test fugu backend health
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Health(healthCtx); err != nil {
		logger.Error(ctx, "fugu server health check failed", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "unhealthy",
			"error":     "Fugu backend unavailable",
			"details":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "search",
		})
		return
	}

	// All checks passed
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
			"search":       true,
			"filters":      true,
			"pagination":   true,
			"namespaces":   true,
			"facets":       true,
			"health_check": true,
		},
		"endpoints": []string{
			"GET/POST /search",
			"GET/POST /search/conversations",
			"GET/POST /search/organizations",
			"GET/POST /search/all",
			"GET /search/info",
			"GET /search/health",
			"GET /search/filters",
			"GET /search/filters/{namespace}",
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode health response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "search health check completed successfully")
}
