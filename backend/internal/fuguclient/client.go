// Package fuguclient provides a Go client for interacting with the FuguDB API
package fuguclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client represents a FuguDB API client
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	debug      bool
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client for the FuguDB client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithDebug enables debug logging for the client
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// NewClient creates a new FuguDB API client
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	client := &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client, nil
}

// ObjectRecord represents a document to be indexed in FuguDB
type ObjectRecord struct {
	ID       string                 `json:"id"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IndexResponse represents the response from the ingest endpoint
type IndexResponse struct {
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
	ObjectID    string                 `json:"object_id"`
	Operation   string                 `json:"operation"`
	UniqueTerms int                    `json:"unique_terms"`
	TextLength  int                    `json:"text_length"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// BatchIndexRequest represents a batch of objects to index
type BatchIndexRequest struct {
	Objects []ObjectRecord `json:"objects"`
}

// BatchIndexResponse represents the response from the batch-ingest endpoint
type BatchIndexResponse struct {
	Status         string        `json:"status"`
	Message        string        `json:"message"`
	TotalObjects   int           `json:"total_objects"`
	Processed      int           `json:"processed"`
	Errors         int           `json:"errors"`
	ObjectsSummary []interface{} `json:"objects"`
}

// FileIngestionRequest represents a request to ingest a file
type FileIngestionRequest struct {
	FilePath  string                 `json:"file_path"`
	ID        *string                `json:"id,omitempty"`
	Namespace *string                `json:"namespace,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// SearchQuery represents a search query
type SearchQuery struct {
	Query     string  `json:"query"`
	Namespace *string `json:"namespace,omitempty"`
}

// AdvancedSearchQuery represents an advanced search query
type AdvancedSearchQuery struct {
	Query     string                 `json:"query"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Namespace *string                `json:"namespace,omitempty"`
	Page      *int                   `json:"page,omitempty"`
	PerPage   *int                   `json:"per_page,omitempty"`
}

// ObjectResponse represents a response containing an object
type ObjectResponse struct {
	ID         string                 `json:"id"`
	Metadata   map[string]interface{} `json:"metadata"`
	Text       string                 `json:"text"`
	TextLength int                    `json:"text_length"`
	Note       string                 `json:"note,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error        string                 `json:"error"`
	ErrorDetails map[string]interface{} `json:"error_details,omitempty"`
}

// ObjectTermsResponse represents the response from the object terms endpoint
type ObjectTermsResponse struct {
	ObjectID   string        `json:"object_id"`
	TotalTerms int           `json:"total_terms"`
	Terms      []interface{} `json:"terms"`
}

// GlobalTermsResponse represents the response from the global terms endpoint
type GlobalTermsResponse struct {
	TotalUniqueTerms int           `json:"total_unique_terms"`
	TermsShown       int           `json:"terms_shown"`
	Terms            []interface{} `json:"terms"`
}

// ObjectsListResponse represents the response from the objects list endpoint
type ObjectsListResponse struct {
	Total   int           `json:"total"`
	Objects []interface{} `json:"objects"`
}

// DemoIndexRequest represents a request to create a demo index
type DemoIndexRequest struct {
	ID *string `json:"id,omitempty"`
}

// request is a helper function to make HTTP requests to the FuguDB API
func (c *Client) request(ctx context.Context, method, path string, body interface{}, response interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)

		if c.debug {
			fmt.Printf("Request body: %s\n", string(jsonBody))
		}
	}

	u, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	reqURL := c.baseURL.ResolveReference(u)

	if c.debug {
		fmt.Printf("Making %s request to %s\n", method, reqURL.String())
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if c.debug {
		fmt.Printf("Response status: %d\n", resp.StatusCode)
		fmt.Printf("Response body: %s\n", string(responseBody))
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(responseBody, &errResp); err != nil {
			return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(responseBody))
		}
		return fmt.Errorf("API error: %s", errResp.Error)
	}

	if response != nil {
		if err := json.Unmarshal(responseBody, response); err != nil {
			return fmt.Errorf("error unmarshaling response: %w", err)
		}
	}

	return nil
}

// HealthCheck checks if the FuguDB API is healthy
func (c *Client) HealthCheck(ctx context.Context) (bool, error) {
	var response string
	err := c.request(ctx, http.MethodGet, "/health", nil, &response)
	if err != nil {
		return false, err
	}
	return response == "OK", nil
}

// IndexObject indexes a single object in FuguDB
func (c *Client) IndexObject(ctx context.Context, object ObjectRecord) (*IndexResponse, error) {
	req := struct {
		Data ObjectRecord `json:"data"`
	}{
		Data: object,
	}

	var response IndexResponse
	err := c.request(ctx, http.MethodPost, "/ingest", req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// BatchIndexObjects indexes multiple objects in FuguDB
func (c *Client) BatchIndexObjects(ctx context.Context, objects []ObjectRecord) (*BatchIndexResponse, error) {
	req := BatchIndexRequest{
		Objects: objects,
	}

	var response BatchIndexResponse
	err := c.request(ctx, http.MethodPost, "/batch-ingest", req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// IngestFile ingests a file into FuguDB
func (c *Client) IngestFile(ctx context.Context, fileReq FileIngestionRequest) (*IndexResponse, error) {
	var response IndexResponse
	err := c.request(ctx, http.MethodPost, "/ingest-file", fileReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Search performs a search in FuguDB
func (c *Client) Search(ctx context.Context, query SearchQuery) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodPost, "/search", query, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SearchNamespace performs a search in a specific namespace in FuguDB
func (c *Client) SearchNamespace(ctx context.Context, namespace string, payload map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodPost, fmt.Sprintf("/search/%s", namespace), payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListNamespaces lists all namespaces in FuguDB
func (c *Client) ListNamespaces(ctx context.Context) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodGet, "/namespaces", nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetObject gets an object by ID from FuguDB
func (c *Client) GetObject(ctx context.Context, objectID string) (*ObjectResponse, error) {
	var response ObjectResponse
	err := c.request(ctx, http.MethodGet, fmt.Sprintf("/objects/%s", objectID), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ListObjects lists all objects in FuguDB
func (c *Client) ListObjects(ctx context.Context) (*ObjectsListResponse, error) {
	var response ObjectsListResponse
	err := c.request(ctx, http.MethodGet, "/objects", nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetObjectTerms gets the terms for an object in FuguDB
func (c *Client) GetObjectTerms(ctx context.Context, objectID string) (*ObjectTermsResponse, error) {
	var response ObjectTermsResponse
	err := c.request(ctx, http.MethodGet, fmt.Sprintf("/objects/%s/terms", objectID), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetGlobalTerms gets all terms across all objects in FuguDB
func (c *Client) GetGlobalTerms(ctx context.Context) (*GlobalTermsResponse, error) {
	var response GlobalTermsResponse
	err := c.request(ctx, http.MethodGet, "/terms", nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateDemoIndex creates a demo object index in FuguDB
func (c *Client) CreateDemoIndex(ctx context.Context, req *DemoIndexRequest) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodPost, "/demo-index", req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// QueryText queries FuguDB using a text query
func (c *Client) QueryText(ctx context.Context, query string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodGet, fmt.Sprintf("/query/%s", url.QueryEscape(query)), nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// QueryJSON queries FuguDB using a JSON query
func (c *Client) QueryJSON(ctx context.Context, query map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodPost, "/query", query, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// QueryAdvanced performs an advanced query in FuguDB
func (c *Client) QueryAdvanced(ctx context.Context, query AdvancedSearchQuery) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := c.request(ctx, http.MethodPost, "/query/advanced", query, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
