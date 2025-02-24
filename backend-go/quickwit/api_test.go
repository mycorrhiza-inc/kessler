package quickwit

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const (
	testIndexID = "api_test_index"
	testBaseURL = "http://100.86.5.114:7280"
)

func TestEndToEndOperations(t *testing.T) {
	client := NewClient(testBaseURL, context.Background())

	// Step 1: Create index
	t.Log("Creating test index...")
	err := createTestIndex(&client)
	if err != nil {
		t.Fatalf("Failed to create index: %v", err)
	}

	// Step 2: Ingest test documents
	t.Log("Ingesting test documents...")
	err = ingestTestDocuments(&client)
	if err != nil {
		t.Fatalf("Failed to ingest documents: %v", err)
	}

	// Wait for documents to be searchable
	time.Sleep(2 * time.Second)

	// Step 3: Search documents
	t.Log("Searching documents...")
	err = searchTestDocuments(&client)
	if err != nil {
		t.Fatalf("Failed to search documents: %v", err)
	}

	// Step 4: Create delete task
	t.Log("Creating delete task...")
	err = createTestDeleteTask(&client)
	if err != nil {
		t.Fatalf("Failed to create delete task: %v", err)
	}

	// Wait for deletion to take effect
	time.Sleep(2 * time.Second)

	// Step 5: Verify deletion with search
	t.Log("Verifying deletion...")
	err = verifyDeletion(&client)
	if err != nil {
		t.Fatalf("Failed to verify deletion: %v", err)
	}

	// Final Step: Delete index
	t.Log("Cleaning up - deleting test index...")
	err = client.DeleteIndex(testIndexID)
	if err != nil {
		t.Fatalf("Failed to delete index: %v", err)
	}
}

func createTestIndex(client *QuickwitClient) error {
	// Create a basic index configuration
	docMapping := map[string]interface{}{
		"field_mappings": []map[string]interface{}{
			{
				"name": "title",
				"type": "text",
				"fast": true,
			},
			{
				"name":   "body",
				"type":   "text",
				"record": "position",
			},
			{
				"name":           "timestamp",
				"type":           "datetime",
				"fast":           true,
				"input_formats":  []string{"unix_timestamp"},
				"fast_precision": "seconds",
			},
		},
		"timestamp_field": "timestamp",
	}

	docMappingJSON, err := json.Marshal(docMapping)
	if err != nil {
		return err
	}

	searchSettings := map[string]interface{}{
		"default_search_fields": []string{"title", "body"},
	}
	searchSettingsJSON, err := json.Marshal(searchSettings)
	if err != nil {
		return err
	}

	config := IndexConfig{
		Version:        "0.7",
		IndexID:        testIndexID,
		DocMapping:     docMappingJSON,
		SearchSettings: searchSettingsJSON,
	}

	return client.CreateIndex(config)
}

func ingestTestDocuments(client *QuickwitClient) error {
	now := time.Now().Format(time.RFC3339)
	docs := []interface{}{
		map[string]interface{}{
			"title":     "Test Document 1",
			"body":      "This is a test document with some content to be deleted",
			"timestamp": now,
		},
		map[string]interface{}{
			"title":     "Test Document 2",
			"body":      "This is another test document that should remain",
			"timestamp": now,
		},
	}

	return client.IngestDocuments(testIndexID, docs, "wait_for")
}

func searchTestDocuments(client *QuickwitClient) error {
	maxHits := 10
	params := SearchParams{
		Query:   "test",
		MaxHits: &maxHits,
	}

	resp, err := client.Search(testIndexID, params)
	if err != nil {
		return err
	}

	if resp.NumHits != 2 {
		return fmt.Errorf("expected 2 hits, got %d", resp.NumHits)
	}

	return nil
}

func createTestDeleteTask(client *QuickwitClient) error {
	task := DeleteTask{
		Query: "body:deleted",
	}

	return client.CreateDeleteTask(testIndexID, task)
}

func verifyDeletion(client *QuickwitClient) error {
	maxHits := 10
	params := SearchParams{
		Query:   "body:deleted",
		MaxHits: &maxHits,
	}

	resp, err := client.Search(testIndexID, params)
	if err != nil {
		return err
	}

	if resp.NumHits != 0 {
		return fmt.Errorf("expected 0 hits after deletion, got %d", resp.NumHits)
	}

	return nil
}
