package quickwit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kessler/objects/files"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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
				{
					"name": "verified",
					"type": "bool",
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

func CreateRFC3339FromString(dateStr string) (string, error) {
	if dateStr == "" {
		return "", errors.New("empty date string")
	}
	dateParts := strings.Split(dateStr, "/")
	if len(dateParts) != 3 {
		return "", errors.New("date string must be in the format MM/DD/YYYY")
	}
	month := dateParts[0]
	day := dateParts[1]
	year := dateParts[2]

	parsedDate, err := time.Parse("01/02/2006", fmt.Sprintf("%s/%s/%s", month, day, year))
	if err != nil {
		return "", err
	}
	return parsedDate.Format(time.RFC3339), nil
}

func ResolveFileSchemaForDocketIngest(complete_files []files.CompleteFileSchema) ([]QuickwitFileUploadData, error) {
	createEnrichedMetadata := func(input_file files.CompleteFileSchema) map[string]interface{} {
		metadata := input_file.Mdata
		metadata["source_id"] = input_file.ID
		metadata["source"] = "ny-puc-energyefficiency-filedocs"
		metadata["conversation_uuid"] = input_file.Conversation.ID.String()
		author_uuids := make([]uuid.UUID, len(input_file.Authors))
		for i, author := range input_file.Authors {
			author_uuids[i] = author.AuthorID
		}
		metadata["author_uuids"] = author_uuids
		// FIXME: IMPLEMENT A date_published FIELD IN PG AND RENDER THIS BASED ON THAT
		dateStr := metadata["date"].(string)
		parsedDate, err := CreateRFC3339FromString(dateStr)
		if err != nil {
			metadata["date_filed"] = parsedDate
		}
		return metadata
	}
	var data []QuickwitFileUploadData
	for _, file := range complete_files {
		newRecord := QuickwitFileUploadData{}
		englishText, err := files.EnglishTextFromCompleteFile(file)
		if err != nil {
			continue
		}
		if englishText == "" {
			englishText = "No English text found in file, this is some example text so quickwit doesnt exclude it, please ignore."
		}
		newRecord.Text = englishText
		newRecord.SourceID = file.ID
		newRecord.Verified = file.Verified

		newRecord.Metadata = createEnrichedMetadata(file)
		newRecord.Name = file.Name
		date, err := CreateRFC3339FromString(file.Mdata["date"].(string))
		if err == nil {
			newRecord.DateFiled = date
		}

		newRecord.Timestamp = time.Now().Unix()
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
