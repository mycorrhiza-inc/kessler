package quickwit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a Quickwit API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Quickwit client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// SearchParams represents parameters for a search request
type SearchParams struct {
	Query          string          `json:"query"`
	StartTimestamp *int64          `json:"start_timestamp,omitempty"`
	EndTimestamp   *int64          `json:"end_timestamp,omitempty"`
	StartOffset    *int            `json:"start_offset,omitempty"`
	MaxHits        *int            `json:"max_hits,omitempty"`
	SearchFields   []string        `json:"search_fields,omitempty"`
	SnippetFields  []string        `json:"snippet_fields,omitempty"`
	SortBy         []string        `json:"sort_by,omitempty"`
	Format         string          `json:"format,omitempty"`
	Aggregations   json.RawMessage `json:"aggs,omitempty"`
}

// SearchResponse represents the response from a search request
type SearchResponse struct {
	Hits              []json.RawMessage `json:"hits"`
	NumHits           int               `json:"num_hits"`
	ElapsedTimeMicros int64             `json:"elapsed_time_micros"`
}

// Search performs a search request on a single index
func (c *Client) Search(ctx context.Context, indexID string, params SearchParams) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/%s/search", url.PathEscape(indexID))

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshaling search params: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}

// IndexConfig represents the configuration for creating a new index
type IndexConfig struct {
	Version          string          `json:"version"`
	IndexID          string          `json:"index_id"`
	IndexURI         string          `json:"index_uri,omitempty"`
	DocMapping       json.RawMessage `json:"doc_mapping"`
	IndexingSettings json.RawMessage `json:"indexing_settings,omitempty"`
	SearchSettings   json.RawMessage `json:"search_settings,omitempty"`
	Retention        json.RawMessage `json:"retention,omitempty"`
}

// CreateIndex creates a new index with the given configuration
func (c *Client) CreateIndex(ctx context.Context, config IndexConfig) error {
	endpoint := "/api/v1/indexes"

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling index config: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteTask represents a delete task configuration
type DeleteTask struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_fields,omitempty"`
	StartTimestamp *int64   `json:"start_timestamp,omitempty"`
	EndTimestamp   *int64   `json:"end_timestamp,omitempty"`
}

// CreateDeleteTask creates a new delete task for the specified index
func (c *Client) CreateDeleteTask(ctx context.Context, indexID string, task DeleteTask) error {
	endpoint := fmt.Sprintf("/api/v1/%s/delete-tasks", url.PathEscape(indexID))

	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshaling delete task: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// IngestDocuments ingests documents into the specified index
func (c *Client) IngestDocuments(ctx context.Context, indexID string, documents []interface{}, commitBehavior string) error {
	endpoint := fmt.Sprintf("/api/v1/%s/ingest", url.PathEscape(indexID))
	if commitBehavior != "" {
		endpoint += "?commit=" + url.QueryEscape(commitBehavior)
	}

	// Convert documents to NDJSON format
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	for _, doc := range documents {
		if err := encoder.Encode(doc); err != nil {
			return fmt.Errorf("encoding document: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, &buf)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteIndex deletes an index
func (c *Client) DeleteIndex(ctx context.Context, indexID string) error {
	endpoint := fmt.Sprintf("/api/v1/indexes/%s", url.PathEscape(indexID))

	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
