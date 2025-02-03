package quickwit

import (
	"encoding/json"
	"fmt"
	"kessler/objects/conversations"
	"kessler/objects/organizations"
	"net/http"
	"strings"
	"time"
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
	maxIngestItems := 1000
	fmt.Println("Initiating ingest into index")

	for i := 0; i < len(data); i += maxIngestItems {
		// for i := 0; i < maxIngestItems; i += maxIngestItems {
		end := i + maxIngestItems
		if end > len(data) {
			end = len(data)
		}

		var records []string
		for _, record := range data[i:end] {
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
				return fmt.Errorf("error marshaling modified record: %v", err)
			}
			records = append(records, string(jsonStr))
		}

		dataToPost := strings.Join(records, "\n")
		fmt.Printf("Ingesting %d data entries into quickwit index:\"%v\"(batch %d-%d) \n", len(records), indexName, i+1, end)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/%s/ingest?commit=force", quickwitEndpoint, indexName),
			"application/x-ndjson",
			strings.NewReader(dataToPost),
		)
		if err != nil {
			return fmt.Errorf("error submitting data to quickwit: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Ingesting data into quickwit gave bad response code: %v\n", resp.StatusCode)
		}
		defer resp.Body.Close()

		printResponse(resp)
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
	fmt.Printf("Tail response (first %d chars): %s\n", n, string(body[:n]))

	// describeUrl := fmt.Sprintf("%s/api/v1/%s/describe", quickwitEndpoint, indexName)
	// resp, err = http.Get(describeUrl)
	// if err != nil {
	// 	return fmt.Errorf("Error describing index: %v\n", err)
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("Describing index gave bad response code: %v\n", resp.StatusCode)
	// }

	body = make([]byte, 1000)
	n, err = resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("Error reading response body: %v", err)
	}
	fmt.Printf("Describe response (first %d chars): %s\n", n, string(body[:n]))

	return nil
}
