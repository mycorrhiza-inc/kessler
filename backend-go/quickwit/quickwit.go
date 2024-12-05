package quickwit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kessler/objects/files"
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

func IngestIntoIndex(indexName string, data []QuickwitFileUploadData) error {
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

func ResolveFileSchemaForDocketIngest(complete_files []files.CompleteFileSchema) ([]QuickwitFileUploadData, error) {
	createEnrichedMetadata := func(input_file files.CompleteFileSchema) map[string]interface{} {
		metadata := input_file.Mdata
		metadata["source_id"] = input_file.ID
		metadata["source"] = "ny-puc-energyefficiency-filedocs"
		metadata["conversation_uuid"] = input_file.Conversation.ID
		return metadata
	}
	var data []QuickwitFileUploadData
	for _, file := range complete_files {
		newRecord := make(map[string]interface{})
		englishText, err := files.EnglishTextFromCompleteFile(file)
		if err != nil {
			continue
		}
		newRecord["text"] = englishText
		newRecord["source_id"] = file.ID

		newRecord["metadata"] = createEnrichedMetadata(file)
		newRecord["name"] = file.Name

		newRecord["timestamp"] = time.Now().Unix()
		data = append(data, newRecord)
	}

	// log.Printf("reformatted data:\n\n%+v\n\n", data)
	return data, nil
}

func MigrateDocketToNYPUC() error {
	query := "metadata.source:(ny-puc-energyefficiency-filedocs)"
	requestData := map[string]interface{}{
		"query":        query,
		"start_offset": 0,
		"max_hits":     40,
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
