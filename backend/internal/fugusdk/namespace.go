// namespace.go
package fugusdk

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// NamespaceFacetInfo represents namespace facet information
type NamespaceFacetInfo struct {
	Path      string `json:"path"`
	Count     int64  `json:"count"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"` // "organization", "conversation", "data"
	Value     string `json:"value,omitempty"`
}

// NamespaceInfo represents namespace information
type NamespaceInfo struct {
	Name          string   `json:"name"`
	Organizations []string `json:"organizations,omitempty"`
	Conversations []string `json:"conversations,omitempty"`
	DataTypes     []string `json:"data_types,omitempty"`
}

// GenerateNamespaceFacets creates namespace facets based on the object's fields Format: /namespace/{namespace}/{type} (NOT including specific values)
func (o *ObjectRecord) GenerateNamespaceFacets() []string {
	var facets []string

	if o.Namespace != "" {
		// Add base namespace facet
		facets = append(facets, fmt.Sprintf("namespace/%s", o.Namespace))

		// Add organization type facet if present (but not the specific organization)
		if o.Organization != "" {
			facets = append(facets, fmt.Sprintf("namespace/%s/organization", o.Namespace))
		}

		// Add conversation type facet if present (but not the specific conversation)
		if o.ConversationID != "" {
			facets = append(facets, fmt.Sprintf("namespace/%s/conversation", o.Namespace))
		}

		// Add data type facet if present (but not the specific data type)
		if o.DataType != "" {
			facets = append(facets, fmt.Sprintf("namespace/%s/data", o.Namespace))
		}
	}

	return facets
}

// Enhanced Client methods for namespace facet support

