package quickwit

import (
	"encoding/json"
	"fmt"
	"kessler/objects/conversations"
	"kessler/objects/organizations"
	"net/http"
	"strings"
)

func IngestIntoIndex[V QuickwitFileUploadData | conversations.ConversationInformation | organizations.OrganizationSchemaComplete](indexName string, data []V) error {
	maxIngestItems := 1000
	fmt.Println("Initiating ingest into index")

	for i := 0; i < len(data); i += maxIngestItems {
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
			fmt.Sprintf("%s/api/v1/%s/ingest", quickwitEndpoint, indexName),
			"application/x-ndjson",
			strings.NewReader(dataToPost),
		)
		if err != nil {
			return fmt.Errorf("error submitting data to quickwit: %v", err)
		}
		defer resp.Body.Close()

		printResponse(resp)
	}
	tailUrl := fmt.Sprintf("%s/api/v1/%s/tail", quickwitEndpoint, indexName)
	resp, err := http.Get(tailUrl)
	if err != nil {
		fmt.Printf("Error tailing index: %v\n", err)
	}
	printResponse(resp)
	return nil
}
