package quickwit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kessler/objects/files"
	"kessler/objects/timestamp"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
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

func CreateQuickwitProceedingIndex(indexName string) error {
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
				{Name: "name", Type: "text", Fast: true},
				{Name: "caseNumber", Type: "text", Fast: true},
				{Name: "description", Type: "text", Fast: true},
				{Name: "uuid", Type: "text", Fast: true},
			},
		},
		SearchSettings: SearchSettings{
			DefaultSearchFields: []string{"name"},
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

func ResolveFileSchemaForDocketIngest(complete_files []files.CompleteFileSchema) ([]QuickwitFileUploadData, error) {
	createEnrichedMetadata := func(input_file files.CompleteFileSchema) map[string]interface{} {
		metadata := input_file.Mdata
		metadata["source_id"] = input_file.ID
		metadata["source"] = "ny-puc-energyefficiency-filedocs"
		metadata["conversation_uuid"] = input_file.Conversation.ID.String()
		author_uuids := make([]string, len(input_file.Authors))
		if len(input_file.Authors) == 0 {
			log.Info(fmt.Sprintf("No authors found in file: %v\n", input_file.ID))
		}
		for i, author := range input_file.Authors {
			author_uuids[i] = author.AuthorID.String()
		}
		metadata["author_uuids"] = author_uuids
		// FIXME: IMPLEMENT A date_published FIELD IN PG AND RENDER THIS BASED ON THAT
		dateStr := metadata["date"].(string)
		parsedDate, err := timestamp.CreateRFC3339FromString(dateStr)
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
			englishText = ""
		}
		newRecord.Text = englishText
		newRecord.SourceID = file.ID
		newRecord.Verified = file.Verified

		newRecord.Metadata = createEnrichedMetadata(file)
		newRecord.Name = file.Name
		date, err := timestamp.CreateRFC3339FromString(file.Mdata["date"].(string))
		if err == nil {
			newRecord.DateFiled = date
		}

		newRecord.Timestamp = time.Now().Unix()
		data = append(data, newRecord)
	}
	// log.Info(fmt.Sprintf("len(data) is %d, len(complete_files) is %d\n", len(data), len(complete_files)))
	if len(data) != len(complete_files) {
		log.Info(fmt.Sprintf("len(data) is %d, len(complete_files) is %d. THEY ARE NOT EUQAL THIS SHOULD NEVER HAPPEN\n", len(data), len(complete_files)))
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

func CreateDeleteTask(indexName string, task DeleteTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("error marshalling delete task request: %v\n\ntask: %s", err, task)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/%s/delete-tasks", quickwitEndpoint, indexName),
		"application/x-ndjson",
		strings.NewReader(string(data)),
	)
	if err != nil {
		return fmt.Errorf("error submitting delete task request: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
	return nil
}
