package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kessler/objects/authors"
	"kessler/objects/conversations"
	"kessler/objects/files"
	"kessler/objects/networking"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

var quickwitURL = os.Getenv("QUICKWIT_ENDPOINT")

type Hit struct {
	CreatedAt     string              `json:"created_at"`
	Doctype       string              `json:"doctype"`
	Hash          string              `json:"hash"`
	Lang          string              `json:"lang"`
	DateFiled     string              `json:"updated_at"`
	Metadata      networking.Metadata `json:"metadata"`
	Name          string              `json:"name"`
	SaOrmSentinel *string             `json:"sa_orm_sentinel"`
	ShortSummary  *string             `json:"short_summary"`
	Source        string              `json:"source"`
	SourceID      string              `json:"source_id"`
	Stage         string              `json:"stage"`
	Summary       *string             `json:"summary"`
	Text          string              `json:"text"`
	Timestamp     string              `json:"timestamp"`
	URL           string              `json:"url"`
}

type Snippet struct {
	Text []string `json:"text"`
}

type quickwitSearchResponse struct {
	Hits     []Hit     `json:"hits"`
	Snippets []Snippet `json:"snippets"`
}

func (s quickwitSearchResponse) String() string {
	// Marshal the struct to JSON format
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
}

type QuickwitSearchRequest struct {
	Query         string `json:"query,omitempty"`
	SnippetFields string `json:"snippet_fields,omitempty"`
	MaxHits       int    `json:"max_hits"`
	StartOffset   int    `json:"start_offset"`
	SortBy        string `json:"sort_by,omitempty"`
}

type SearchData struct {
	Name     string `json:"name"`
	Snippet  string `json:"snippet,omitempty"`
	DocID    string `json:"docID"`
	SourceID string `json:"sourceID"`
}

func (s SearchData) String() string {
	// Marshal the struct to JSON format
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
}

func SearchQuickwit(r SearchRequest) ([]SearchDataHydrated, error) {
	fmt.Printf("searching quickwit:\n%s", r)
	r.Index = "NY_PUC"
	search_index := r.Index
	// ===== construct search request =====
	query := r.Query

	var queryFilters networking.FilterFields = r.SearchFilters
	var metadataFilters networking.MetadataFilterFields = queryFilters.MetadataFilters
	var uuidFilters networking.UUIDFilterFields = queryFilters.UUIDFilters
	log.Printf("zzxxcc: %v\n", uuidFilters)

	var queryString string
	dateQuery, err := ConstructDateQuery(metadataFilters.DateFrom, metadataFilters.DateTo)
	if err != nil {
		log.Printf("error constructing date query: %v", err)
	}
	if len(r.Query) >= 0 {
		queryString = fmt.Sprintf("((text:(%s) OR name:(%s)) AND %s)", query, query, dateQuery)
		// queryString = fmt.Sprintf("((text:(%s) OR name:(%s)) AND verified:true AND %s)", query, query, dateQuery)
	}

	filtersString := constructQuickwitMetadataQueryString(metadataFilters.Metadata)
	uuidFilterString := constructQuickwitUUIDMetadataQueryString(uuidFilters)

	log.Printf(
		"!!!!!!!!!!\nquery: %s\nfilters: %s\nuuid filters: %s\n!!!!!!!!!!\n",
		queryString,
		filtersString,
		uuidFilterString,
	)
	queryString = queryString + filtersString + uuidFilterString
	log.Printf("full query string: %s\n", queryString)

	// construct sortby string
	sortbyStr := "date_filed"
	// Changing this to be a greater than or equal to 1, instead of a less than or equal to 1, trying to track down an out of index thing - Nic
	if len(r.SortBy) >= 1 {
		sortbyStr = r.SortBy[0]
	} else {
		for _, sortby := range r.SortBy {
			sortbyStr += fmt.Sprintf("metadata.%s,", sortby)
		}
	}
	snippetFields := "text"
	if !r.GetText {
		snippetFields = ""
	}

	if r.MaxHits == 0 {
		r.MaxHits = 40
	}
	// construct request
	request := QuickwitSearchRequest{
		Query:         queryString,
		MaxHits:       r.MaxHits,
		SnippetFields: snippetFields,
		StartOffset:   r.StartOffset,
		SortBy:        sortbyStr,
	}

	jsonData, err := json.Marshal(request)

	// ===== submit request to quickwit =====
	log.Printf("jsondata: \n%s", jsonData)
	if err != nil {
		log.Printf("Error Marshalling quickwit request: %s", err)
		return []SearchDataHydrated{}, err
	}

	request_url := fmt.Sprintf("%s/api/v1/%s/search", quickwitURL, search_index)
	resp, err := http.Post(
		request_url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Error sending request to quickwit: %s", err)
		errturn(err)
	}

	defer resp.Body.Close()

	// ===== handle response =====
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %v", resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error: received status code %v", resp.StatusCode)
			return []SearchDataHydrated{}, fmt.Errorf("received status code %v", resp.StatusCode)
		}
	}
	var searchResponse quickwitSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		log.Printf("Error unmarshalling quickwit response: %s", err)
		return []SearchDataHydrated{}, err
	}
	fmt.Printf("quickwit response: %v\n", resp)

	data, err := ExtractSearchData(searchResponse)
	if err != nil {
		log.Printf("Error creating response data: %s", err)
		return []SearchDataHydrated{}, err
	}

	return data, nil
}

