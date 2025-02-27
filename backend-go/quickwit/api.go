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

// QuickwitClient represents a Quickwit APIQuickwitClient
type QuickwitClient struct {
	baseURL    string
	httpClient *http.Client
	ctx        context.Context
}

// NewClient creates a new QuickwitQuickwitClient
func NewClient(baseURL string, ctx context.Context) QuickwitClient {
	// TODO: good logging with traces
	qwCtx := ctx

	return QuickwitClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		ctx: qwCtx,
	}
}

// SearchParams represents parameters for a search request
type SearchParams struct {
	Query          string          `json:"query"`
	StartTimestamp string          `json:"start_timestamp,omitempty"`
	EndTimestamp   string          `json:"end_timestamp,omitempty"`
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
func (c QuickwitClient) Search(indexID string, params SearchParams) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/%s/search", url.PathEscape(indexID))

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshaling search params: %w", err)
	}

	req, err := http.NewRequestWithContext(c.ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
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
func (c QuickwitClient) CreateIndex(config IndexConfig) error {
	endpoint := "/api/v1/indexes"

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling index config: %w", err)
	}

	req, err := http.NewRequestWithContext(c.ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
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

// IndexMetadata represents the metadata of an index
type IndexMetadata struct {
	IndexConfig     IndexConfig                `json:"index_config"`
	Checkpoint      map[string]json.RawMessage `json:"checkpoint"`
	CreateTimestamp int64                      `json:"create_timestamp"`
	Sources         []json.RawMessage          `json:"sources"`
}

// GetIndexMetadata retrieves metadata for a specific index
func (c QuickwitClient) GetIndexMetadata(indexID string) (*IndexMetadata, error) {
	endpoint := fmt.Sprintf("/api/v1/indexes/%s", url.PathEscape(indexID))

	req, err := http.NewRequestWithContext(c.ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var metadata IndexMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &metadata, nil
}

// ValidateFieldName checks if a field name is valid and exists in the given index
func (c QuickwitClient) ValidateFieldName(indexID string, fieldName string) error {
	if fieldName == "" {
		return fmt.Errorf("field name cannot be empty")
	}

	metadata, err := c.GetIndexMetadata(indexID)
	if err != nil {
		return fmt.Errorf("getting index metadata: %w", err)
	}

	// Parse the doc_mapping to check field existence
	var docMapping struct {
		FieldMappings []struct {
			Name string `json:"name"`
		} `json:"field_mappings"`
	}

	if err := json.Unmarshal(metadata.IndexConfig.DocMapping, &docMapping); err != nil {
		return fmt.Errorf("parsing doc mapping: %w", err)
	}

	// Check if the field exists
	for _, field := range docMapping.FieldMappings {
		if field.Name == fieldName {
			return nil // Field found
		}
	}

	return fmt.Errorf("field '%s' not found in index '%s'", fieldName, indexID)
}

// DeleteTask represents a delete task configuration
type DeleteTask struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_fields,omitempty"`
	StartTimestamp string   `json:"start_timestamp,omitempty"`
	EndTimestamp   string   `json:"end_timestamp,omitempty"`
}

// CreateDeleteTask creates a new delete task for the specified index
func (c QuickwitClient) CreateDeleteTask(indexID string, task DeleteTask) error {
	endpoint := fmt.Sprintf("/api/v1/%s/delete-tasks", url.PathEscape(indexID))

	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshaling delete task: %w", err)
	}

	req, err := http.NewRequestWithContext(c.ctx, "POST", c.baseURL+endpoint, bytes.NewReader(body))
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
func (c QuickwitClient) IngestDocuments(indexID string, documents []interface{}, commitBehavior string) error {
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

	req, err := http.NewRequestWithContext(c.ctx, "POST", c.baseURL+endpoint, &buf)
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
func (c QuickwitClient) DeleteIndex(indexID string) error {
	endpoint := fmt.Sprintf("/api/v1/indexes/%s", url.PathEscape(indexID))

	req, err := http.NewRequestWithContext(c.ctx, "DELETE", c.baseURL+endpoint, nil)
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
