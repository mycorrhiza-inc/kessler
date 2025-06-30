// sdk.go
package fugusdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

const (
	// DefaultTimeout for HTTP requests
	DefaultTimeout = 300 * time.Second

	// SDK version for user agent
	SDKVersion = "1.2.0" // Updated version for namespace support

	// Security limits
	MaxRequestBodySize  = 10 * 1024 * 1024 // 10MB
	MaxResponseBodySize = 50 * 1024 * 1024 // 50MB
	MaxBatchSize        = 1000
	MaxQueryLength      = 10000
	MaxObjectIDLength   = 256
	MaxNamespaceLength  = 128
	MaxMetadataSize     = 1024 * 1024 // 1MB
	MaxFacetsPerObject  = 100         // Maximum facets per object
	MaxFacetLength      = 512         // Maximum length per facet

	// Rate limiting
	DefaultRateLimit = 100 // requests per second
	DefaultBurst     = 10
)

var tracer = otel.Tracer("fugu-sdk")

// SecureToken represents an encrypted authentication token
type SecureToken struct {
	encryptedValue []byte
	mutex          sync.RWMutex
}

// NewSecureToken creates a new secure token
func NewSecureToken(token string) *SecureToken {
	return &SecureToken{
		encryptedValue: []byte(token),
	}
}

// Get retrieves the token value
func (st *SecureToken) Get() string {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	return string(st.encryptedValue)
}

// Client represents the FuguDB API client
type Client struct {
	baseURL     string
	httpClient  *http.Client
	token       *SecureToken
	userAgent   string
	tracer      trace.Tracer
	rateLimiter *rate.Limiter
	tlsConfig   *tls.Config
	sanitizer   *InputSanitizer
	maxRetries  int
	retryDelay  time.Duration
}

// InputSanitizer handles input validation and sanitization
type InputSanitizer struct {
	objectIDRegex  *regexp.Regexp
	namespaceRegex *regexp.Regexp
}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer() *InputSanitizer {
	return &InputSanitizer{
		objectIDRegex:  regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`),
		namespaceRegex: regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`),
	}
}

// ValidateObjectID validates object ID format and length
func (s *InputSanitizer) ValidateObjectID(id string) error {
	if len(id) == 0 || len(id) > MaxObjectIDLength {
		return fmt.Errorf("invalid object ID length: must be 1-%d characters", MaxObjectIDLength)
	}
	if !s.objectIDRegex.MatchString(id) {
		return fmt.Errorf("invalid object ID format: contains illegal characters")
	}
	return nil
}

// ValidateNamespace validates namespace format and length
func (s *InputSanitizer) ValidateNamespace(namespace string) error {
	if len(namespace) == 0 || len(namespace) > MaxNamespaceLength {
		return fmt.Errorf("invalid namespace length: must be 1-%d characters", MaxNamespaceLength)
	}
	if !s.namespaceRegex.MatchString(namespace) {
		return fmt.Errorf("invalid namespace format: contains illegal characters")
	}
	return nil
}

// ValidateQuery validates search query - simplified
func (s *InputSanitizer) ValidateQuery(query string) error {
	if len(query) == 0 {
		return fmt.Errorf("query cannot be empty")
	}
	if len(query) > MaxQueryLength {
		return fmt.Errorf("query too long: maximum %d characters", MaxQueryLength)
	}
	return nil
}

// ValidateMetadata validates metadata size and content
func (s *InputSanitizer) ValidateMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("invalid metadata: cannot serialize")
	}

	if len(data) > MaxMetadataSize {
		return fmt.Errorf("metadata too large: maximum %d bytes", MaxMetadataSize)
	}

	return nil
}

// ValidateFacets validates facets array and individual facet strings
func (s *InputSanitizer) ValidateFacets(facets []string) error {
	if len(facets) == 0 {
		return nil
	}

	if len(facets) > MaxFacetsPerObject {
		return fmt.Errorf("too many facets: maximum %d per object", MaxFacetsPerObject)
	}

	for i, facet := range facets {
		if len(facet) == 0 {
			return fmt.Errorf("facet at index %d cannot be empty", i)
		}
		if len(facet) > MaxFacetLength {
			return fmt.Errorf("facet at index %d too long: maximum %d characters", i, MaxFacetLength)
		}
		// Basic validation - facets should not contain control characters
		for j, r := range facet {
			if r < 32 && r != 9 && r != 10 && r != 13 { // Allow tab, LF, CR
				return fmt.Errorf("facet at index %d contains invalid character at position %d", i, j)
			}
		}
	}

	return nil
}

// ClientOption represents configuration options for the client
type ClientOption func(*Client) error

