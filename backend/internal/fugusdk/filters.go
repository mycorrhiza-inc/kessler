package fugusdk

// Updated Go SDK methods to work with filter config system

import (
	"context"
	"fmt"
)

// FilterBuilder helps build validated filters for Fugu search
type FilterBuilder struct {
	filters   map[string]string
	namespace string
	client    *Client
}

// NewFilterBuilder creates a new filter builder with client for validation
func (c *Client) NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		filters: make(map[string]string),
		client:  c,
	}
}

// AddMetadataFilter adds a metadata field filter with optional validation
func (fb *FilterBuilder) AddMetadataFilter(field, value string) *FilterBuilder {
	fb.filters[field] = value
	return fb
}

// AddNamespaceFilter sets the namespace for filtering
func (fb *FilterBuilder) AddNamespaceFilter(namespace string) *FilterBuilder {
	fb.namespace = namespace
	return fb
}

// AddFilter adds a generic filter
func (fb *FilterBuilder) AddFilter(field, value string) *FilterBuilder {
	fb.filters[field] = value
	return fb
}

// Validate validates all filters using the filter config service
func (fb *FilterBuilder) Validate(ctx context.Context) (*FilterValidationResult, error) {
	if len(fb.filters) == 0 {
		return &FilterValidationResult{IsValid: true}, nil
	}

	validateReq := FilterValidateRequest{
		Filters: fb.filters,
	}

	resp, err := fb.client.makeRequest(ctx, "POST", "/search/filters/validate", validateReq)
	if err != nil {
		return nil, err
	}

	var result FilterValidationResult
	err = fb.client.handleResponse(resp, &result)
	return &result, err
}

// BuildWithValidation builds filters after validating them
func (fb *FilterBuilder) BuildWithValidation(ctx context.Context) ([]string, error) {
	// Validate filters first
	validation, err := fb.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("filter validation failed: %w", err)
	}

	if !validation.IsValid {
		return nil, fmt.Errorf("invalid filters: %v", validation.Errors)
	}

	// Convert to backend format
	return fb.Build(ctx)
}

// Build converts filters to backend format using filter config service
func (fb *FilterBuilder) Build(ctx context.Context) ([]string, error) {
	if len(fb.filters) == 0 {
		if fb.namespace != "" {
			return []string{fb.namespace}, nil
		}
		return []string{}, nil
	}

	// Use filter config service to convert filters
	convertReq := FilterConvertRequest{
		Filters: fb.filters,
	}

	resp, err := fb.client.makeRequest(ctx, "POST", "/search/filters/convert", convertReq)
	if err != nil {
		// Fallback to simple conversion if service unavailable
		return fb.buildFallback(), nil
	}

	var convertResp FilterConvertResponse
	if err := fb.client.handleResponse(resp, &convertResp); err != nil {
		// Fallback to simple conversion
		return fb.buildFallback(), nil
	}

	// Convert backend filters map to string slice
	var filterStrings []string

	// Add namespace first if specified
	if fb.namespace != "" {
		filterStrings = append(filterStrings, fb.namespace)
	}

	// Add converted filters
	for field, value := range convertResp.BackendFilters {
		if valueStr, ok := value.(string); ok && valueStr != "" {
			filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", field, valueStr))
		}
	}

	return filterStrings, nil
}

// buildFallback provides simple filter conversion when config service unavailable
func (fb *FilterBuilder) buildFallback() []string {
	var filterStrings []string

	// Add namespace first if specified
	if fb.namespace != "" {
		filterStrings = append(filterStrings, fb.namespace)
	}

	// Simple field:value conversion
	for field, value := range fb.filters {
		if value != "" {
			filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", field, value))
		}
	}

	return filterStrings
}

// Filter config related types
type FilterConvertRequest struct {
	Filters map[string]string `json:"filters"`
}

type FilterConvertResponse struct {
	BackendFilters map[string]interface{} `json:"backendFilters"`
}

type FilterValidateRequest struct {
	Filters map[string]string `json:"filters"`
}

