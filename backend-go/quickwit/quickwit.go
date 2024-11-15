package quickwit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var quickwitEndpoint = os.Getenv("QUICKWIT_ENDPOINT")

func printResponse(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	log.Printf("quickwit_api_call:\nstatus: %d\nresponse:\n%s", resp.StatusCode, string(body))
}

func CreateDocketsQuickwitIndex(indexName string) error {
	requestData := map[string]interface{}{
		"version":  "0.7",
		"index_id": indexName,
		"doc_mapping": map[string]interface{}{
			"mode": "dynamic",
			"dynamic_mapping": map[string]interface{}{
				"indexed":     true,
				"stored":      true,
				"tokenizer":   "default",
				"record":      "basic",
				"expand_dots": true,
				"fast":        true,
			},
			"field_mappings": []map[string]interface{}{
				{
					"name": "text",
					"type": "text",
					"fast": true,
				},
				{
					"name":           "timestamp",
					"type":           "datetime",
					"input_formats":  []string{"unix_timestamp"},
					"fast_precision": "seconds",
					"fast":           true,
				},
				{
					"name": "date_filed",
					"type": "datetime",
					"fast": true,
				},
			},
			"timestamp_field": "timestamp",
		},
		"search_settings": map[string]interface{}{
			"default_search_fields": []string{
				"text", "state", "city", "country",
			},
		},
		"indexing_settings": map[string]interface{}{
			"merge_policy": map[string]interface{}{
				"type":             "limit_merge",
				"max_merge_ops":    3,
				"merge_factor":     10,
				"max_merge_factor": 12,
			},
			"resources": map[string]interface{}{
				"max_merge_write_throughput": "80mb",
			},
		},
		"retention": map[string]interface{}{
			"period":   "10 years",
			"schedule": "yearly",
		},
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshaling request data: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/indexes", quickwitEndpoint),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
	return nil
}

func ClearIndex(indexName string) error {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/indexes/%s/clear", quickwitEndpoint, indexName), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
	return nil
}

func IngestIntoIndex(indexName string, data []map[string]interface{}) error {
	fmt.Println("Initiating ingest into index")
	var records []string
	for _, record := range data {
		jsonStr, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("error marshaling records for index: %v", err)
		}
		records = append(records, string(jsonStr))
	}

	dataToPost := strings.Join(records, "\n")
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/%s/ingest", quickwitEndpoint, indexName),
		"application/x-ndjson",
		strings.NewReader(dataToPost),
	)
	if err != nil {
		return fmt.Errorf("error submitting data to quickwit: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
	return nil
}

func ResolveFileSchemaForDocketIngest(records []map[string]interface{}) []map[string]interface{} {
	var data []map[string]interface{}
	for _, record := range records {
		newRecord := make(map[string]interface{})
		newRecord["text"] = record["english_text"]
		newRecord["source_id"] = record["id"]

		if metadata, ok := record["mdata"].(string); ok {
			var parsedMetadata interface{}
			if err := json.Unmarshal([]byte(metadata), &parsedMetadata); err == nil {
				newRecord["metadata"] = parsedMetadata
			}
		}

		newRecord["timestamp"] = time.Now().Unix()
		data = append(data, newRecord)
	}

	log.Printf("reformatted data:\n\n%+v\n\n", data)
	return data
}

func MigrateDocketToNYPUC() error {
	query := "metadata.source:(ny-puc-energyefficiency-filedocs)"
	requestData := map[string]interface{}{
		"query":        query,
		"start_offset": 0,
		"max_hits":     10,
		"limit":        1000,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshaling request data: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/dockets/search", quickwitEndpoint),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
	return nil
}
