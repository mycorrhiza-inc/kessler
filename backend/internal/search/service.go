package search

import (
	"context"
	"fmt"
	"kessler/internal/fugusdk"
	"kessler/internal/search/filter"
	"kessler/pkg/logger"
	"strings"
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
}

// NewSearchService creates a new search service
func NewSearchService(fuguServerURL string, filterService *filter.Service) *SearchService {
	return &SearchService{
		fuguServerURL: fuguServerURL,
		filterService: filterService,
	}
}

// convertFiltersToBackend converts frontend filters to backend facet format using filter service
func (s *SearchService) convertFiltersToBackend(ctx context.Context, filters map[string]string, namespace string) ([]string, error) {
	ctx, span := serviceTracer.Start(ctx, "search-service:convert-filters-to-backend")
	defer span.End()

	var filterStrings []string

	// Add namespace filter if specified
	if namespace != "" {
		filterStrings = append(filterStrings, namespace)
	}

	if len(filters) == 0 {
		return filterStrings, nil
	}

	logger.Info(ctx, "converting filters to backend facet format",
		zap.Int("filter_count", len(filters)),
		zap.Any("filters", filters),
		zap.String("namespace", namespace))

	// Use filter service to convert frontend filters to backend format
	backendFilters, err := s.filterService.ConvertFiltersToBackend(ctx, filters)
	if err != nil {
		logger.Warn(ctx, "failed to convert filters using filter service, falling back to simple conversion", zap.Error(err))
		// Fallback to simple conversion if filter service fails
		return s.fallbackFilterConversion(filters, namespace), nil
	}

	// Convert the backend filters map to facet strings that Fugu expects
	for key, value := range backendFilters {
		if value != nil {
			if strValue, ok := value.(string); ok && strValue != "" {
				// Convert to facet format: "metadata/field_name:value"
				facetFilter := fmt.Sprintf("metadata/%s:%s", key, strValue)
				filterStrings = append(filterStrings, facetFilter)
			}
		}
	}

	logger.Info(ctx, "filters converted to facet format using filter service",
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
		filterStrings = append(filterStrings, namespace)
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

// validateFilters validates filters using the filter service before processing
func (s *SearchService) validateFilters(ctx context.Context, filters map[string]string) error {
	ctx, span := serviceTracer.Start(ctx, "search-service:validate-filters")
	defer span.End()

	if len(filters) == 0 {
		return nil
	}

	// Use filter service to validate filters
	validation, err := s.filterService.ValidateFilters(ctx, filters)
	if err != nil {
		logger.Warn(ctx, "filter validation failed, proceeding without validation", zap.Error(err))
		return nil // Don't fail search if validation service is down
	}

	if !validation.IsValid {
		// Log validation errors but don't fail the search
		logger.Warn(ctx, "filter validation found issues",
			zap.Int("error_count", len(validation.Errors)),
			zap.Int("warning_count", len(validation.Warnings)))

		for _, validationError := range validation.Errors {
			logger.Warn(ctx, "filter validation error",
				zap.String("field_id", validationError.FieldID),
				zap.String("message", validationError.Message),
				zap.String("type", validationError.Type))
		}

		// For now, log warnings but allow the search to proceed
		// In a stricter implementation, you might want to return an error here
	}

	return nil
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

// ProcessSearch processes a search request with namespace support and filter validation
func (s *SearchService) ProcessSearch(ctx context.Context, query string, filters map[string]string, pagination PaginationParams, namespace string) (*SearchResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "search-service:process-search")
	defer span.End()

	startTime := time.Now()

	logger.Info(ctx, "starting search processing",
		zap.String("query", query),
		zap.String("namespace", namespace),
		zap.String("fugu_url", s.fuguServerURL))

	// Validate filters using filter service
	if err := s.validateFilters(ctx, filters); err != nil {
		logger.Error(ctx, "filter validation failed", zap.Error(err))
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	// Create fugu client with timeout
	clientCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := fugusdk.NewClient(clientCtx, s.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	logger.Info(ctx, "fugu client created successfully")

	// Convert filters to backend format with namespace using filter service
	backendFilters, err := s.convertFiltersToBackend(ctx, filters, namespace)
	if err != nil {
		logger.Warn(ctx, "failed to convert filters, proceeding with fallback", zap.Error(err))
		backendFilters = s.fallbackFilterConversion(filters, namespace)
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

	// Get available filters from filter service
	config, err := s.filterService.BuildFilterConfiguration(ctx)
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
			// Look for namespace facets (those without metadata/ prefix)
			for _, facet := range result.Facets {
				if !strings.HasPrefix(facet, "metadata/") {
					item.Namespace = facet
					break
				}
			}
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
