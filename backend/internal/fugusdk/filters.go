// filters.go
package fugusdk

// Updated Go SDK methods to work with filter config system and namespace facets

import (
	"context"
	"fmt"
	"strings"
)

// FilterBuilder helps build validated filters for Fugu search with namespace support
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

// AddNamespaceFacetFilter adds a namespace facet filter
func (fb *FilterBuilder) AddNamespaceFacetFilter(namespace, facetType, value string) *FilterBuilder {
	facetPath := fmt.Sprintf("namespace/%s/%s/%s", namespace, facetType, value)
	fb.filters[facetPath] = value
	return fb
}

// AddOrganizationFilter is a convenience method for adding organization namespace facets
func (fb *FilterBuilder) AddOrganizationFilter(namespace, organization string) *FilterBuilder {
	return fb.AddNamespaceFacetFilter(namespace, "organization", organization)
}

// AddConversationFilter is a convenience method for adding conversation namespace facets
func (fb *FilterBuilder) AddConversationFilter(namespace, conversationID string) *FilterBuilder {
	return fb.AddNamespaceFacetFilter(namespace, "conversation", conversationID)
}

// AddDataTypeFilter is a convenience method for adding data type namespace facets
func (fb *FilterBuilder) AddDataTypeFilter(namespace, dataType string) *FilterBuilder {
	return fb.AddNamespaceFacetFilter(namespace, "data", dataType)
}