type FilterValidationResult struct {
	IsValid  bool                `json:"isValid"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
}

type ValidationError struct {
	FieldID string `json:"fieldId"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ValidationWarning struct {
	FieldID string `json:"fieldId"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// GetFilterConfiguration retrieves the current filter configuration
func (c *Client) GetFilterConfiguration(ctx context.Context) (*FilterConfiguration, error) {
	resp, err := c.makeRequest(ctx, "GET", "/search/filters/configuration", nil)
	if err != nil {
		return nil, err
	}

	var config FilterConfiguration
	err = c.handleResponse(resp, &config)
	return &config, err
}

// GetDynamicFilterOptions gets dynamic options for a specific filter field
func (c *Client) GetDynamicFilterOptions(ctx context.Context, fieldID string, context map[string]string, namespace string) ([]FilterOption, error) {
	optionsReq := FilterOptionsRequest{
		FieldID:   fieldID,
		Context:   context,
		Namespace: namespace,
	}

	resp, err := c.makeRequest(ctx, "POST", "/search/filters/options", optionsReq)
	if err != nil {
		return nil, err
	}

	var optionsResp FilterOptionsResponse
	err = c.handleResponse(resp, &optionsResp)
	return optionsResp.Options, err
}

// Filter configuration types
type FilterConfiguration struct {
	Fields     []FilterFieldDefinition `json:"fields"`
	Categories []FilterCategory        `json:"categories"`
	Config     FilterGlobalConfig      `json:"config"`
}

type FilterFieldDefinition struct {
	ID           string            `json:"id"`
	BackendKey   string            `json:"backendKey"`
	DisplayName  string            `json:"displayName"`
	Description  string            `json:"description"`
	InputType    string            `json:"inputType"`
	Required     bool              `json:"required,omitempty"`
	Placeholder  string            `json:"placeholder,omitempty"`
	Order        int               `json:"order"`
	Category     string            `json:"category"`
	Validation   *FilterValidation `json:"validation,omitempty"`
	Options      []FilterOption    `json:"options,omitempty"`
	DefaultValue string            `json:"defaultValue,omitempty"`
	Enabled      bool              `json:"enabled"`
}

type FilterValidation struct {
	MinLength       *int   `json:"minLength,omitempty"`
	MaxLength       *int   `json:"maxLength,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
	Min             *int   `json:"min,omitempty"`
	Max             *int   `json:"max,omitempty"`
	CustomValidator string `json:"customValidator,omitempty"`
}

type FilterOption struct {
	Value    string                 `json:"value"`
	Label    string                 `json:"label"`
	Disabled bool                   `json:"disabled,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type FilterCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`
	Collapsible bool   `json:"collapsible,omitempty"`
}

type FilterGlobalConfig struct {
	Version         string `json:"version"`
	LastUpdated     string `json:"lastUpdated"`
	DefaultCategory string `json:"defaultCategory"`
}

type FilterOptionsRequest struct {
	FieldID   string            `json:"fieldId"`
	Context   map[string]string `json:"context,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
}

type FilterOptionsResponse struct {
	Options []FilterOption `json:"options"`
}

// Enhanced search methods using filter config system

// SearchWithValidatedFilters performs a search with validated filters
func (c *Client) SearchWithValidatedFilters(ctx context.Context, query string, filters map[string]string, namespace string, page, perPage int) (*SanitizedResponse, error) {
	// Build and validate filters
	filterBuilder := c.NewFilterBuilder()
	if namespace != "" {
		filterBuilder.AddNamespaceFilter(namespace)
	}
	for field, value := range filters {
		filterBuilder.AddMetadataFilter(field, value)
	}

	// Build with validation
	backendFilters, err := filterBuilder.BuildWithValidation(ctx)
	if err != nil {
		return nil, fmt.Errorf("filter validation failed: %w", err)
	}

	return c.SearchWithFilters(ctx, query, backendFilters, page, perPage)
}

// SearchWithFilters performs a search with pre-built filters
func (c *Client) SearchWithFilters(ctx context.Context, query string, filters []string, page, perPage int) (*SanitizedResponse, error) {
	searchQuery := FuguSearchQuery{
		Query: query,
	}

	if len(filters) > 0 {
		searchQuery.Filters = &filters
	}

	if page >= 0 || perPage > 0 {
		pagination := &Pagination{}
		if page >= 0 {
			pagination.Page = &page
		}
		if perPage > 0 {
			pagination.PerPage = &perPage
		}
		searchQuery.Page = pagination
	}

	return c.Search(ctx, searchQuery)
}

// SearchByMetadata is a convenience method for searching by metadata fields with validation
func (c *Client) SearchByMetadata(ctx context.Context, query string, metadataFilters map[string]string, page, perPage int) (*SanitizedResponse, error) {
	return c.SearchWithValidatedFilters(ctx, query, metadataFilters, "", page, perPage)
}

// SearchInNamespace searches within a specific namespace with optional metadata filters and validation
func (c *Client) SearchInNamespace(ctx context.Context, query, namespace string, metadataFilters map[string]string, page, perPage int) (*SanitizedResponse, error) {
	return c.SearchWithValidatedFilters(ctx, query, metadataFilters, namespace, page, perPage)
}

// BuildSmartFilters creates filters using the configuration system
func (c *Client) BuildSmartFilters(ctx context.Context) (*SmartFilterBuilder, error) {
	config, err := c.GetFilterConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get filter configuration: %w", err)
	}

	return &SmartFilterBuilder{
		client: c,
		config: config,
		values: make(map[string]string),
	}, nil
}

// SmartFilterBuilder provides intelligent filter building based on configuration
type SmartFilterBuilder struct {
	client    *Client
	config    *FilterConfiguration
	values    map[string]string
	namespace string
}

// SetFilter sets a filter value with automatic validation based on configuration
func (sfb *SmartFilterBuilder) SetFilter(fieldID, value string) error {
	// Find field definition
	var fieldDef *FilterFieldDefinition
	for _, field := range sfb.config.Fields {
		if field.ID == fieldID {
			fieldDef = &field
			break
		}
	}

	if fieldDef == nil {
		return fmt.Errorf("unknown filter field: %s", fieldID)
	}

	if !fieldDef.Enabled {
		return fmt.Errorf("filter field is disabled: %s", fieldID)
	}

	// Basic validation
	if fieldDef.Required && value == "" {
		return fmt.Errorf("field %s is required", fieldDef.DisplayName)
	}

	// Store the value
	sfb.values[fieldID] = value
	return nil
}

// SetNamespace sets the namespace filter
func (sfb *SmartFilterBuilder) SetNamespace(namespace string) *SmartFilterBuilder {
	sfb.namespace = namespace
	return sfb
}

// Build creates the final filter list with full validation
func (sfb *SmartFilterBuilder) Build(ctx context.Context) ([]string, error) {
	// Use the filter builder for final conversion
	filterBuilder := sfb.client.NewFilterBuilder()
	if sfb.namespace != "" {
		filterBuilder.AddNamespaceFilter(sfb.namespace)
	}
	for field, value := range sfb.values {
		filterBuilder.AddFilter(field, value)
	}

	return filterBuilder.BuildWithValidation(ctx)
}

// GetAvailableFields returns all available filter fields
func (sfb *SmartFilterBuilder) GetAvailableFields() []FilterFieldDefinition {
	var available []FilterFieldDefinition
	for _, field := range sfb.config.Fields {
		if field.Enabled {
			available = append(available, field)
		}
	}
	return available
}

// GetFieldsByCategory returns fields grouped by category
func (sfb *SmartFilterBuilder) GetFieldsByCategory() map[string][]FilterFieldDefinition {
	fieldsByCategory := make(map[string][]FilterFieldDefinition)
	for _, field := range sfb.config.Fields {
		if field.Enabled {
			fieldsByCategory[field.Category] = append(fieldsByCategory[field.Category], field)
		}
	}
	return fieldsByCategory
}
