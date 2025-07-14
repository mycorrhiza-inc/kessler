// service.go
package search

import (
	"context"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/internal/search/filter"
	"kessler/pkg/logger"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var serviceTracer = otel.Tracer("search-service")

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
	Data        []CardData `json:"data"`
	Total       int        `json:"total,omitempty"`
	Page        int        `json:"page,omitempty"`
	PerPage     int        `json:"per_page,omitempty"`
	Query       string     `json:"query,omitempty"`
	Namespace   string     `json:"namespace,omitempty"`
	ProcessTime string     `json:"process_time,omitempty"`
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

// CardData represents the frontend card format

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
	fuguServerURL string
	filterService *filter.Service
	db            dbstore.DBTX
	cacheCtrl     cache.CacheController
	cacheEnabled  bool
}

// NewSearchService creates a new search service
func NewSearchService(fuguServerURL string, filterService *filter.Service, db dbstore.DBTX) (*SearchService, error) {
	cacheCtrl, err := cache.NewCacheController()
	cacheEnabled := err == nil

	if !cacheEnabled {
		logger.Warn(context.Background(), "failed to initialize search cache controller", zap.Error(err))
	}

	return &SearchService{
		fuguServerURL: fuguServerURL,
		filterService: filterService,
		db:            db,
		cacheCtrl:     cacheCtrl,
		cacheEnabled:  cacheEnabled,
	}, nil
}

// ProcessSearch processes a search request with namespace support
func (s *SearchService) ProcessSearch(ctx context.Context, query string, metadataFilters map[string]string, pagination PaginationParams, namespace string) (*SearchResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "search-service:process-search")
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

	rawFilters := convertMetadataFiltersToRaw(metadataFilters)

	// Convert filters to backend format with namespace
	backendFilters, err := convertFiltersToBackend(ctx, rawFilters, namespace)
	if err != nil {
		logger.Warn(ctx, "failed to convert filters, proceeding with fallback", zap.Error(err))
		backendFilters = s.fallbackFilterConversion(rawFilters, namespace)
	}

	// Create fugu search query using SDK types
	fuguQuery := s.createFuguSearchQuery(query, backendFilters, pagination)

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
	frontendResponse, err := s.transformSearchResponse(ctx, fuguResponse, query, namespace, pagination, time.Since(startTime))
	if err != nil {
		logger.Error(ctx, "failed to transform search response", zap.Error(err))
		return nil, fmt.Errorf("failed to transform response: %w", err)
	}

	logger.Info(ctx, "search processing completed successfully",
		zap.Int("final_result_count", len(frontendResponse.Data)))

	return frontendResponse, nil
}

var metadataFilterRenameDict map[string]string = map[string]string{
	"convo_id": "conversation_id",
}

// convertMetadataFiltersToRaw takes a map of metadata filters and returns a new map
// with the same values but with "metadata/" prepended to each key.
func convertMetadataFiltersToRaw(metadataFilters map[string]string) map[string]string {
	rawMetadataFilters := make(map[string]string)
	for key, value := range metadataFilters {
		if renamedKey, ok := metadataFilterRenameDict[key]; ok {
			rawMetadataFilters["metadata/"+renamedKey] = value
		} else {
			rawMetadataFilters["metadata/"+key] = value
		}
	}
	return rawMetadataFilters
}

// convertFiltersToBackend converts frontend filters to backend facet format
func convertFiltersToBackend(ctx context.Context, filters map[string]string, namespace string) ([]string, error) {
	ctx, span := serviceTracer.Start(ctx, "search-service:convert-filters-to-backend")
	defer span.End()

	var filterStrings []string

	// Add namespace filter if specified
	if namespace != "" {
		filterStrings = append(filterStrings, fmt.Sprintf("namespace/%s", namespace))
	}

	if len(filters) == 0 {
		return filterStrings, nil
	}

	logger.Info(ctx, "converting filters to backend facet format",
		zap.Int("filter_count", len(filters)),
		zap.Any("filters", filters),
		zap.String("namespace", namespace))

	// Simple conversion: convert each filter to facet format
	for key, value := range filters {
		if value != "" && key != "q" && key != "page" && key != "per_page" && key != "limit" && key != "namespace" {
			// Convert to facet format
			facetFilter := fmt.Sprintf("%s/%s", key, value)
			filterStrings = append(filterStrings, facetFilter)
		}
	}

	logger.Info(ctx, "filters converted to facet format",
		zap.Int("original_filter_count", len(filters)),
		zap.Int("backend_filter_count", len(filterStrings)),
		zap.Strings("backend_filters", filterStrings))

	return filterStrings, nil
}

// fallbackFilterConversion provides simple filter conversion when filter service is unavailable
func (s *SearchService) fallbackFilterConversion(filters map[string]string, namespace string) []string {
	var filterStrings []string

	// Add namespace filter if specified
	if namespace != "" {
		filterStrings = append(filterStrings, fmt.Sprintf("namespace/%s", namespace))
	}

	// Simple field:value conversion with metadata prefix
	for key, value := range filters {
		if value != "" && key != "q" && key != "page" && key != "per_page" && key != "limit" && key != "namespace" {
			// Convert to facet format: "metadata/field_name:value"
			facetFilter := fmt.Sprintf("metadata/%s:%s", key, value)
			filterStrings = append(filterStrings, facetFilter)
		}
	}

	return filterStrings
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
	ctx, span := serviceTracer.Start(ctx, "search-service:execute-search")
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

// extractTitle extracts a meaningful title from the search result
func (s *SearchService) extractTitle(result fugusdk.FuguSearchResult) string {
	// Try to extract title from metadata first
	if result.Metadata != nil {
		if title, ok := result.Metadata["title"].(string); ok && title != "" {
			return title
		}
		if fileName, ok := result.Metadata["file_name"].(string); ok && fileName != "" {
			return fileName
		}
		if subject, ok := result.Metadata["subject"].(string); ok && subject != "" {
			return subject
		}
		if name, ok := result.Metadata["name"].(string); ok && name != "" {
			return name
		}
	}

	// Fallback to truncated text content
	if len(result.Text) > 100 {
		return result.Text[:100] + "..."
	}

	// Final fallback to ID
	if result.Text != "" {
		return result.Text
	}

	return fmt.Sprintf("Document %s", result.ID)
}

// GetSearchInfo provides information about search capabilities and status
func (s *SearchService) GetSearchInfo(ctx context.Context) (*SearchInfo, error) {
	ctx, span := serviceTracer.Start(ctx, "search-service:get-search-info")
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

	// Get available filters from the simplified filter service
	var availableFilters []string
	filters, err := s.filterService.GetAllFilters(ctx)
	if err != nil {
		logger.Warn(ctx, "failed to get filters", zap.Error(err))
		availableFilters = []string{}
	} else {
		for _, filter := range filters {
			availableFilters = append(availableFilters, filter.FilterPath)
		}
	}

	// Build search info response
	info := &SearchInfo{
		Status:      "operational",
		Version:     "1.1.0",
		LastUpdated: "2024-01-01T00:00:00Z", // Use static timestamp
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
			TotalDocuments:   0,                // Would need to query fugu for this
			IndexedFields:    availableFilters, // Use filter paths as indexed fields
			AvailableFilters: availableFilters,
			BackendStatus:    backendStatus,
		},
	}

	return info, nil
}