// Validate validates all filters using the filter config service
func (fb *FilterBuilder) Validate(ctx context.Context) (*FilterValidationResult, error) {
	if len(fb.filters) == 0 {
		return &FilterValidationResult{IsValid: true}, nil
	}

	validateReq := FilterValidateRequest{
		Filters: fb.filters,
	}

	// Try namespace-specific validation endpoint if namespace is set
	endpoint := "/search/filters/validate"
	if fb.namespace != "" {
		endpoint = fmt.Sprintf("/search/filters/namespace/%s/validate", fb.namespace)
	}

	resp, err := fb.client.makeRequest(ctx, "POST", endpoint, validateReq)
	if err != nil {
		// Fallback to general validation endpoint
		if fb.namespace != "" {
			resp, err = fb.client.makeRequest(ctx, "POST", "/search/filters/validate", validateReq)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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

// Build converts filters to backend format using filter config service with namespace support
func (fb *FilterBuilder) Build(ctx context.Context) ([]string, error) {
	if len(fb.filters) == 0 {
		if fb.namespace != "" {
			return []string{fmt.Sprintf("namespace/%s", fb.namespace)}, nil
		}
		return []string{}, nil
	}

	// Check if we have namespace facets in filters
	hasNamespaceFacets := false
	for filterKey := range fb.filters {
		if strings.HasPrefix(filterKey, "namespace/") {
			hasNamespaceFacets = true
			break
		}
	}

	// Use namespace-aware filter conversion if available
	convertReq := FilterConvertRequest{
		Filters: fb.filters,
	}

	endpoint := "/search/filters/convert"
	if fb.namespace != "" {
		endpoint = fmt.Sprintf("/search/filters/namespace/%s/convert", fb.namespace)
	}

	resp, err := fb.client.makeRequest(ctx, "POST", endpoint, convertReq)
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

	// Add namespace first if specified and not already in filters
	if fb.namespace != "" && !hasNamespaceFacets {
		filterStrings = append(filterStrings, fmt.Sprintf("namespace/%s", fb.namespace))
	}

	// Add converted filters
	for field, value := range convertResp.BackendFilters {
		if valueStr, ok := value.(string); ok && valueStr != "" {
			// Handle namespace facet format
			if strings.HasPrefix(field, "namespace/") {
				filterStrings = append(filterStrings, field)
			} else {
				filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", field, valueStr))
			}
		} else if valueSlice, ok := value.([]string); ok {
			// Handle multiple values (for namespace facets)
			for _, v := range valueSlice {
				if strings.HasPrefix(field, "namespace/") {
					filterStrings = append(filterStrings, fmt.Sprintf("%s/%s", field, v))
				} else {
					filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", field, v))
				}
			}
		}
	}

	return filterStrings, nil
}

// buildFallback provides simple filter conversion when config service unavailable
func (fb *FilterBuilder) buildFallback() []string {
	var filterStrings []string

	// Add namespace first if specified
	if fb.namespace != "" {
		filterStrings = append(filterStrings, fmt.Sprintf("namespace/%s", fb.namespace))
	}

	// Convert filters - handle namespace facets specially
	for field, value := range fb.filters {
		if value != "" {
			if strings.HasPrefix(field, "namespace/") {
				// This is already a namespace facet path
				filterStrings = append(filterStrings, field)
			} else {
				// Regular metadata filter
				filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", field, value))
			}
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

// GetNamespaceFilterConfiguration retrieves filter configuration for a specific namespace
func (c *Client) GetNamespaceFilterConfiguration(ctx context.Context, namespace string) (*FilterConfiguration, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/search/filters/namespace/%s/configuration", namespace)
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		// Fallback to general configuration
		return c.GetFilterConfiguration(ctx)
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

	endpoint := "/search/filters/options"
	if namespace != "" {
		endpoint = fmt.Sprintf("/search/filters/namespace/%s/options", namespace)
	}

	resp, err := c.makeRequest(ctx, "POST", endpoint, optionsReq)
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

// Enhanced search methods using filter config system with namespace support

// SearchWithValidatedFilters performs a search with validated filters and namespace support
func (c *Client) SearchWithValidatedFilters(ctx context.Context, query string, filters map[string]string, namespace string, page, perPage int) (*SanitizedResponse, error) {
	// Build and validate filters using namespace-aware builder
	filterBuilder := c.NewFilterBuilder()
	if namespace != "" {
		filterBuilder.AddNamespaceFilter(namespace)
	}
	for field, value := range filters {
		// Check if this is a namespace facet filter
		if strings.HasPrefix(field, "namespace/") {
			filterBuilder.AddFilter(field, value)
		} else {
			filterBuilder.AddMetadataFilter(field, value)
		}
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

// BuildSmartFilters creates filters using the configuration system with namespace support
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

// BuildSmartFiltersForNamespace creates filters using namespace-specific configuration
func (c *Client) BuildSmartFiltersForNamespace(ctx context.Context, namespace string) (*SmartFilterBuilder, error) {
	config, err := c.GetNamespaceFilterConfiguration(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace filter configuration: %w", err)
	}

	return &SmartFilterBuilder{
		client:    c,
		config:    config,
		values:    make(map[string]string),
		namespace: namespace,
	}, nil
}

// SmartFilterBuilder provides intelligent filter building based on configuration with namespace support
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

// SetOrganization sets the organization namespace facet filter
func (sfb *SmartFilterBuilder) SetOrganization(organization string) error {
	if sfb.namespace == "" {
		return fmt.Errorf("namespace must be set before setting organization")
	}
	facetKey := fmt.Sprintf("namespace/%s/organization", sfb.namespace)
	sfb.values[facetKey] = organization
	return nil
}

// SetConversation sets the conversation namespace facet filter
func (sfb *SmartFilterBuilder) SetConversation(conversationID string) error {
	if sfb.namespace == "" {
		return fmt.Errorf("namespace must be set before setting conversation")
	}
	facetKey := fmt.Sprintf("namespace/%s/conversation", sfb.namespace)
	sfb.values[facetKey] = conversationID
	return nil
}

// SetDataType sets the data type namespace facet filter
func (sfb *SmartFilterBuilder) SetDataType(dataType string) error {
	if sfb.namespace == "" {
		return fmt.Errorf("namespace must be set before setting data type")
	}
	facetKey := fmt.Sprintf("namespace/%s/data", sfb.namespace)
	sfb.values[facetKey] = dataType
	return nil
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

// GetNamespaceFields returns fields specific to namespace facets
func (sfb *SmartFilterBuilder) GetNamespaceFields() []FilterFieldDefinition {
	var namespaceFields []FilterFieldDefinition
	for _, field := range sfb.config.Fields {
		if field.Enabled && field.Category == "namespace" {
			namespaceFields = append(namespaceFields, field)
		}
	}
	return namespaceFields
}

// Enhanced namespace-aware filter builder methods

// CreateNamespaceAwareFilterBuilder creates a filter builder with namespace context
func (c *Client) CreateNamespaceAwareFilterBuilder(ctx context.Context, namespace string) (*NamespaceAwareFilterBuilder, error) {
	// Get namespace info to populate available options
	namespaceInfo, err := c.GetNamespaceInfo(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace info: %w", err)
	}

	// Get filter configuration for this namespace
	config, err := c.GetNamespaceFilterConfiguration(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace filter configuration: %w", err)
	}

	return &NamespaceAwareFilterBuilder{
		client:        c,
		namespace:     namespace,
		namespaceInfo: namespaceInfo,
		config:        config,
		filters:       make(map[string]string),
	}, nil
}

// NamespaceAwareFilterBuilder provides namespace-specific filter building capabilities
type NamespaceAwareFilterBuilder struct {
	client        *Client
	namespace     string
	namespaceInfo *NamespaceInfo
	config        *FilterConfiguration
	filters       map[string]string
}

// SetOrganizationFilter adds an organization filter with validation
func (nafb *NamespaceAwareFilterBuilder) SetOrganizationFilter(organization string) error {
	// Validate that the organization exists in this namespace
	found := false
	for _, org := range nafb.namespaceInfo.Organizations {
		if org == organization {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("organization '%s' not found in namespace '%s'", organization, nafb.namespace)
	}

	nafb.filters["organization"] = organization
	return nil
}

// SetConversationFilter adds a conversation filter with validation
func (nafb *NamespaceAwareFilterBuilder) SetConversationFilter(conversationID string) error {
	// Validate that the conversation exists in this namespace
	found := false
	for _, conv := range nafb.namespaceInfo.Conversations {
		if conv == conversationID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("conversation '%s' not found in namespace '%s'", conversationID, nafb.namespace)
	}

	nafb.filters["conversation"] = conversationID
	return nil
}

// SetDataTypeFilter adds a data type filter with validation
func (nafb *NamespaceAwareFilterBuilder) SetDataTypeFilter(dataType string) error {
	// Validate that the data type exists in this namespace
	found := false
	for _, dt := range nafb.namespaceInfo.DataTypes {
		if dt == dataType {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("data type '%s' not found in namespace '%s'", dataType, nafb.namespace)
	}

	nafb.filters["data_type"] = dataType
	return nil
}

// SetMetadataFilter adds a metadata filter
func (nafb *NamespaceAwareFilterBuilder) SetMetadataFilter(field, value string) *NamespaceAwareFilterBuilder {
	nafb.filters[field] = value
	return nafb
}

// GetAvailableOrganizations returns all organizations in this namespace
func (nafb *NamespaceAwareFilterBuilder) GetAvailableOrganizations() []string {
	return nafb.namespaceInfo.Organizations
}

// GetAvailableConversations returns all conversations in this namespace
func (nafb *NamespaceAwareFilterBuilder) GetAvailableConversations() []string {
	return nafb.namespaceInfo.Conversations
}

// GetAvailableDataTypes returns all data types in this namespace
func (nafb *NamespaceAwareFilterBuilder) GetAvailableDataTypes() []string {
	return nafb.namespaceInfo.DataTypes
}

// Build creates the final validated filter list
func (nafb *NamespaceAwareFilterBuilder) Build(ctx context.Context) ([]string, error) {
	// Convert to namespace facet builder for final processing
	facetBuilder := nafb.client.NewNamespaceFacetBuilder(nafb.namespace)

	// Add namespace-specific filters
	for field, value := range nafb.filters {
		switch field {
		case "organization":
			facetBuilder.AddOrganizationFilter(value)
		case "conversation":
			facetBuilder.AddConversationFilter(value)
		case "data_type":
			facetBuilder.AddDataTypeFilter(value)
		default:
			// Add as custom facet for metadata filters
			facetBuilder.AddCustomFacet(fmt.Sprintf("metadata/%s/%s", field, value))
		}
	}

	return facetBuilder.GetFilters(), nil
}

// Search performs a search with the built namespace-aware filters
func (nafb *NamespaceAwareFilterBuilder) Search(ctx context.Context, query string, page, perPage int) (*SanitizedResponse, error) {
	filters, err := nafb.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build filters: %w", err)
	}

	return nafb.client.SearchWithNamespaceFacets(ctx, query, filters, page, perPage)
}

// Utility methods for namespace filter integration

// ValidateNamespaceFilters validates that all namespace facet filters are valid for the given namespace
func (c *Client) ValidateNamespaceFilters(ctx context.Context, namespace string, filters map[string]string) error {
	namespaceInfo, err := c.GetNamespaceInfo(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to get namespace info: %w", err)
	}

	for field, value := range filters {
		switch {
		case strings.HasSuffix(field, "/organization"):
			found := false
			for _, org := range namespaceInfo.Organizations {
				if org == value {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("organization '%s' not found in namespace '%s'", value, namespace)
			}
		case strings.HasSuffix(field, "/conversation"):
			found := false
			for _, conv := range namespaceInfo.Conversations {
				if conv == value {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("conversation '%s' not found in namespace '%s'", value, namespace)
			}
		case strings.HasSuffix(field, "/data"):
			found := false
			for _, dt := range namespaceInfo.DataTypes {
				if dt == value {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("data type '%s' not found in namespace '%s'", value, namespace)
			}
		}
	}

	return nil
}
