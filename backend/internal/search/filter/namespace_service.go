// namespace_service.go
package filter

import (
	"context"
	"fmt"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"strings"
	"time"

	"go.uber.org/zap"
)

// NamespaceService handles namespace-specific filter configuration and operations
type NamespaceService struct {
	*Service // Embed the base service
}

// NewNamespaceService creates a new namespace filter service
func NewNamespaceService(fuguServerURL string) *NamespaceService {
	return &NamespaceService{
		Service: NewService(fuguServerURL),
	}
}

// BuildNamespaceFilterConfiguration builds filter configuration for a specific namespace
func (ns *NamespaceService) BuildNamespaceFilterConfiguration(ctx context.Context, namespace string) (*FilterConfiguration, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:build-configuration")
	defer span.End()

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Get namespace-specific information
	namespaceInfo, err := client.GetNamespaceInfo(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace info: %w", err)
	}

	// Build namespace-specific filter fields
	fields := ns.buildNamespaceFilterFields(namespaceInfo)

	// Get regular facets and add them
	response, err := client.GetNamespaceFilters(ctx, namespace)
	if err != nil {
		logger.Warn(ctx, "failed to get namespace filters, continuing with basic fields", zap.Error(err))
	} else {
		regularFields := ns.buildRegularFilterFields(response, namespace)
		fields = append(fields, regularFields...)
	}

	// Create categories including namespace-specific ones
	categories := ns.createNamespaceFilterCategories()

	// Build configuration
	config := &FilterConfiguration{
		Fields:     fields,
		Categories: categories,
		Config: FilterGlobalConfig{
			Version:         "1.0.0-namespace",
			LastUpdated:     time.Now().Format(time.RFC3339),
			DefaultCategory: "namespace",
		},
	}

	return config, nil
}

// ConvertNamespaceFiltersToBackend converts namespace-aware filters to backend format
func (ns *NamespaceService) ConvertNamespaceFiltersToBackend(ctx context.Context, filters map[string]string, namespace string) ([]string, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:convert-filters")
	defer span.End()

	var backendFilters []string

	// Create fugu client for dynamic conversion
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Use the namespace facet builder for conversion
	builder := client.NewNamespaceFacetBuilder(namespace)

	for fieldID, value := range filters {
		if value == "" {
			continue
		}

		switch fieldID {
		case "organization":
			builder.AddOrganizationFilter(value)
		case "conversation":
			builder.AddConversationFilter(value)
		case "data_type":
			builder.AddDataTypeFilter(value)
		case "namespace":
			// Namespace is already set in builder
			continue
		default:
			// Handle metadata filters
			metadataFacet := ns.convertMetadataFilter(fieldID, value, namespace)
			if metadataFacet != "" {
				builder.AddCustomFacet(metadataFacet)
			}
		}
	}

	// Build the filters
	builtFilters := builder.GetFilters()
	backendFilters = append(backendFilters, builtFilters...)

	return backendFilters, nil
}

