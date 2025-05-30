package fugusdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
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
	DefaultTimeout = 30 * time.Second

	// SDK version for user agent
	SDKVersion = "1.0.0"

	// Security limits
	MaxRequestBodySize  = 10 * 1024 * 1024 // 10MB
	MaxResponseBodySize = 50 * 1024 * 1024 // 50MB
	MaxBatchSize        = 1000
	MaxQueryLength      = 10000
	MaxObjectIDLength   = 256
	MaxNamespaceLength  = 128
	MaxMetadataSize     = 1024 * 1024 // 1MB

	// Rate limiting
	DefaultRateLimit = 100 // requests per second
	DefaultBurst     = 10
)

// SecureToken represents an encrypted authentication token
type SecureToken struct {
	encryptedValue []byte
	mutex          sync.RWMutex
}

// NewSecureToken creates a new secure token (simplified - in production use proper encryption)
func NewSecureToken(token string) *SecureToken {
	// In production, implement proper encryption with AES-GCM or similar
	return &SecureToken{
		encryptedValue: []byte(token), // Simplified for example
	}
}

// Get retrieves the token value (in production, decrypt here)
func (st *SecureToken) Get() string {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	return string(st.encryptedValue) // Simplified for example
}

// Client represents the hardened FuguDB API client
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
	queryRegex     *regexp.Regexp
}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer() *InputSanitizer {
	return &InputSanitizer{
		objectIDRegex:  regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{1,256}$`),
		namespaceRegex: regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,128}$`),
		queryRegex:     regexp.MustCompile(`^[\p{L}\p{N}\p{P}\p{S}\s]{1,10000}$`),
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

// ValidateQuery validates search query
func (s *InputSanitizer) ValidateQuery(query string) error {
	if len(query) == 0 || len(query) > MaxQueryLength {
		return fmt.Errorf("invalid query length: must be 1-%d characters", MaxQueryLength)
	}
	if !s.queryRegex.MatchString(query) {
		return fmt.Errorf("invalid query format: contains illegal characters")
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

// WithTLSConfig sets TLS configuration for secure connections
func WithTLSConfig(tlsConfig *tls.Config) ClientOption {
	return func(c *Client) error {
		if tlsConfig == nil {
			return fmt.Errorf("TLS config cannot be nil")
		}
		c.tlsConfig = tlsConfig
		return nil
	}
}

// WithCertificatePinning enables certificate pinning with provided certificates
func WithCertificatePinning(certs []string) ClientOption {
	return func(c *Client) error {
		certPool := x509.NewCertPool()
		for _, cert := range certs {
			if !certPool.AppendCertsFromPEM([]byte(cert)) {
				return fmt.Errorf("failed to parse certificate")
			}
		}

		c.tlsConfig = &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}
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

// WithUserAgent sets a custom user agent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		if len(userAgent) == 0 {
			return fmt.Errorf("user agent cannot be empty")
		}
		c.userAgent = userAgent
		return nil
	}
}

// WithTracer sets a custom OpenTelemetry tracer
func WithTracer(tracer trace.Tracer) ClientOption {
	return func(c *Client) error {
		if tracer == nil {
			return fmt.Errorf("tracer cannot be nil")
		}
		c.tracer = tracer
		return nil
	}
}

// NewClient creates a new hardened FuguDB API client
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	// Validate base URL
	if len(baseURL) == 0 {
		return nil, fmt.Errorf("base URL cannot be empty")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("only HTTPS URLs are allowed")
	}

	// Create secure HTTP client with sensible defaults
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			},
		},
		DisableKeepAlives:     false,
		DisableCompression:    false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: transport,
		},
		userAgent:   fmt.Sprintf("fugusdk-go/%s", SDKVersion),
		tracer:      otel.Tracer("fugusdk"),
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

	// Update transport with custom TLS config if provided
	if client.tlsConfig != nil {
		transport.TLSClientConfig = client.tlsConfig
	}

	return client, nil
}