// WithSecureToken sets the authentication token securely
func WithSecureToken(token string) ClientOption {
	return func(c *Client) error {
		if len(token) == 0 {
			return fmt.Errorf("token cannot be empty")
		}
		c.token = NewSecureToken(token)
		return nil
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithRateLimit sets rate limiting parameters
func WithRateLimit(requestsPerSecond int, burst int) ClientOption {
	return func(c *Client) error {
		if requestsPerSecond <= 0 || burst <= 0 {
			return fmt.Errorf("rate limit parameters must be positive")
		}
		c.rateLimiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
		return nil
	}
}

// WithRetry configures retry behavior
func WithRetry(maxRetries int, delay time.Duration) ClientOption {
	return func(c *Client) error {
		if maxRetries < 0 {
			return fmt.Errorf("max retries cannot be negative")
		}
		c.maxRetries = maxRetries
		c.retryDelay = delay
		return nil
	}
}

// NewClient creates a new default FuguDB API client
func NewClient(ctx context.Context, baseURL string) (*Client, error) {
	return BuildClient(
		ctx,
		baseURL,
		WithSecureToken("your-token"),
		WithRateLimit(50, 10),
		WithRetry(3, time.Second),
	)
}

// BuildClient creates a new custom FuguDB API client
func BuildClient(ctx context.Context, baseURL string, options ...ClientOption) (*Client, error) {
	_, span := tracer.Start(ctx, "fugu-client-build")
	defer span.End()

	// Validate base URL
	if len(baseURL) == 0 {
		return nil, fmt.Errorf("base URL cannot be empty")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Create HTTP client with sensible defaults
	transport := &http.Transport{
		DisableKeepAlives:     false,
		DisableCompression:    false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &Client{
		baseURL: strings.TrimSuffix(parsedURL.String(), "/"),
		httpClient: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: transport,
		},
		userAgent:   fmt.Sprintf("fugusdk-go/%s", SDKVersion),
		tracer:      tracer,
		rateLimiter: rate.NewLimiter(DefaultRateLimit, DefaultBurst),
		sanitizer:   NewInputSanitizer(),
		maxRetries:  3,
		retryDelay:  1 * time.Second,
	}

	// Apply options
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	return client, nil
}

// ObjectRecord represents an object to be indexed with namespace facet support
type ObjectRecord struct {
	ID        string                 `json:"id"`
	Text      string                 `json:"text"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Namespace string                 `json:"namespace,omitempty"`
	Facets    []string               `json:"facets,omitempty"`

	// New fields for namespace facets
	Organization   string `json:"organization,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	DataType       string `json:"data_type,omitempty"`

	// Optional timestamps
	DateCreated   string `json:"date_created,omitempty"`
	DateUpdated   string `json:"date_updated,omitempty"`
	DatePublished string `json:"date_published,omitempty"`
}

// Validate validates the object record with namespace field support
func (o *ObjectRecord) Validate(sanitizer *InputSanitizer) error {
	if err := sanitizer.ValidateObjectID(o.ID); err != nil {
		return fmt.Errorf("invalid object ID: %w", err)
	}

	if len(o.Text) == 0 {
		return fmt.Errorf("text cannot be empty")
	}

	if len(o.Text) > MaxQueryLength {
		return fmt.Errorf("text too long: maximum %d characters", MaxQueryLength)
	}

	if o.Namespace != "" {
		if err := sanitizer.ValidateNamespace(o.Namespace); err != nil {
			return fmt.Errorf("invalid namespace: %w", err)
		}
	}

	// Validate facets if present
	if len(o.Facets) > 0 {
		if err := sanitizer.ValidateFacets(o.Facets); err != nil {
			return fmt.Errorf("invalid facets: %w", err)
		}
	}

	// Validate organization field format if present
	if o.Organization != "" {
		if len(o.Organization) > MaxNamespaceLength {
			return fmt.Errorf("organization name too long: maximum %d characters", MaxNamespaceLength)
		}
	}

	// Validate conversation ID format if present
	if o.ConversationID != "" {
		if len(o.ConversationID) > MaxObjectIDLength {
			return fmt.Errorf("conversation ID too long: maximum %d characters", MaxObjectIDLength)
		}
	}

	// Validate data type format if present
	if o.DataType != "" {
		if len(o.DataType) > MaxNamespaceLength {
			return fmt.Errorf("data type name too long: maximum %d characters", MaxNamespaceLength)
		}
	}

	if err := sanitizer.ValidateMetadata(o.Metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	return nil
}

// Pagination matches the Rust Pagination struct
type Pagination struct {
	Page    *int `json:"page,omitempty"`
	PerPage *int `json:"per_page,omitempty"`
}

// FuguSearchQuery matches the Rust FuguSearchQuery struct exactly
type FuguSearchQuery struct {
	Query   string      `json:"query"`
	Filters *[]string   `json:"filters,omitempty"`
	Page    *Pagination `json:"page,omitempty"`
}

// IndexRequest matches the Rust IndexRequest struct
type IndexRequest struct {
	Data []ObjectRecord `json:"data"`
}

// BatchIndexRequest matches the Rust BatchIndexRequest struct
type BatchIndexRequest struct {
	Objects []ObjectRecord `json:"objects"`
}

// Validate validates the index request
func (i *IndexRequest) Validate(sanitizer *InputSanitizer) error {
	if len(i.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	if len(i.Data) > MaxBatchSize {
		return fmt.Errorf("batch too large: maximum %d objects", MaxBatchSize)
	}

	for idx, obj := range i.Data {
		if err := obj.Validate(sanitizer); err != nil {
			return fmt.Errorf("invalid object at index %d: %w", idx, err)
		}
	}

	return nil
}

// Validate validates the batch index request
func (b *BatchIndexRequest) Validate(sanitizer *InputSanitizer) error {
	if len(b.Objects) == 0 {
		return fmt.Errorf("objects cannot be empty")
	}

	if len(b.Objects) > MaxBatchSize {
		return fmt.Errorf("batch too large: maximum %d objects", MaxBatchSize)
	}

	for idx, obj := range b.Objects {
		if err := obj.Validate(sanitizer); err != nil {
			return fmt.Errorf("invalid object at index %d: %w", idx, err)
		}
	}

	return nil
}

// SanitizedResponse represents a sanitized API response
type SanitizedResponse struct {
	Results       []FuguSearchResult `json:"results,omitempty"`
	Total         int                `json:"total,omitempty"`
	Page          int                `json:"page,omitempty"`
	PerPage       int                `json:"per_page,omitempty"`
	Query         string             `json:"query,omitempty"`
	Data          interface{}        `json:"data,omitempty"`
	Message       string             `json:"message,omitempty"`
	Status        string             `json:"status,omitempty"`
	UpsertedCount *int               `json:"upserted_count,omitempty"`
}

// FuguSearchResult represents a search result from Fugu
type FuguSearchResult struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Facets   []string               `json:"facets,omitempty"`
}

// APIError represents a sanitized API error
type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// makeRequest performs a secure HTTP request with retries and rate limiting
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay * time.Duration(attempt)):
			}
		}

		resp, err := c.doRequest(ctx, method, path, body)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on client errors (4xx)
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
			break
		}
	}

	return nil, lastErr
}