// IngestObjectsWithNamespaceFacets ingests objects with enhanced namespace facet support
func (c *Client) IngestObjectsWithNamespaceFacets(ctx context.Context, objects []ObjectRecord) (*SanitizedResponse, error) {
	req := IndexRequest{Data: objects}

	// Validate request
	if err := req.Validate(c.sanitizer); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Log the namespace facets that will be generated (optional debug)
	for i, obj := range objects {
		facets := obj.GenerateNamespaceFacets()
		if len(facets) > 0 {
			fmt.Printf("Object %d will generate facets: %v\n", i, facets)
		}
	}

	resp, err := c.makeRequest(ctx, "POST", "/ingest/namespace", req)
	if err != nil {
		// Fallback to regular ingest if namespace endpoint not available
		resp, err = c.makeRequest(ctx, "POST", "/ingest", req)
		if err != nil {
			return nil, err
		}
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// SearchWithNamespaceFacets performs search with namespace facet filtering
func (c *Client) SearchWithNamespaceFacets(ctx context.Context, query string, namespaceFacets []string, page, perPage int) (*SanitizedResponse, error) {
	// Validate query
	if err := c.sanitizer.ValidateQuery(query); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	searchQuery := FuguSearchQuery{
		Query: query,
	}

	if len(namespaceFacets) > 0 {
		searchQuery.Filters = &namespaceFacets
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

	// Try namespace-specific search endpoint first
	resp, err := c.makeRequest(ctx, "POST", "/search/namespace", searchQuery)
	if err != nil {
		// Fallback to regular search if namespace endpoint not available
		resp, err = c.makeRequest(ctx, "POST", "/search", searchQuery)
		if err != nil {
			return nil, err
		}
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// GetAvailableNamespaces retrieves all available namespaces
func (c *Client) GetAvailableNamespaces(ctx context.Context) ([]string, error) {
	resp, err := c.makeRequest(ctx, "GET", "/namespaces", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status     string   `json:"status"`
		Namespaces []string `json:"namespaces"`
	}

	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Namespaces, nil
}

// GetNamespaceFacets retrieves facets for a specific namespace
func (c *Client) GetNamespaceFacets(ctx context.Context, namespace string) ([]NamespaceFacetInfo, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/namespaces/%s/facets", url.PathEscape(namespace))
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status    string `json:"status"`
		Namespace string `json:"namespace"`
		Facets    []struct {
			Path  string `json:"path"`
			Count int64  `json:"count"`
		} `json:"facets"`
	}

	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	// Convert to NamespaceFacetInfo
	var facetInfos []NamespaceFacetInfo
	for _, facet := range result.Facets {
		info := NamespaceFacetInfo{
			Path:      facet.Path,
			Count:     facet.Count,
			Namespace: namespace,
		}

		// Parse facet type and value from path
		if parsed := parseNamespaceFacetPath(facet.Path, namespace); parsed != nil {
			info.Type = parsed.Type
			info.Value = parsed.Value
		}

		facetInfos = append(facetInfos, info)
	}

	return facetInfos, nil
}

// GetNamespaceOrganizations retrieves organizations for a namespace
func (c *Client) GetNamespaceOrganizations(ctx context.Context, namespace string) ([]string, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/namespaces/%s/organizations", url.PathEscape(namespace))
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status        string   `json:"status"`
		Namespace     string   `json:"namespace"`
		Organizations []string `json:"organizations"`
	}

	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Organizations, nil
}

// GetNamespaceConversations retrieves conversations for a namespace
func (c *Client) GetNamespaceConversations(ctx context.Context, namespace string) ([]string, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/namespaces/%s/conversations", url.PathEscape(namespace))
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status        string   `json:"status"`
		Namespace     string   `json:"namespace"`
		Conversations []string `json:"conversations"`
	}

	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Conversations, nil
}

// GetNamespaceDataTypes retrieves data types for a namespace
func (c *Client) GetNamespaceDataTypes(ctx context.Context, namespace string) ([]string, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/namespaces/%s/data", url.PathEscape(namespace))
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status    string   `json:"status"`
		Namespace string   `json:"namespace"`
		DataTypes []string `json:"data_types"`
	}

	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.DataTypes, nil
}

// GetNamespaceInfo retrieves comprehensive information about a namespace
func (c *Client) GetNamespaceInfo(ctx context.Context, namespace string) (*NamespaceInfo, error) {
	// Get organizations, conversations, and data types in parallel
	orgChan := make(chan []string, 1)
	convChan := make(chan []string, 1)
	dataChan := make(chan []string, 1)
	errChan := make(chan error, 3)

	// Get organizations
	go func() {
		orgs, err := c.GetNamespaceOrganizations(ctx, namespace)
		if err != nil {
			errChan <- err
			return
		}
		orgChan <- orgs
	}()

	// Get conversations
	go func() {
		convs, err := c.GetNamespaceConversations(ctx, namespace)
		if err != nil {
			errChan <- err
			return
		}
		convChan <- convs
	}()

	// Get data types
	go func() {
		types, err := c.GetNamespaceDataTypes(ctx, namespace)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- types
	}()

	// Collect results
	info := &NamespaceInfo{Name: namespace}
	for i := 0; i < 3; i++ {
		select {
		case orgs := <-orgChan:
			info.Organizations = orgs
		case convs := <-convChan:
			info.Conversations = convs
		case types := <-dataChan:
			info.DataTypes = types
		case err := <-errChan:
			return nil, err
		}
	}

	return info, nil
}

// Enhanced FilterBuilder with namespace facet support
type NamespaceFacetBuilder struct {
	client    *Client
	namespace string
	filters   []string
}

// NewNamespaceFacetBuilder creates a new namespace facet builder
func (c *Client) NewNamespaceFacetBuilder(namespace string) *NamespaceFacetBuilder {
	return &NamespaceFacetBuilder{
		client:    c,
		namespace: namespace,
		filters:   make([]string, 0),
	}
}

// AddOrganizationFilter adds an organization filter
func (nfb *NamespaceFacetBuilder) AddOrganizationFilter(organization string) *NamespaceFacetBuilder {
	facet := fmt.Sprintf("namespace/%s/organization/%s", nfb.namespace, organization)
	nfb.filters = append(nfb.filters, facet)
	return nfb
}

// AddConversationFilter adds a conversation filter
func (nfb *NamespaceFacetBuilder) AddConversationFilter(conversationID string) *NamespaceFacetBuilder {
	facet := fmt.Sprintf("namespace/%s/conversation/%s", nfb.namespace, conversationID)
	nfb.filters = append(nfb.filters, facet)
	return nfb
}

// AddDataTypeFilter adds a data type filter
func (nfb *NamespaceFacetBuilder) AddDataTypeFilter(dataType string) *NamespaceFacetBuilder {
	facet := fmt.Sprintf("namespace/%s/data/%s", nfb.namespace, dataType)
	nfb.filters = append(nfb.filters, facet)
	return nfb
}

// AddCustomFacet adds a custom facet filter
func (nfb *NamespaceFacetBuilder) AddCustomFacet(facet string) *NamespaceFacetBuilder {
	nfb.filters = append(nfb.filters, facet)
	return nfb
}

// Search performs search with the built filters
func (nfb *NamespaceFacetBuilder) Search(ctx context.Context, query string, page, perPage int) (*SanitizedResponse, error) {
	return nfb.client.SearchWithNamespaceFacets(ctx, query, nfb.filters, page, perPage)
}

// GetFilters returns the current filter list
func (nfb *NamespaceFacetBuilder) GetFilters() []string {
	return nfb.filters
}

// Helper function to parse namespace facet paths
func parseNamespaceFacetPath(path, namespace string) *struct {
	Type  string
	Value string
} {
	// Expected formats:
	// /namespace/NYPUC/organization/someorg
	// /namespace/NYPUC/conversation/conv123
	// /namespace/NYPUC/data/sometype

	prefix := fmt.Sprintf("/namespace/%s/", namespace)
	if !strings.HasPrefix(path, prefix) {
		return nil
	}

	remaining := strings.TrimPrefix(path, prefix)
	parts := strings.Split(remaining, "/")

	if len(parts) >= 2 {
		return &struct {
			Type  string
			Value string
		}{
			Type:  parts[0], // "organization", "conversation", "data"
			Value: parts[1], // the actual value
		}
	}

	if len(parts) == 1 {
		return &struct {
			Type  string
			Value string
		}{
			Type:  parts[0],
			Value: "",
		}
	}

	return nil
}

// Convenience methods for common namespace operations

// CreateOrganizationObject creates an object record for an organization
func CreateOrganizationObject(id, text, namespace, organization string, metadata map[string]interface{}) ObjectRecord {
	return ObjectRecord{
		ID:           id,
		Text:         text,
		Namespace:    namespace,
		Organization: organization,
		Metadata:     metadata,
	}
}

// CreateConversationObject creates an object record for a conversation
func CreateConversationObject(id, text, namespace, conversationID string, metadata map[string]interface{}) ObjectRecord {
	return ObjectRecord{
		ID:             id,
		Text:           text,
		Namespace:      namespace,
		ConversationID: conversationID,
		Metadata:       metadata,
	}
}

// CreateDataObject creates an object record for data
func CreateDataObject(id, text, namespace, dataType string, metadata map[string]interface{}) ObjectRecord {
	return ObjectRecord{
		ID:        id,
		Text:      text,
		Namespace: namespace,
		DataType:  dataType,
		Metadata:  metadata,
	}
}

// SearchByOrganization searches for documents by organization within a namespace
func (c *Client) SearchByOrganization(ctx context.Context, query, namespace, organization string, page, perPage int) (*SanitizedResponse, error) {
	builder := c.NewNamespaceFacetBuilder(namespace)
	builder.AddOrganizationFilter(organization)
	return builder.Search(ctx, query, page, perPage)
}

// SearchByConversation searches for documents by conversation within a namespace
func (c *Client) SearchByConversation(ctx context.Context, query, namespace, conversationID string, page, perPage int) (*SanitizedResponse, error) {
	builder := c.NewNamespaceFacetBuilder(namespace)
	builder.AddConversationFilter(conversationID)
	return builder.Search(ctx, query, page, perPage)
}

// SearchByDataType searches for documents by data type within a namespace
func (c *Client) SearchByDataType(ctx context.Context, query, namespace, dataType string, page, perPage int) (*SanitizedResponse, error) {
	builder := c.NewNamespaceFacetBuilder(namespace)
	builder.AddDataTypeFilter(dataType)
	return builder.Search(ctx, query, page, perPage)
}