// Function to create search data array
func ExtractSearchData(data quickwitSearchResponse) ([]SearchDataHydrated, error) {
	var result []SearchDataHydrated

	// Map snippets text to hit names
	for i, hit := range data.Hits {
		var snippet string
		if len(data.Snippets) > i {
			if len(data.Snippets[i].Text) > 0 {
				snippet = data.Snippets[i].Text[0]
			}
		}
		author_infos := make([]authors.AuthorInformation, len(hit.Metadata.AuthorUUIDs))
		for index, author_id := range hit.Metadata.AuthorUUIDs {
			name := ""
			if len(hit.Metadata.Authors) > index {
				name = hit.Metadata.Authors[index]
			}
			author_infos[index] = authors.AuthorInformation{
				AuthorID:   author_id,
				AuthorName: name,
			}
		}
		file_id, err := uuid.Parse(hit.SourceID)
		if err != nil {
			return []SearchDataHydrated{}, err
		}
		convo_id := hit.Metadata.ConversationUUID
		convo_info := conversations.ConversationInformation{
			DocketGovID: hit.Metadata.DocketID,
			ID:          convo_id,
		}
		file_schema := files.CompleteFileSchema{
			ID:           file_id,
			Authors:      author_infos,
			Conversation: convo_info,
		}
		sdata := SearchDataHydrated{
			Name:     hit.Name,
			Snippet:  snippet,
			DocID:    hit.Metadata.DocketID,
			SourceID: hit.SourceID,
			File:     file_schema,
		}
		result = append(result, sdata)
	}

	return result, nil
}

func ExtractSearchDataPlain(data quickwitSearchResponse) ([]SearchData, error) {
	var result []SearchData

	// Map snippets text to hit names
	for i, hit := range data.Hits {
		var snippet string
		if len(data.Snippets) > i {
			if len(data.Snippets[i].Text) > 0 {
				snippet = data.Snippets[i].Text[0]
			}
		}
		sdata := SearchData{
			Name:     hit.Name,
			Snippet:  snippet,
			DocID:    hit.Metadata.DocketID,
			SourceID: hit.SourceID,
		}
		result = append(result, sdata)
	}

	return result, nil
}

func errturn(err error) ([]SearchData, error) {
	return nil, err
}

func convertToRFC3339(date string) (string, error) {
	layout := "2006-01-02"

	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	parsedDateString := parsedDate.Format(time.RFC3339)
	return parsedDateString, nil
}

func ConstructDateQuery(DateFrom string, DateTo string) (string, error) {
	// construct date query
	fromDate := "*"
	toDate := "*"
	log.Printf("building date from: %s\n", DateFrom)
	log.Printf("building date to: %s\n", DateTo)

	if DateFrom != "" || DateTo != "" {
		var err error
		if DateFrom != "" {
			fromDate, err = convertToRFC3339(DateFrom)
			if err != nil {
				return "date_filed:[* TO *]", fmt.Errorf("invalid date format for DateFrom: %v", err)
			}
			DateFrom = ""
		}
		if DateTo != "" {
			toDate, err = convertToRFC3339(DateTo)
			if err != nil {
				return "date_filed:[* TO *]", fmt.Errorf("invalid date format for DateTo: %v", err)
			}
			DateTo = ""
		}
	}
	dateQuery := fmt.Sprintf("date_filed:[%s TO %s]", fromDate, toDate)
	return dateQuery, nil
}

// THESE ARE THE IMPORTANT MAPPINGS
var QuickwitFilterMapping = map[string]string{
	"DateFrom": "date_filed",
}

func constructQuickwitMetadataQueryString(meta networking.Metadata) string {
	values := reflect.ValueOf(meta)
	types := reflect.TypeOf(meta)
	return constructGenericFilterQuery(values, types)
}

func constructQuickwitUUIDMetadataQueryString(meta networking.UUIDFilterFields) string {
	values := reflect.ValueOf(meta)
	types := reflect.TypeOf(meta)
	return constructGenericFilterQuery(values, types)
}

func constructGenericFilterQuery(values reflect.Value, types reflect.Type) string {
	var filterQuery string
	filters := []string{}

	fmt.Printf("values: %v\n", values)
	fmt.Printf("types: %v\n", types)

	// ===== iterate over metadata for filter =====
	for i := 0; i < types.NumField(); i++ {
		// get the field and value
		field := types.Field(i)
		value := values.Field(i)
		tag := field.Tag.Get("json")
		if strings.Contains(tag, ",omitempty") {
			tag = strings.Split(tag, ",")[0]
		}

		fmt.Printf("tag: %v\nfield: %v\nvalue: %v\n", tag, field, value)

		if tag == "fileuuid" {
			tag = "source_id"
		}
		s := fmt.Sprintf("metadata.%s:(%s)", tag, value)

		// exlude empty values
		if strings.Contains(s, "00000000-0000-0000-0000-000000000000") {
			continue
		}
		log.Printf("new filter: %s\n", s)
		filters = append(filters, s)
	}
	// concat all filters with AND clauses
	for _, f := range filters {
		filterQuery += fmt.Sprintf(" AND (%s)", f)
	}
	fmt.Printf("filter query: %s\n", filterQuery)
	return filterQuery
}

func SearchMilvus(request SearchRequest) ([]SearchData, error) {
	return []SearchData{}, fmt.Errorf("not implemented")
}

func FormatSearchResults(searchResults []SearchData, query string) string {
	searchResultsString := fmt.Sprintf("Query: %s\n", query)
	for _, result := range searchResults {
		searchResultsString += fmt.Sprintf("Name: %s\n", result.Name)
		// searchResultsString += fmt.Sprintf("Text: %s\n", result.Text)
		searchResultsString += fmt.Sprintf("DocID: %s\n", result.DocID)
		searchResultsString += fmt.Sprintf("SourceID: %s\n", result.SourceID)
	}
	return searchResultsString
}
