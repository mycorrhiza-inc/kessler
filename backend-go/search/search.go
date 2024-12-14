package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kessler/objects/networking"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
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

func SearchQuickwit(r SearchRequest) ([]SearchData, error) {
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
		return nil, err
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
		a, e := errturn(fmt.Errorf("received status code %v", resp.StatusCode))
		return a, e
	}
	var searchResponse quickwitSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		log.Printf("Error unmarshalling quickwit response: %s", err)
		errturn(err)
	}
	validated_response := ValidateSearchRequest(r, searchResponse)

	data, err := ExtractSearchData(validated_response)
	if err != nil {
		log.Printf("Error creating response data: %s", err)
		errturn(err)
	}

	return data, nil
}

// write a function that will take in a searchRequest and searchResults and create a new quickwitSearchResponse then for each hit and snippet in the passed in search results, make sure all the filters in search request are valid for that, if it is valid append it to the return searchResponse, else skip it and print a scary error message, then return the list of validated results
func ValidateSearchRequest(searchRequest SearchRequest, searchResults quickwitSearchResponse) quickwitSearchResponse {
	global_filters := searchRequest.SearchFilters
	filters := global_filters.MetadataFilters
	metadata_filters := filters.Metadata
	var validatedResponse quickwitSearchResponse

	for i, hit := range searchResults.Hits {
		isValid := true

		// Validate query matches if present
		if searchRequest.Query != "" {
			if !strings.Contains(strings.ToLower(hit.Text), strings.ToLower(searchRequest.Query)) &&
				!strings.Contains(strings.ToLower(hit.Name), strings.ToLower(searchRequest.Query)) {
				log.Printf("❌ WARNING: Hit %d failed query validation", i)
				isValid = false
			}
		}

		// Claude writes some interesting go code lol - nic
		// Validate metadata filters
		v := reflect.ValueOf(metadata_filters)
		t := reflect.TypeOf(metadata_filters)
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			value := v.Field(j)

			// Skip empty string fields
			if value.Kind() == reflect.String && value.String() != "" {
				hitValue := reflect.ValueOf(hit.Metadata).FieldByName(field.Name)
				if !hitValue.IsValid() || hitValue.String() != value.String() {
					log.Printf("❌ WARNING: Hit %d failed metadata validation for field %s", i, field.Name)
					isValid = false
					break
				}
			}
		}

		// If all validations pass, append to response
		if isValid {
			validatedResponse.Hits = append(validatedResponse.Hits, hit)
			if i < len(searchResults.Snippets) {
				validatedResponse.Snippets = append(validatedResponse.Snippets, searchResults.Snippets[i])
			}
		}
	}

	return validatedResponse
}

// Function to create search data array
func ExtractSearchData(data quickwitSearchResponse) ([]SearchData, error) {
	var result []SearchData

	// Map snippets text to hit names
	for i, hit := range data.Hits {
		var snippet string
		if len(data.Snippets) > 0 {
			snippet = data.Snippets[i].Text[0]
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
