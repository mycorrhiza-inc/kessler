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
	fmt.Printf("Ingesting %d initial data entries into index\n", len(data))
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
