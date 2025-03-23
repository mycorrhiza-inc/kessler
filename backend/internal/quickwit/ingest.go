package quickwit

import (
	"encoding/json"
	"fmt"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/organizations"
	"kessler/pkg/util"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type GenericQuickwitSearchSchema interface {
	QuickwitFileUploadData | conversations.ConversationInformation | organizations.OrganizationQuickwitSchema
}

func IngestIntoIndex[V GenericQuickwitSearchSchema](indexName string, data []V, clear_index bool) error {
	if clear_index {
		clear_url := fmt.Sprintf("%s/api/v1/indexes/%s/clear", quickwitEndpoint, indexName)
		req, err := http.NewRequest(http.MethodPut, clear_url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error clearing index: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error clearing index, got bad status code: %v", resp.StatusCode)
		}
		defer resp.Body.Close()
	}
	maxIngestItems := 100
	log.Info("Initiating ingest into index")

	var subIngestLists [][]V

	for i := 0; i < len(data); i += maxIngestItems {
		end := i + maxIngestItems
		if end > len(data) {
			end = len(data)
		}
		subIngestLists = append(subIngestLists, data[i:end])
	}
	ingestWrapedFunc := func(data []V) (int, error) {
		return 0, IngestMinimalIntoQuickwit(indexName, data)
	}
	workers := 15
	_, err := util.ConcurrentMapError(subIngestLists, ingestWrapedFunc, workers)
	if err != nil {
		return fmt.Errorf("Error ingesting into index: %v\n", err)
	}

	tailUrl := fmt.Sprintf("%s/api/v1/%s/tail", quickwitEndpoint, indexName)
	resp, err := http.Get(tailUrl)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("Error tailing index: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Tailing index gave bad response code: %v\n", resp.StatusCode)
	}

	body := make([]byte, 1000)
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("Error reading response body: %v", err)
	}
	log.Info(fmt.Sprintf("Tail response (first %d chars): %s\n", n, string(body[:n])))

	body = make([]byte, 1000)
	n, err = resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("Error reading response body: %v", err)
	}
	log.Info(fmt.Sprintf("Describe response (first %d chars): %s\n", n, string(body[:n])))

	return nil
}

func IngestMinimalIntoQuickwit[V GenericQuickwitSearchSchema](indexName string, data []V) error {
	records := make([]string, 0, len(data))
	for _, record := range data {
		// Convert the record to a map to check for the timestamp field
		var recordMap map[string]interface{}
		jsonBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("error marshaling record: %v", err)
		}
		if err := json.Unmarshal(jsonBytes, &recordMap); err != nil {
			return fmt.Errorf("error unmarshaling record into map: %v", err)
		}

		// Check if the timestamp field exists; if not, set it
		if _, exists := recordMap["timestamp"]; !exists {
			timestamp_value := time.Now().UTC().Unix()
			// recordMap["timestamp"] = timestamp.KesslerTime(time.Now())
			recordMap["timestamp"] = timestamp_value
		}

		// Marshal the modified map into a JSON string
		jsonStr, err := json.Marshal(recordMap)
		if err != nil {
			log.Info(fmt.Sprintf("error marshaling modified record: %v\n", err))
			return fmt.Errorf("error marshaling modified record: %v", err)
		}
		records = append(records, string(jsonStr))
	}

	dataToPost := strings.Join(records, "\n")
	log.Info(fmt.Sprintf("Ingesting %d data entries into quickwit index:\"%v\" \n", len(records), indexName))
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/%s/ingest?commit=force", quickwitEndpoint, indexName),
		"application/x-ndjson",
		strings.NewReader(dataToPost),
	)
	if err != nil {
		log.Info(err)
		return fmt.Errorf("error submitting data to quickwit: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Info(fmt.Sprintf("Encountered bad response code: %v\n", resp.StatusCode))
		return fmt.Errorf("Ingesting data into quickwit gave bad response code: %v\n", resp.StatusCode)
	}
	printResponse(resp)
	defer resp.Body.Close()

	return nil
}