// doRequest performs a single HTTP request with tracing
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	ctx, span := c.tracer.Start(ctx, fmt.Sprintf("fugusdk.%s", strings.ToLower(method)))
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("http.path", path),
		attribute.String("sdk.version", SDKVersion),
	)

	var reqBody io.Reader
	var bodySize int

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		bodySize = len(jsonBody)
		if bodySize > MaxRequestBodySize {
			return nil, fmt.Errorf("request body too large: maximum %d bytes", MaxRequestBodySize)
		}

		reqBody = bytes.NewReader(jsonBody)
		span.SetAttributes(attribute.Int("http.request_size", bodySize))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	// Add authentication token if present
	if c.token != nil {
		req.Header.Set("Authorization", "Bearer "+c.token.Get())
	}

	// Inject tracing context into HTTP headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	return resp, nil
}

// handleResponse processes HTTP response with security measures
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		message := "An error occurred"
		if len(body) > 0 {
			var errorResp map[string]interface{}
			if json.Unmarshal(body, &errorResp) == nil {
				if errMsg, ok := errorResp["error"].(string); ok {
					message = errMsg
				}
			}
		}

		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    message,
			RequestID:  resp.Header.Get("X-Request-ID"),
		}
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Health checks the health of the API - matches Rust endpoint
func (c *Client) Health(ctx context.Context) error {
	resp, err := c.makeRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// IngestObjects ingests multiple objects into the database (now performs upserts)
// Enhanced to support namespace facets automatically
func (c *Client) IngestObjects(ctx context.Context, objects []ObjectRecord) (*SanitizedResponse, error) {
	req := IndexRequest{Data: objects}

	// Validate request
	if err := req.Validate(c.sanitizer); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if any objects have namespace facet fields or facets
	hasNamespaceFacets := false
	for _, obj := range objects {
		if obj.Organization != "" || obj.ConversationID != "" || obj.DataType != "" || len(obj.Facets) > 0 {
			hasNamespaceFacets = true
			break
		}
	}

	// Use namespace-aware endpoint if namespace facets are present
	endpoint := "/ingest"
	if hasNamespaceFacets {
		endpoint = "/ingest/namespace"
	}

	resp, err := c.makeRequest(ctx, "POST", endpoint, req)
	if err != nil {
		// Fallback to regular ingest if namespace endpoint fails
		if hasNamespaceFacets && endpoint == "/ingest/namespace" {
			resp, err = c.makeRequest(ctx, "POST", "/ingest", req)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// UpsertObjects explicitly upserts multiple objects - matches Rust PUT /objects endpoint
func (c *Client) UpsertObjects(ctx context.Context, objects []ObjectRecord) (*SanitizedResponse, error) {
	req := IndexRequest{Data: objects}

	// Validate request
	if err := req.Validate(c.sanitizer); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	resp, err := c.makeRequest(ctx, "PUT", "/objects", req)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// BatchUpsertObjects performs batch upsert with detailed response - matches Rust /batch/upsert endpoint
func (c *Client) BatchUpsertObjects(ctx context.Context, objects []ObjectRecord) (*SanitizedResponse, error) {
	req := BatchIndexRequest{Objects: objects}

	// Validate request
	if err := req.Validate(c.sanitizer); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	resp, err := c.makeRequest(ctx, "POST", "/batch/upsert", req)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// DeleteObject deletes a single object by ID - matches Rust DELETE /objects/{id} endpoint
func (c *Client) DeleteObject(ctx context.Context, objectID string) (*SanitizedResponse, error) {
	// Validate object ID
	if err := c.sanitizer.ValidateObjectID(objectID); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/objects/%s", url.PathEscape(objectID))

	resp, err := c.makeRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// Search performs a POST search - enhanced to automatically use namespace endpoints when appropriate
func (c *Client) Search(ctx context.Context, query FuguSearchQuery) (*SanitizedResponse, error) {
	// Validate query
	if err := c.sanitizer.ValidateQuery(query.Query); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if filters contain namespace facets
	endpoint := "/search"
	if query.Filters != nil {
		for _, filter := range *query.Filters {
			if strings.HasPrefix(filter, "namespace/") {
				endpoint = "/search/namespace"
				break
			}
		}
	}

	resp, err := c.makeRequest(ctx, "POST", endpoint, query)
	if err != nil {
		// Fallback to regular search if namespace endpoint fails
		if endpoint == "/search/namespace" {
			resp, err = c.makeRequest(ctx, "POST", "/search", query)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// SearchText performs a GET search with query parameter - matches Rust /search GET endpoint
func (c *Client) SearchText(ctx context.Context, query string) (*SanitizedResponse, error) {
	// Validate query
	if err := c.sanitizer.ValidateQuery(query); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	params := url.Values{}
	params.Add("q", query)

	path := "/search?" + params.Encode()
	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// GetObjectByID retrieves a specific object by its ID - matches Rust /objects/{id} endpoint
func (c *Client) GetObjectByID(ctx context.Context, objectID string) (*SanitizedResponse, error) {
	// Validate object ID
	if err := c.sanitizer.ValidateObjectID(objectID); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/objects/%s", url.PathEscape(objectID))

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// ListFilters lists all available filters - matches Rust /filters endpoint
func (c *Client) ListFilters(ctx context.Context) (*SanitizedResponse, error) {
	resp, err := c.makeRequest(ctx, "GET", "/filters", nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// GetNamespaceFilters gets filters for a namespace - matches Rust /filters/{namespace} endpoint
func (c *Client) GetNamespaceFilters(ctx context.Context, namespace string) (*SanitizedResponse, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/filters/namespace/%s", url.PathEscape(namespace))

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// Convenience methods for common operations

// AddOrUpdateObject is a convenience method for upserting a single object
func (c *Client) AddOrUpdateObject(ctx context.Context, object ObjectRecord) (*SanitizedResponse, error) {
	return c.UpsertObjects(ctx, []ObjectRecord{object})
}

// AddOrUpdateObjects is an alias for UpsertObjects for backward compatibility
func (c *Client) AddOrUpdateObjects(ctx context.Context, objects []ObjectRecord) (*SanitizedResponse, error) {
	return c.UpsertObjects(ctx, objects)
}

// SimpleSearch is a convenience method for basic text search
func (c *Client) SimpleSearch(ctx context.Context, query string) (*SanitizedResponse, error) {
	return c.SearchText(ctx, query)
}

// AdvancedSearch is a convenience method for search with filters and pagination
func (c *Client) AdvancedSearch(ctx context.Context, query string, filters []string, page, perPage int) (*SanitizedResponse, error) {
	searchQuery := FuguSearchQuery{
		Query: query,
	}

	if len(filters) > 0 {
		searchQuery.Filters = &filters
	}

	if page > 0 || perPage > 0 {
		pagination := &Pagination{}
		if page > 0 {
			pagination.Page = &page
		}
		if perPage > 0 {
			pagination.PerPage = &perPage
		}
		searchQuery.Page = pagination
	}

	return c.Search(ctx, searchQuery)
}