// ValidateNamespaceFilters validates filters for a specific namespace
func (ns *NamespaceService) ValidateNamespaceFilters(ctx context.Context, filters map[string]string, namespace string) (*ValidateFiltersResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:validate-filters")
	defer span.End()

	// Get namespace-specific configuration
	config, err := ns.BuildNamespaceFilterConfiguration(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace configuration: %w", err)
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
				Message: "Unknown filter field for namespace",
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

		// Additional namespace-specific validation
		if err := ns.validateNamespaceFieldValue(ctx, field, value, namespace); err != nil {
			errors = append(errors, ValidationError{
				FieldID: fieldID,
				Message: err.Error(),
				Type:    "namespace_validation",
			})
			continue
		}

		// Regular field validation
		fieldErrors := ns.validateFieldValue(field, value)
		errors = append(errors, fieldErrors...)
	}

	return &ValidateFiltersResponse{
		IsValid:  len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// GetNamespaceDynamicOptions gets dynamic options for namespace-specific fields
func (ns *NamespaceService) GetNamespaceDynamicOptions(ctx context.Context, fieldID string, namespace string, context map[string]string) ([]FilterOption, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:get-dynamic-options")
	defer span.End()

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Handle namespace-specific fields directly
	switch fieldID {
	case "organization":
		orgs, err := client.GetNamespaceOrganizations(ctx, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations: %w", err)
		}

		options := make([]FilterOption, len(orgs))
		for i, org := range orgs {
			options[i] = FilterOption{
				Value: org,
				Label: ns.formatOptionLabel(org, FilterInputSelect),
			}
		}
		return options, nil

	case "conversation":
		convs, err := client.GetNamespaceConversations(ctx, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversations: %w", err)
		}

		options := make([]FilterOption, len(convs))
		for i, conv := range convs {
			options[i] = FilterOption{
				Value: conv,
				Label: ns.formatConversationLabel(conv),
			}
		}
		return options, nil

	case "data_type":
		dataTypes, err := client.GetNamespaceDataTypes(ctx, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get data types: %w", err)
		}

		options := make([]FilterOption, len(dataTypes))
		for i, dt := range dataTypes {
			options[i] = FilterOption{
				Value: dt,
				Label: ns.formatDataTypeLabel(dt),
			}
		}
		return options, nil

	default:
		// Fall back to regular dynamic options
		return ns.GetDynamicOptions(ctx, fieldID, context, namespace)
	}
}

// buildNamespaceFilterFields creates filter fields for namespace facets
func (ns *NamespaceService) buildNamespaceFilterFields(namespaceInfo *fugusdk.NamespaceInfo) []FilterFieldDefinition {
	var fields []FilterFieldDefinition
	order := 0

	// Add namespace field itself
	fields = append(fields, FilterFieldDefinition{
		ID:           "namespace",
		BackendKey:   "namespace",
		DisplayName:  "Namespace",
		Description:  fmt.Sprintf("Filter by namespace: %s", namespaceInfo.Name),
		InputType:    FilterInputText,
		Required:     false,
		Placeholder:  "Enter namespace...",
		Order:        order,
		Category:     "namespace",
		DefaultValue: namespaceInfo.Name,
		Enabled:      true,
		Validation: &FilterValidation{
			Pattern: "^[a-zA-Z0-9_-]+$",
		},
		Options: []FilterOption{
			{
				Value: namespaceInfo.Name,
				Label: namespaceInfo.Name,
			},
		},
	})
	order++

	// Add organization filter if organizations exist
	if len(namespaceInfo.Organizations) > 0 {
		orgOptions := make([]FilterOption, len(namespaceInfo.Organizations))
		for i, org := range namespaceInfo.Organizations {
			orgOptions[i] = FilterOption{
				Value: org,
				Label: ns.formatOptionLabel(org, FilterInputSelect),
			}
		}

		fields = append(fields, FilterFieldDefinition{
			ID:          "organization",
			BackendKey:  fmt.Sprintf("namespace/%s/organization", namespaceInfo.Name),
			DisplayName: "Organization",
			Description: fmt.Sprintf("Filter by organization within %s namespace", namespaceInfo.Name),
			InputType:   FilterInputSelect,
			Required:    false,
			Placeholder: "Select organization...",
			Order:       order,
			Category:    "namespace",
			Options:     orgOptions,
			Enabled:     true,
		})
		order++
	}

	// Add conversation filter if conversations exist
	if len(namespaceInfo.Conversations) > 0 {
		convOptions := make([]FilterOption, len(namespaceInfo.Conversations))
		for i, conv := range namespaceInfo.Conversations {
			convOptions[i] = FilterOption{
				Value: conv,
				Label: ns.formatConversationLabel(conv),
			}
		}

		fields = append(fields, FilterFieldDefinition{
			ID:          "conversation",
			BackendKey:  fmt.Sprintf("namespace/%s/conversation", namespaceInfo.Name),
			DisplayName: "Conversation",
			Description: fmt.Sprintf("Filter by conversation within %s namespace", namespaceInfo.Name),
			InputType:   FilterInputSelect,
			Required:    false,
			Placeholder: "Select conversation...",
			Order:       order,
			Category:    "namespace",
			Options:     convOptions,
			Enabled:     true,
		})
		order++
	}

	// Add data type filter if data types exist
	if len(namespaceInfo.DataTypes) > 0 {
		dataOptions := make([]FilterOption, len(namespaceInfo.DataTypes))
		for i, dataType := range namespaceInfo.DataTypes {
			dataOptions[i] = FilterOption{
				Value: dataType,
				Label: ns.formatDataTypeLabel(dataType),
			}
		}

		fields = append(fields, FilterFieldDefinition{
			ID:          "data_type",
			BackendKey:  fmt.Sprintf("namespace/%s/data", namespaceInfo.Name),
			DisplayName: "Data Type",
			Description: fmt.Sprintf("Filter by data type within %s namespace", namespaceInfo.Name),
			InputType:   FilterInputSelect,
			Required:    false,
			Placeholder: "Select data type...",
			Order:       order,
			Category:    "namespace",
			Options:     dataOptions,
			Enabled:     true,
		})
		order++
	}

	return fields
}

// buildRegularFilterFields builds fields from regular facets (metadata fields)
func (ns *NamespaceService) buildRegularFilterFields(response *fugusdk.SanitizedResponse, namespace string) []FilterFieldDefinition {
	var fields []FilterFieldDefinition

	// Parse facets from response similar to existing logic
	facets, err := ns.parseFacetsFromResponse(response)
	if err != nil {
		return fields
	}

	// Filter out namespace facets (we handle those separately)
	var metadataFacets []FacetInfo
	for _, facet := range facets {
		if !strings.HasPrefix(facet.Path, fmt.Sprintf("/namespace/%s/", namespace)) {
			metadataFacets = append(metadataFacets, facet)
		}
	}

	// Convert metadata facets to fields
	metadataFields := ns.convertFacetsToFields(metadataFacets)

	// Adjust order to come after namespace fields
	baseOrder := 100
	for i := range metadataFields {
		metadataFields[i].Order = baseOrder + i
		metadataFields[i].Category = "metadata"
	}

	return metadataFields
}

// createNamespaceFilterCategories creates categories including namespace-specific ones
func (ns *NamespaceService) createNamespaceFilterCategories() []FilterCategory {
	categories := []FilterCategory{
		{
			ID:          "namespace",
			Name:        "Namespace Filters",
			Description: "Filters based on namespace facets",
			Order:       0,
			Collapsible: false,
		},
		{
			ID:          "metadata",
			Name:        "Metadata Filters",
			Description: "Filters based on document metadata",
			Order:       1,
			Collapsible: true,
		},
	}

	// Add existing categories with adjusted order
	existingCategories := ns.createFilterCategories()
	for _, cat := range existingCategories {
		if cat.ID != "general" { // Skip general category
			cat.Order += 10
			categories = append(categories, cat)
		}
	}

	return categories
}

// validateNamespaceFieldValue performs namespace-specific validation
func (ns *NamespaceService) validateNamespaceFieldValue(ctx context.Context, field FilterFieldDefinition, value, namespace string) error {
	// Create fugu client for validation
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Validate namespace-specific fields
	switch field.ID {
	case "organization":
		orgs, err := client.GetNamespaceOrganizations(ctx, namespace)
		if err != nil {
			return fmt.Errorf("failed to get organizations: %w", err)
		}
		for _, org := range orgs {
			if org == value {
				return nil
			}
		}
		return fmt.Errorf("organization '%s' not found in namespace '%s'", value, namespace)

	case "conversation":
		convs, err := client.GetNamespaceConversations(ctx, namespace)
		if err != nil {
			return fmt.Errorf("failed to get conversations: %w", err)
		}
		for _, conv := range convs {
			if conv == value {
				return nil
			}
		}
		return fmt.Errorf("conversation '%s' not found in namespace '%s'", value, namespace)

	case "data_type":
		dataTypes, err := client.GetNamespaceDataTypes(ctx, namespace)
		if err != nil {
			return fmt.Errorf("failed to get data types: %w", err)
		}
		for _, dt := range dataTypes {
			if dt == value {
				return nil
			}
		}
		return fmt.Errorf("data type '%s' not found in namespace '%s'", value, namespace)
	}

	return nil
}

// convertMetadataFilter converts a metadata filter to facet format
func (ns *NamespaceService) convertMetadataFilter(fieldID, value, namespace string) string {
	// Convert field ID back to metadata path
	// This is a simplified conversion - you may need more sophisticated mapping
	metadataPath := fmt.Sprintf("/metadata/%s/%s", fieldID, value)
	return metadataPath
}

// formatConversationLabel formats conversation labels for display
func (ns *NamespaceService) formatConversationLabel(conversationID string) string {
	// Convert conversation IDs to human-readable labels
	// Example: "hearing-2024-03-15" -> "Hearing (Mar 15, 2024)"

	parts := strings.Split(conversationID, "-")
	if len(parts) >= 4 {
		convType := strings.Title(parts[0])
		year := parts[1]
		month := parts[2]
		day := parts[3]

		// Convert month number to name
		monthNames := map[string]string{
			"01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr",
			"05": "May", "06": "Jun", "07": "Jul", "08": "Aug",
			"09": "Sep", "10": "Oct", "11": "Nov", "12": "Dec",
		}

		if monthName, ok := monthNames[month]; ok {
			return fmt.Sprintf("%s (%s %s, %s)", convType, monthName, day, year)
		}
	}

	// Fallback to formatted version
	return ns.formatOptionLabel(conversationID, FilterInputSelect)
}

// formatDataTypeLabel formats data type labels for display
func (ns *NamespaceService) formatDataTypeLabel(dataType string) string {
	// Convert data types to human-readable labels
	// Example: "consumption_stats" -> "Consumption Statistics"

	commonTypes := map[string]string{
		"consumption_stats":    "Consumption Statistics",
		"renewable_report":     "Renewable Energy Report",
		"environmental_report": "Environmental Report",
		"rate_case":            "Rate Case Filing",
		"regulatory_filing":    "Regulatory Filing",
		"research_report":      "Research Report",
		"energy_statistics":    "Energy Statistics",
		"compliance_filing":    "Compliance Filing",
		"technical_report":     "Technical Report",
		"financial_report":     "Financial Report",
		"market_analysis":      "Market Analysis",
		"policy_document":      "Policy Document",
		"meeting_minutes":      "Meeting Minutes",
		"public_comment":       "Public Comment",
		"tariff_schedule":      "Tariff Schedule",
		"service_agreement":    "Service Agreement",
	}

	if label, ok := commonTypes[dataType]; ok {
		return label
	}

	// Fallback to general formatting
	return ns.formatOptionLabel(dataType, FilterInputSelect)
}

// GetAvailableNamespaces retrieves all available namespaces from the system
func (ns *NamespaceService) GetAvailableNamespaces(ctx context.Context) ([]string, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:get-available-namespaces")
	defer span.End()

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	namespaces, err := client.GetAvailableNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available namespaces: %w", err)
	}

	return namespaces, nil
}

// GetNamespaceStatistics returns statistics about namespace usage
func (ns *NamespaceService) GetNamespaceStatistics(ctx context.Context, namespace string) (map[string]interface{}, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:get-statistics")
	defer span.End()

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Get namespace info
	info, err := client.GetNamespaceInfo(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace info: %w", err)
	}

	// Build statistics
	stats := map[string]interface{}{
		"namespace":          namespace,
		"organization_count": len(info.Organizations),
		"conversation_count": len(info.Conversations),
		"data_type_count":    len(info.DataTypes),
		"organizations":      info.Organizations,
		"conversations":      info.Conversations,
		"data_types":         info.DataTypes,
		"last_updated":       time.Now().Format(time.RFC3339),
	}

	// Get document counts for each category
	if len(info.Organizations) > 0 {
		orgStats := make(map[string]int)
		for _, org := range info.Organizations {
			response, err := client.SearchByOrganization(ctx, "*", namespace, org, 0, 1000)
			if err == nil {
				orgStats[org] = len(response.Results)
			}
		}
		stats["organization_document_counts"] = orgStats
	}

	if len(info.Conversations) > 0 {
		convStats := make(map[string]int)
		for _, conv := range info.Conversations {
			response, err := client.SearchByConversation(ctx, "*", namespace, conv, 0, 1000)
			if err == nil {
				convStats[conv] = len(response.Results)
			}
		}
		stats["conversation_document_counts"] = convStats
	}

	if len(info.DataTypes) > 0 {
		dataStats := make(map[string]int)
		for _, dataType := range info.DataTypes {
			response, err := client.SearchByDataType(ctx, "*", namespace, dataType, 0, 1000)
			if err == nil {
				dataStats[dataType] = len(response.Results)
			}
		}
		stats["data_type_document_counts"] = dataStats
	}

	return stats, nil
}

// SuggestFiltersForNamespace analyzes a namespace and suggests useful filter combinations
func (ns *NamespaceService) SuggestFiltersForNamespace(ctx context.Context, namespace string) ([]map[string]string, error) {
	ctx, span := serviceTracer.Start(ctx, "namespace-service:suggest-filters")
	defer span.End()

	// Get namespace info
	client, err := fugusdk.NewClient(ctx, ns.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	info, err := client.GetNamespaceInfo(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace info: %w", err)
	}

	var suggestions []map[string]string

	// Suggest popular organizations
	if len(info.Organizations) > 0 {
		for _, org := range info.Organizations[:min(3, len(info.Organizations))] {
			suggestions = append(suggestions, map[string]string{
				"name":         fmt.Sprintf("Documents from %s", ns.formatOptionLabel(org, FilterInputSelect)),
				"organization": org,
				"description":  fmt.Sprintf("All documents from %s organization", org),
			})
		}
	}

	// Suggest data type combinations
	if len(info.DataTypes) > 0 {
		for _, dataType := range info.DataTypes[:min(3, len(info.DataTypes))] {
			suggestions = append(suggestions, map[string]string{
				"name":        fmt.Sprintf("%s Documents", ns.formatDataTypeLabel(dataType)),
				"data_type":   dataType,
				"description": fmt.Sprintf("All %s documents in %s", ns.formatDataTypeLabel(dataType), namespace),
			})
		}
	}

	// Suggest recent conversations
	if len(info.Conversations) > 0 {
		for _, conv := range info.Conversations[:min(2, len(info.Conversations))] {
			suggestions = append(suggestions, map[string]string{
				"name":         ns.formatConversationLabel(conv),
				"conversation": conv,
				"description":  fmt.Sprintf("Documents from %s", ns.formatConversationLabel(conv)),
			})
		}
	}

	return suggestions, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
