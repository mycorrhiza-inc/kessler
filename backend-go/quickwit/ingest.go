package quickwit

import (
	"encoding/json"
	"fmt"
	"kessler/objects/conversations"
	"kessler/objects/organizations"
	"net/http"
	"strings"
)

func IngestIntoIndex[V QuickwitFileUploadData | conversations.ConversationInformation | organizations.OrganizationSchemaComplete](indexName string, data []V, clear_index bool) error {
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

	// for i := 0; i < len(data); i += maxIngestItems {
	for i := 0; i < maxIngestItems; i += maxIngestItems {
		end := i + maxIngestItems
		if end > len(data) {
			end = len(data)
		}

		var records []string
		for _, record := range data[i:end] {
			jsonStr, err := json.Marshal(record)
			if err != nil {
				return fmt.Errorf("error marshaling records for index: %v", err)
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

	describeUrl := fmt.Sprintf("%s/api/v1/%s/describe", quickwitEndpoint, indexName)
	resp, err = http.Get(describeUrl)
	if err != nil {
		return fmt.Errorf("Error describing index: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Describing index gave bad response code: %v\n", resp.StatusCode)
	}

	body = make([]byte, 1000)
	n, err = resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("Error reading response body: %v", err)
	}
	fmt.Printf("Describe response (first %d chars): %s\n", n, string(body[:n]))

	return nil
}