// ObjectRecord represents a sanitized object to be indexed
type ObjectRecord struct {
	ID        string                 `json:"id"`
	Text      string                 `json:"text"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Namespace string                 `json:"namespace,omitempty"`
}

// Validate validates the object record
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

	if err := sanitizer.ValidateMetadata(o.Metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	return nil
}

// SearchQuery represents a sanitized search request
type SearchQuery struct {
	Query   string   `json:"query"`
	Filters []string `json:"filters,omitempty"`
}

// Validate validates the search query
func (s *SearchQuery) Validate(sanitizer *InputSanitizer) error {
	if err := sanitizer.ValidateQuery(s.Query); err != nil {
		return fmt.Errorf("invalid query: %w", err)
	}

	for _, filter := range s.Filters {
		if len(filter) > MaxQueryLength {
			return fmt.Errorf("filter too long: maximum %d characters", MaxQueryLength)
		}
	}

	return nil
}

// IndexRequest represents a request to index objects
type IndexRequest struct {
	Data []ObjectRecord `json:"data"`
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

// SanitizedResponse represents a sanitized API response
type SanitizedResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
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

// sanitizeErrorMessage removes sensitive information from error messages
func sanitizeErrorMessage(message string) string {
	// Remove potential sensitive patterns
	sensitive := []string{
		"password", "token", "key", "secret", "auth",
		"database", "internal", "stack trace", "sql",
	}

	lower := strings.ToLower(message)
	for _, word := range sensitive {
		if strings.Contains(lower, word) {
			return "An internal error occurred. Please contact support."
		}
	}

	// Limit message length
	if len(message) > 200 {
		return message[:200] + "..."
	}

	return message
}

// limitedReader wraps an io.Reader to limit the number of bytes read
type limitedReader struct {
	reader io.Reader
	limit  int64
	read   int64
}

func (lr *limitedReader) Read(p []byte) (int, error) {
	if lr.read >= lr.limit {
		return 0, fmt.Errorf("response body too large: maximum %d bytes", lr.limit)
	}

	remaining := lr.limit - lr.read
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err := lr.reader.Read(p)
	lr.read += int64(n)
	return n, err
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

// doRequest performs a single HTTP request with tracing and security measures
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Start tracing span (with minimal sensitive data)
	ctx, span := c.tracer.Start(ctx, fmt.Sprintf("fugusdk.%s", strings.ToLower(method)))
	defer span.End()

	// Add safe attributes to span
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

	// Set secure headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("X-Content-Type-Options", "nosniff")
	req.Header.Set("X-Frame-Options", "DENY")
	req.Header.Set("X-XSS-Protection", "1; mode=block")

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

	// Add safe response attributes to span
	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
	)

	return resp, nil
}

// handleResponse processes HTTP response with security measures
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Limit response body size
	limitedBody := &limitedReader{
		reader: resp.Body,
		limit:  MaxResponseBodySize,
	}

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		// Sanitize error message to prevent information disclosure
		message := "An error occurred"
		if len(body) > 0 {
			var errorResp map[string]interface{}
			if json.Unmarshal(body, &errorResp) == nil {
				if errMsg, ok := errorResp["error"].(string); ok {
					message = sanitizeErrorMessage(errMsg)
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

// Health checks the health of the API
func (c *Client) Health(ctx context.Context) error {
	resp, err := c.makeRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// IngestObjects ingests multiple objects into the database with validation
func (c *Client) IngestObjects(ctx context.Context, objects []ObjectRecord) error {
	req := IndexRequest{Data: objects}

	// Validate request
	if err := req.Validate(c.sanitizer); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	resp, err := c.makeRequest(ctx, "POST", "/ingest", req)
	if err != nil {
		return err
	}

	var result SanitizedResponse
	return c.handleResponse(resp, &result)
}

// Search performs a validated search across all namespaces
func (c *Client) Search(ctx context.Context, query SearchQuery) (*SanitizedResponse, error) {
	// Validate query
	if err := query.Validate(c.sanitizer); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	resp, err := c.makeRequest(ctx, "POST", "/search", query)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// SearchText performs a simple validated text search
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

// AdvancedQuery performs a validated advanced query
func (c *Client) AdvancedQuery(ctx context.Context, queryData interface{}) (*SanitizedResponse, error) {
	// Basic validation of query data size
	if jsonData, err := json.Marshal(queryData); err != nil {
		return nil, fmt.Errorf("invalid query data: %w", err)
	} else if len(jsonData) > MaxRequestBodySize {
		return nil, fmt.Errorf("query data too large: maximum %d bytes", MaxRequestBodySize)
	}

	resp, err := c.makeRequest(ctx, "POST", "/query/advanced", queryData)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// GetObjectByID retrieves a specific object by its validated ID
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

// ListFilters lists all available filters
func (c *Client) ListFilters(ctx context.Context) (*SanitizedResponse, error) {
	resp, err := c.makeRequest(ctx, "GET", "/filters", nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// GetNamespaceFilters gets filters for a validated namespace
func (c *Client) GetNamespaceFilters(ctx context.Context, namespace string) (*SanitizedResponse, error) {
	// Validate namespace
	if err := c.sanitizer.ValidateNamespace(namespace); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	path := fmt.Sprintf("/filters/%s", url.PathEscape(namespace))

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SanitizedResponse
	err = c.handleResponse(resp, &result)
	return &result, err
}

// BatchClient provides secure batch operations
type BatchClient struct {
	*Client
	objects []ObjectRecord
	mutex   sync.Mutex
}

// NewBatch creates a new secure batch client
func (c *Client) NewBatch() *BatchClient {
	return &BatchClient{
		Client:  c,
		objects: make([]ObjectRecord, 0),
	}
}

// AddObject adds a validated object to the batch
func (b *BatchClient) AddObject(object ObjectRecord) error {
	// Validate object
	if err := object.Validate(b.sanitizer); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.objects) >= MaxBatchSize {
		return fmt.Errorf("batch size limit exceeded: maximum %d objects", MaxBatchSize)
	}

	b.objects = append(b.objects, object)
	return nil
}

// Execute executes the batch operation with validation
func (b *BatchClient) Execute(ctx context.Context) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.objects) == 0 {
		return nil
	}

	err := b.IngestObjects(ctx, b.objects)
	if err != nil {
		return err
	}

	// Clear the batch after successful execution
	b.objects = b.objects[:0]
	return nil
}

// Size returns the current batch size safely
func (b *BatchClient) Size() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return len(b.objects)
}

// Clear clears the batch safely
func (b *BatchClient) Clear() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.objects = b.objects[:0]
}
