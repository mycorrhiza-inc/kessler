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

type MergePolicy struct {
	Type           string `json:"type"`
	MaxMergeOps    int    `json:"max_merge_ops"`
	MergeFactor    int    `json:"merge_factor"`
	MaxMergeFactor int    `json:"max_merge_factor"`
}

type Resources struct {
	MaxMergeWriteThroughput string `json:"max_merge_write_throughput"`
}

type IndexingSettings struct {
	MergePolicy MergePolicy `json:"merge_policy"`
	Resources   Resources   `json:"resources"`
}

type SearchSettings struct {
	DefaultSearchFields []string `json:"default_search_fields"`
}

type FieldMapping struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	InputFormats  []string `json:"input_formats,omitempty"`
	FastPrecision string   `json:"fast_precision,omitempty"`
	Fast          bool     `json:"fast"`
}

type DynamicMapping struct {
	Indexed    bool   `json:"indexed"`
	Stored     bool   `json:"stored"`
	Tokenizer  string `json:"tokenizer"`
	Record     string `json:"record"`
	ExpandDots bool   `json:"expand_dots"`
	Fast       bool   `json:"fast"`
}

type DocMapping struct {
	Mode           string         `json:"mode"`
	DynamicMapping DynamicMapping `json:"dynamic_mapping"`
	FieldMappings  []FieldMapping `json:"field_mappings"`
	TimestampField string         `json:"timestamp_field"`
}

type Retention struct {
	Period   string `json:"period"`
	Schedule string `json:"schedule"`
}

type QuickwitIndex struct {
	Version          string           `json:"version"`
	IndexID          string           `json:"index_id"`
	DocMapping       DocMapping       `json:"doc_mapping"`
	SearchSettings   SearchSettings   `json:"search_settings"`
	IndexingSettings IndexingSettings `json:"indexing_settings"`
	Retention        Retention        `json:"retention"`
}

func CreateIndex(index QuickwitIndex) error {
	jsonData, err := json.Marshal(index)
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

func CreateDocketsIndex(indexName string) error {
	requestData := QuickwitIndex{
		Version: "0.7",
		IndexID: indexName,
		DocMapping: DocMapping{
			Mode: "dynamic",
			DynamicMapping: DynamicMapping{
				Indexed:    true,
				Stored:     true,
				Tokenizer:  "default",
				Record:     "basic",
				ExpandDots: true,
				Fast:       true,
			},
			FieldMappings: []FieldMapping{
				{Name: "text", Type: "text", Fast: true},
				{Name: "name", Type: "text", Fast: true},
				{Name: "verified", Type: "bool", Fast: true},
				{Name: "timestamp", Type: "datetime", InputFormats: []string{"unix_timestamp"}, FastPrecision: "seconds", Fast: true},
				{Name: "date_filed", Type: "datetime", Fast: true},
			},
			TimestampField: "timestamp",
		},
		SearchSettings: SearchSettings{
			DefaultSearchFields: []string{"text", "state", "city", "country"},
		},
		IndexingSettings: IndexingSettings{
			MergePolicy: MergePolicy{
				Type:           "limit_merge",
				MaxMergeOps:    3,
				MergeFactor:    10,
				MaxMergeFactor: 12,
			},
			Resources: Resources{
				MaxMergeWriteThroughput: "80mb",
			},
		},
		Retention: Retention{
			Period:   "10 years",
			Schedule: "yearly",
		},
	}
	err := CreateIndex(requestData)
	return err

}

func ClearIndex(indexName string, nuke bool) error {
	var req *http.Request
	var err error
	if nuke {
		req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/indexes/%s", quickwitEndpoint, indexName), nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
	} else {
		req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/indexes/%s/clear", quickwitEndpoint, indexName), nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
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
	fmt.Println("Ingesting %d initial data entries into index\n", len(data))
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
			// Do nothing, an error here means to text was found.
		}
		if englishText == "" {
			englishText = "Example Text!"
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
	// fmt.Printf("len(data) is %d, len(complete_files) is %d\n", len(data), len(complete_files))
	if len(data) != len(complete_files) {
		fmt.Printf("len(data) is %d, len(complete_files) is %d. THEY ARE NOT EUQAL THIS SHOULD NEVER HAPPEN\n", len(data), len(complete_files))
		return nil, errors.New("ASSERTION ERROR: len(data) is not equal to len(complete_files), dispite no logical way for that to happen!")
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
