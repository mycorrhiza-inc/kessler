// service.go
package filter

import (
	"context"
	"fmt"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var serviceTracer = otel.Tracer("filter-service")

// Frontend filter types
type FilterInputType string

const (
	FilterInputText        FilterInputType = "text"
	FilterInputSelect      FilterInputType = "select"
	FilterInputMultiSelect FilterInputType = "multiselect"
	FilterInputDate        FilterInputType = "date"
	FilterInputDateRange   FilterInputType = "daterange"
	FilterInputNumber      FilterInputType = "number"
	FilterInputBoolean     FilterInputType = "boolean"
	FilterInputUUID        FilterInputType = "uuid"
	FilterInputCustom      FilterInputType = "custom"
	FilterInputHidden      FilterInputType = "hidden"
)

// FilterFieldDefinition represents a filter field configuration for the frontend
type FilterFieldDefinition struct {
	ID           string            `json:"id"`
	BackendKey   string            `json:"backendKey"`
	DisplayName  string            `json:"displayName"`
	Description  string            `json:"description"`
	InputType    FilterInputType   `json:"inputType"`
	Required     bool              `json:"required,omitempty"`
	Placeholder  string            `json:"placeholder,omitempty"`
	Order        int               `json:"order"`
	Category     string            `json:"category"`
	Validation   *FilterValidation `json:"validation,omitempty"`
	Options      []FilterOption    `json:"options,omitempty"`
	DefaultValue string            `json:"defaultValue,omitempty"`
	Enabled      bool              `json:"enabled"`
}

// FilterValidation represents validation rules for filter fields
type FilterValidation struct {
	MinLength       *int   `json:"minLength,omitempty"`
	MaxLength       *int   `json:"maxLength,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
	Min             *int   `json:"min,omitempty"`
	Max             *int   `json:"max,omitempty"`
	CustomValidator string `json:"customValidator,omitempty"`
}

// FilterOption represents an option for select/multiselect inputs
type FilterOption struct {
	Value    string                 `json:"value"`
	Label    string                 `json:"label"`
	Disabled bool                   `json:"disabled,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// FilterCategory represents a category grouping for filters
type FilterCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`
	Collapsible bool   `json:"collapsible,omitempty"`
}

// FilterConfiguration represents the complete filter configuration
type FilterConfiguration struct {
	Fields     []FilterFieldDefinition `json:"fields"`
	Categories []FilterCategory        `json:"categories"`
	Config     FilterGlobalConfig      `json:"config"`
}

// FilterGlobalConfig represents global configuration settings
type FilterGlobalConfig struct {
	Version         string `json:"version"`
	LastUpdated     string `json:"lastUpdated"`
	DefaultCategory string `json:"defaultCategory"`
}

// ConvertFiltersRequest represents a request to convert frontend filters to backend format
type ConvertFiltersRequest struct {
	Filters map[string]string `json:"filters"`
}

// ConvertFiltersResponse represents the response from filter conversion
type ConvertFiltersResponse struct {
	BackendFilters map[string]interface{} `json:"backendFilters"`
}

// ValidateFiltersRequest represents a request to validate filters
type ValidateFiltersRequest struct {
	Filters map[string]string `json:"filters"`
}

// ValidationError represents a validation error
type ValidationError struct {
	FieldID string `json:"fieldId"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	FieldID string `json:"fieldId"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ValidateFiltersResponse represents the response from filter validation
type ValidateFiltersResponse struct {
	IsValid  bool                `json:"isValid"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
}

// GetOptionsRequest represents a request for dynamic filter options
type GetOptionsRequest struct {
	FieldID   string            `json:"fieldId"`
	Context   map[string]string `json:"context,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
}

// GetOptionsResponse represents the response for dynamic filter options
type GetOptionsResponse struct {
	Options []FilterOption `json:"options"`
}

// FacetInfo represents facet information from fugu
type FacetInfo struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}

// Service handles the business logic for filter configuration
type Service struct {
	fuguServerURL string
}

// NewService creates a new filter configuration service
func NewService(fuguServerURL string) *Service {
	return &Service{
		fuguServerURL: fuguServerURL,
	}
}

// BuildFilterConfiguration builds the complete filter configuration from fugu facets
func (s *Service) BuildFilterConfiguration(ctx context.Context) (*FilterConfiguration, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:build-configuration")
	defer span.End()

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Get all filters (facets) from fugu
	response, err := client.ListFilters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get filters from fugu: %w", err)
	}

	// Parse facets from response
	facets, err := s.parseFacetsFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse facets: %w", err)
	}

	// Convert facets to filter fields
	fields := s.convertFacetsToFields(facets)

	// Create categories
	categories := s.createFilterCategories()

	// Build configuration
	config := &FilterConfiguration{
		Fields:     fields,
		Categories: categories,
		Config: FilterGlobalConfig{
			Version:         "1.0.0",
			LastUpdated:     time.Now().Format(time.RFC3339),
			DefaultCategory: "general",
		},
	}

	return config, nil
}

// ConvertFiltersToBackend converts frontend filters to backend format
func (s *Service) ConvertFiltersToBackend(ctx context.Context, filters map[string]string) (map[string]interface{}, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:convert-to-backend")
	defer span.End()

	// Get current configuration to map field IDs to backend keys
	config, err := s.BuildFilterConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Create field ID to backend key mapping
	fieldMap := make(map[string]string)
	for _, field := range config.Fields {
		fieldMap[field.ID] = field.BackendKey
	}

	// Convert filters
	backendFilters := make(map[string]interface{})
	for fieldID, value := range filters {
		if value == "" {
			continue
		}

		backendKey, exists := fieldMap[fieldID]
		if !exists {
			logger.Warn(ctx, "unknown filter field", zap.String("field_id", fieldID))
			continue
		}

		backendFilters[backendKey] = value
	}

	return backendFilters, nil
}

// ValidateFilters validates filter values against the current configuration
func (s *Service) ValidateFilters(ctx context.Context, filters map[string]string) (*ValidateFiltersResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:validate-filters")
	defer span.End()

	// Get current configuration
	config, err := s.BuildFilterConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Create field ID to field mapping
	fieldMap := make(map[string]FilterFieldDefinition)
	for _, field := range config.Fields {
		fieldMap[field.ID] = field
	}

	var errors []ValidationError
	var warnings []ValidationWarning

	// Validate each filter
	for fieldID, value := range filters {
		field, exists := fieldMap[fieldID]
		if !exists {
			warnings = append(warnings, ValidationWarning{
				FieldID: fieldID,
				Message: "Unknown filter field",
				Type:    "unknown_field",
			})
			continue
		}

		if !field.Enabled {
			warnings = append(warnings, ValidationWarning{
				FieldID: fieldID,
				Message: "Filter field is disabled",
				Type:    "disabled_field",
			})
			continue
		}

		// Validate field value
		fieldErrors := s.validateFieldValue(field, value)
		errors = append(errors, fieldErrors...)
	}

	return &ValidateFiltersResponse{
		IsValid:  len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// GetDynamicOptions gets dynamic options for a specific field
func (s *Service) GetDynamicOptions(ctx context.Context, fieldID string, context map[string]string, namespace string) ([]FilterOption, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:get-dynamic-options")
	defer span.End()

	// Create fugu client to fetch dynamic data
	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Get current filter configuration to understand the field
	var config *FilterConfiguration
	if namespace != "" {
		config, err = s.BuildNamespaceFilterConfiguration(ctx, namespace)
	} else {
		config, err = s.BuildFilterConfiguration(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Find the field definition
	var field *FilterFieldDefinition
	for _, f := range config.Fields {
		if f.ID == fieldID {
			field = &f
			break
		}
	}

	if field == nil {
		return nil, fmt.Errorf("field not found: %s", fieldID)
	}

	// If field already has static options, return them
	if len(field.Options) > 0 {
		return field.Options, nil
	}

	// For dynamic options, we need to query the actual data
	options, err := s.fetchDynamicOptionsFromFugu(ctx, client, field, context, namespace)
	if err != nil {
		logger.Warn(ctx, "failed to fetch dynamic options, returning empty list",
			zap.String("field_id", fieldID),
			zap.Error(err))
		return []FilterOption{}, nil
	}

	return options, nil
}

// parseFacetsFromResponse extracts facet information from fugu response
func (s *Service) parseFacetsFromResponse(response *fugusdk.SanitizedResponse) ([]FacetInfo, error) {
	// Parse the actual fugu response format: {"filters":[["/metadata/field_name",count],...]}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	filtersData, ok := data["filters"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("filters not found in response")
	}

	facets := make([]FacetInfo, 0, len(filtersData))
	for _, filter := range filtersData {
		if filterArray, ok := filter.([]interface{}); ok && len(filterArray) == 2 {
			// Extract path and count from ["/metadata/field_name", count] format
			if path, ok := filterArray[0].(string); ok {
				count := int64(0)
				if countFloat, ok := filterArray[1].(float64); ok {
					count = int64(countFloat)
				} else if countInt, ok := filterArray[1].(int64); ok {
					count = countInt
				}

				facets = append(facets, FacetInfo{
					Path:  path,
					Count: count,
				})
			}
		}
	}

	return facets, nil
}

// convertFacetsToFields converts fugu facets to frontend filter field definitions
func (s *Service) convertFacetsToFields(facets []FacetInfo) []FilterFieldDefinition {
	fields := make([]FilterFieldDefinition, 0, len(facets))

	for i, facet := range facets {
		// Extract filter field from facet path (last segment after /metadata/)
		fieldName := s.extractFieldNameFromFacet(facet.Path)
		if fieldName == "" {
			continue
		}

		// Determine input type based on field name patterns
		inputType := s.determineInputType(fieldName)

		// Create field definition
		field := FilterFieldDefinition{
			ID:           s.generateFieldID(fieldName),
			BackendKey:   fieldName,
			DisplayName:  s.generateDisplayName(fieldName),
			Description:  s.generateDescription(fieldName, facet.Count),
			InputType:    inputType,
			Required:     false,
			Placeholder:  s.generatePlaceholder(fieldName, inputType),
			Order:        i,
			Category:     s.determineCategory(fieldName),
			Validation:   s.createValidation(fieldName, inputType),
			Options:      s.createOptions(fieldName, inputType),
			DefaultValue: "",
			Enabled:      true,
		}

		fields = append(fields, field)
	}

	// Sort fields by order
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Order < fields[j].Order
	})

	return fields
}

// extractFieldNameFromFacet extracts the filter field name from a facet path
func (s *Service) extractFieldNameFromFacet(facetPath string) string {
	// For paths like "/metadata/case_number", extract "case_number"
	// Skip namespace paths like "/namespace/NYPUC/organization"
	parts := strings.Split(strings.Trim(facetPath, "/"), "/")

	if len(parts) < 2 {
		return ""
	}

	// Skip namespace facets - handled separately in namespace_service.go
	if parts[0] == "namespace" {
		return ""
	}

	if parts[0] != "metadata" {
		return ""
	}

	return parts[len(parts)-1] // Get the last segment after /metadata/
}

// determineInputType determines the appropriate input type for a field
func (s *Service) determineInputType(fieldName string) FilterInputType {
	lower := strings.ToLower(fieldName)

	switch {
	case strings.Contains(lower, "date"):
		return FilterInputDate
	case strings.Contains(lower, "time"):
		return FilterInputDate
	case strings.Contains(lower, "id"):
		return FilterInputUUID
	case strings.Contains(lower, "uuid"):
		return FilterInputUUID
	case strings.Contains(lower, "type"), strings.Contains(lower, "status"), strings.Contains(lower, "category"):
		return FilterInputSelect
	case strings.Contains(lower, "tag"), strings.Contains(lower, "label"):
		return FilterInputMultiSelect
	case strings.Contains(lower, "count"), strings.Contains(lower, "size"), strings.Contains(lower, "length"):
		return FilterInputNumber
	case strings.Contains(lower, "enable"), strings.Contains(lower, "active"), strings.Contains(lower, "flag"):
		return FilterInputBoolean
	default:
		return FilterInputText
	}
}

// generateFieldID creates a consistent field ID from field name
func (s *Service) generateFieldID(fieldName string) string {
	// Convert to snake_case and ensure uniqueness
	id := strings.ToLower(fieldName)
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	return id
}

// generateDisplayName creates a human-readable display name
func (s *Service) generateDisplayName(fieldName string) string {
	// Convert snake_case or camelCase to Title Case
	words := strings.FieldsFunc(fieldName, func(c rune) bool {
		return c == '_' || c == '-'
	})

	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}

	return strings.Join(words, " ")
}

// generateDescription creates a description for the field
func (s *Service) generateDescription(fieldName string, count int64) string {
	displayName := s.generateDisplayName(fieldName)
	return fmt.Sprintf("Filter by %s (%d records)", strings.ToLower(displayName), count)
}

// generatePlaceholder creates an appropriate placeholder for the field
func (s *Service) generatePlaceholder(fieldName string, inputType FilterInputType) string {
	displayName := s.generateDisplayName(fieldName)

	switch inputType {
	case FilterInputDate:
		return "Select date..."
	case FilterInputDateRange:
		return "Select date range..."
	case FilterInputSelect:
		return fmt.Sprintf("Select %s...", strings.ToLower(displayName))
	case FilterInputMultiSelect:
		return fmt.Sprintf("Select %s...", strings.ToLower(displayName))
	case FilterInputNumber:
		return "Enter number..."
	case FilterInputUUID:
		return "Enter ID..."
	default:
		return fmt.Sprintf("Enter %s...", strings.ToLower(displayName))
	}
}

// determineCategory determines the category for a field
func (s *Service) determineCategory(fieldName string) string {
	lower := strings.ToLower(fieldName)

	switch {
	case strings.Contains(lower, "date"), strings.Contains(lower, "time"), strings.Contains(lower, "created"), strings.Contains(lower, "updated"):
		return "temporal"
	case strings.Contains(lower, "type"), strings.Contains(lower, "category"), strings.Contains(lower, "status"):
		return "classification"
	case strings.Contains(lower, "tag"), strings.Contains(lower, "label"), strings.Contains(lower, "keyword"):
		return "content"
	case strings.Contains(lower, "id"), strings.Contains(lower, "uuid"), strings.Contains(lower, "reference"):
		return "identifiers"
	case strings.Contains(lower, "source"), strings.Contains(lower, "origin"), strings.Contains(lower, "namespace"):
		return "source"
	default:
		return "general"
	}
}

// createValidation creates validation rules for a field
func (s *Service) createValidation(fieldName string, inputType FilterInputType) *FilterValidation {
	validation := &FilterValidation{}

	switch inputType {
	case FilterInputText:
		maxLen := 255
		validation.MaxLength = &maxLen
	case FilterInputUUID:
		validation.Pattern = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	case FilterInputNumber:
		min := 0
		validation.Min = &min
	}

	// Return nil if no validation rules were set
	if validation.MinLength == nil && validation.MaxLength == nil && validation.Pattern == "" && validation.Min == nil && validation.Max == nil {
		return nil
	}

	return validation
}

// createOptions creates options for select/multiselect fields
func (s *Service) createOptions(fieldName string, inputType FilterInputType) []FilterOption {
	if inputType != FilterInputSelect && inputType != FilterInputMultiSelect {
		return nil
	}

	// For now, return empty slice - options should be loaded dynamically
	return []FilterOption{}
}

// createFilterCategories creates the predefined filter categories
func (s *Service) createFilterCategories() []FilterCategory {
	return []FilterCategory{
		{
			ID:          "general",
			Name:        "General",
			Description: "General purpose filters",
			Order:       0,
			Collapsible: true,
		},
		{
			ID:          "temporal",
			Name:        "Time & Date",
			Description: "Date and time related filters",
			Order:       1,
			Collapsible: true,
		},
		{
			ID:          "classification",
			Name:        "Classification",
			Description: "Type and category filters",
			Order:       2,
			Collapsible: true,
		},
		{
			ID:          "content",
			Name:        "Content",
			Description: "Content-related filters",
			Order:       3,
			Collapsible: true,
		},
		{
			ID:          "identifiers",
			Name:        "Identifiers",
			Description: "ID and reference filters",
			Order:       4,
			Collapsible: true,
		},
		{
			ID:          "source",
			Name:        "Source",
			Description: "Source and origin filters",
			Order:       5,
			Collapsible: true,
		},
	}
}

// validateFieldValue validates a single field value
func (s *Service) validateFieldValue(field FilterFieldDefinition, value string) []ValidationError {
	var errors []ValidationError

	// Required field validation
	if field.Required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			FieldID: field.ID,
			Message: fmt.Sprintf("%s is required", field.DisplayName),
			Type:    "required",
		})
		return errors
	}

	// Skip validation for empty optional fields
	if strings.TrimSpace(value) == "" {
		return errors
	}

	// Validation based on field configuration
	if field.Validation != nil {
		validation := field.Validation

		// Length validation
		if validation.MinLength != nil && len(value) < *validation.MinLength {
			errors = append(errors, ValidationError{
				FieldID: field.ID,
				Message: fmt.Sprintf("%s must be at least %d characters", field.DisplayName, *validation.MinLength),
				Type:    "minLength",
			})
		}

		if validation.MaxLength != nil && len(value) > *validation.MaxLength {
			errors = append(errors, ValidationError{
				FieldID: field.ID,
				Message: fmt.Sprintf("%s must be no more than %d characters", field.DisplayName, *validation.MaxLength),
				Type:    "maxLength",
			})
		}

		// Pattern validation using regex
		if validation.Pattern != "" {
			if matched, err := regexp.MatchString(validation.Pattern, value); err != nil {
				errors = append(errors, ValidationError{
					FieldID: field.ID,
					Message: fmt.Sprintf("Invalid pattern for %s", field.DisplayName),
					Type:    "pattern",
				})
			} else if !matched {
				errors = append(errors, ValidationError{
					FieldID: field.ID,
					Message: fmt.Sprintf("%s format is invalid", field.DisplayName),
					Type:    "pattern",
				})
			}
		}

		// Number validation for number fields
		if field.InputType == FilterInputNumber {
			if num, err := strconv.ParseFloat(value, 64); err != nil {
				errors = append(errors, ValidationError{
					FieldID: field.ID,
					Message: fmt.Sprintf("%s must be a valid number", field.DisplayName),
					Type:    "number",
				})
			} else {
				// Min/max validation for numbers
				if validation.Min != nil && int(num) < *validation.Min {
					errors = append(errors, ValidationError{
						FieldID: field.ID,
						Message: fmt.Sprintf("%s must be at least %d", field.DisplayName, *validation.Min),
						Type:    "min",
					})
				}
				if validation.Max != nil && int(num) > *validation.Max {
					errors = append(errors, ValidationError{
						FieldID: field.ID,
						Message: fmt.Sprintf("%s must be no more than %d", field.DisplayName, *validation.Max),
						Type:    "max",
					})
				}
			}
		}

		// Boolean validation
		if field.InputType == FilterInputBoolean {
			if value != "true" && value != "false" && value != "1" && value != "0" {
				errors = append(errors, ValidationError{
					FieldID: field.ID,
					Message: fmt.Sprintf("%s must be true, false, 1, or 0", field.DisplayName),
					Type:    "boolean",
				})
			}
		}

		// Date validation
		if field.InputType == FilterInputDate {
			if _, err := time.Parse("2006-01-02", value); err != nil {
				// Try alternative formats
				if _, err2 := time.Parse(time.RFC3339, value); err2 != nil {
					if _, err3 := time.Parse("2006-01-02T15:04:05", value); err3 != nil {
						errors = append(errors, ValidationError{
							FieldID: field.ID,
							Message: fmt.Sprintf("%s must be a valid date (YYYY-MM-DD format)", field.DisplayName),
							Type:    "date",
						})
					}
				}
			}
		}
	}

	// Input type specific validation
	switch field.InputType {
	case FilterInputUUID:
		// UUID validation (if no pattern is specified)
		if field.Validation == nil || field.Validation.Pattern == "" {
			uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
			if matched, _ := regexp.MatchString(uuidPattern, value); !matched {
				errors = append(errors, ValidationError{
					FieldID: field.ID,
					Message: fmt.Sprintf("%s must be a valid UUID", field.DisplayName),
					Type:    "uuid",
				})
			}
		}
	case FilterInputSelect:
		// Validate against available options if they exist
		if len(field.Options) > 0 {
			valid := false
			for _, option := range field.Options {
				if option.Value == value && !option.Disabled {
					valid = true
					break
				}
			}
			if !valid {
				errors = append(errors, ValidationError{
					FieldID: field.ID,
					Message: fmt.Sprintf("Invalid option for %s", field.DisplayName),
					Type:    "option",
				})
			}
		}
	case FilterInputMultiSelect:
		// For multi-select, split by comma and validate each option
		if len(field.Options) > 0 {
			values := strings.Split(value, ",")
			for _, v := range values {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				valid := false
				for _, option := range field.Options {
					if option.Value == v && !option.Disabled {
						valid = true
						break
					}
				}
				if !valid {
					errors = append(errors, ValidationError{
						FieldID: field.ID,
						Message: fmt.Sprintf("Invalid option '%s' for %s", v, field.DisplayName),
						Type:    "option",
					})
					break
				}
			}
		}
	}

	return errors
}

// fetchDynamicOptionsFromFugu fetches dynamic options from Fugu based on field configuration
func (s *Service) fetchDynamicOptionsFromFugu(ctx context.Context, client *fugusdk.Client, field *FilterFieldDefinition, context map[string]string, namespace string) ([]FilterOption, error) {
	switch field.InputType {
	case FilterInputSelect, FilterInputMultiSelect:
		return s.fetchSelectOptions(ctx, client, field, namespace)
	case FilterInputBoolean:
		return s.getBooleanOptions(), nil
	default:
		// For text, number, date fields, we typically don't provide dynamic options
		return []FilterOption{}, nil
	}
}

// fetchSelectOptions fetches select options for a field by querying distinct values
func (s *Service) fetchSelectOptions(ctx context.Context, client *fugusdk.Client, field *FilterFieldDefinition, namespace string) ([]FilterOption, error) {
	// Query for documents to extract unique values for this field
	searchQuery := fugusdk.FuguSearchQuery{
		Query: "*", // Wildcard to get all documents
	}

	// Add namespace filter if specified
	if namespace != "" {
		filters := []string{fmt.Sprintf("namespace/%s", namespace)}
		searchQuery.Filters = &filters
	}

	// Limit results to avoid too much data
	page := 0
	perPage := 100
	searchQuery.Page = &fugusdk.Pagination{
		Page:    &page,
		PerPage: &perPage,
	}

	response, err := client.Search(ctx, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search for dynamic options: %w", err)
	}

	// Extract unique values for this field from the results
	uniqueValues := make(map[string]int) // value -> count

	for _, result := range response.Results {
		if result.Metadata != nil {
			if value, exists := result.Metadata[field.BackendKey]; exists {
				if strValue, ok := value.(string); ok && strValue != "" {
					uniqueValues[strValue]++
				}
			}
		}
	}

	// Convert to FilterOption format
	options := make([]FilterOption, 0, len(uniqueValues))
	for value, count := range uniqueValues {
		options = append(options, FilterOption{
			Value: value,
			Label: s.formatOptionLabel(value, field.InputType),
			Metadata: map[string]interface{}{
				"count": count,
			},
		})
	}

	// Sort options by frequency (most common first)
	sort.Slice(options, func(i, j int) bool {
		countI := options[i].Metadata["count"].(int)
		countJ := options[j].Metadata["count"].(int)
		if countI != countJ {
			return countI > countJ
		}
		// If counts are equal, sort alphabetically
		return options[i].Label < options[j].Label
	})

	// Limit to reasonable number of options
	if len(options) > 50 {
		options = options[:50]
	}

	return options, nil
}

// getBooleanOptions returns standard boolean options
func (s *Service) getBooleanOptions() []FilterOption {
	return []FilterOption{
		{
			Value: "true",
			Label: "Yes",
		},
		{
			Value: "false",
			Label: "No",
		},
	}
}

// formatOptionLabel formats the display label for an option value
func (s *Service) formatOptionLabel(value string, inputType FilterInputType) string {
	switch inputType {
	case FilterInputBoolean:
		switch strings.ToLower(value) {
		case "true", "1", "yes", "y":
			return "Yes"
		case "false", "0", "no", "n":
			return "No"
		}
	case FilterInputSelect, FilterInputMultiSelect:
		// Convert snake_case or kebab-case to Title Case
		words := strings.FieldsFunc(value, func(c rune) bool {
			return c == '_' || c == '-' || c == '.'
		})

		for i, word := range words {
			words[i] = strings.Title(strings.ToLower(word))
		}

		if len(words) > 0 {
			return strings.Join(words, " ")
		}
	}

	// Default: return value as-is but with proper capitalization
	if len(value) > 0 {
		return strings.ToUpper(string(value[0])) + strings.ToLower(value[1:])
	}

	return value
}

// BuildNamespaceFilterConfiguration is a convenience method that delegates to NamespaceService
// This allows the base Service to work with namespace functionality when needed
func (s *Service) BuildNamespaceFilterConfiguration(ctx context.Context, namespace string) (*FilterConfiguration, error) {
	// Create a namespace service instance and delegate
	nsService := &NamespaceService{Service: s}
	return nsService.BuildNamespaceFilterConfiguration(ctx, namespace)
}
