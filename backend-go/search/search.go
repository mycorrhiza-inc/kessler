package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

var quickwitURL = os.Getenv("QUICKWIT_ENDPOINT")

type Hit struct {
	CreatedAt     string   `json:"created_at"`
	Doctype       string   `json:"doctype"`
	Hash          string   `json:"hash"`
	Lang          string   `json:"lang"`
	DateFiled     string   `json:"updated_at"`
	Metadata      Metadata `json:"metadata"`
	Name          string   `json:"name"`
	SaOrmSentinel *string  `json:"sa_orm_sentinel"`
	ShortSummary  *string  `json:"short_summary"`
	Source        string   `json:"source"`
	SourceID      string   `json:"source_id"`
	Stage         string   `json:"stage"`
	Summary       *string  `json:"summary"`
	Text          string   `json:"text"`
	Timestamp     string   `json:"timestamp"`
	URL           string   `json:"url"`
}

type Metadata struct {
	Author   string `json:"author"`
	Date     string `json:"date"`
	DocketID string `json:"docket_id"`
	Doctype  string `json:"doctype"`
	Lang     string `json:"lang"`
	Language string `json:"language"`
	Source   string `json:"source"`
	Title    string `json:"title"`
}

type FilterFields struct {
	Metadata
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (m Metadata) String() string {
	// Marshal the struct to JSON format
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
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
	SortBy        string `json:"sort_by"`
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
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	parsedDateString := parsedDate.Format(time.RFC3339)
	return parsedDateString, nil
}

func SearchQuickwit(r SearchRequest) ([]SearchData, error) {
	// TODO: Cleanup default value if not gotten from request, just for bugfixing - Nicole
	if r.MaxHits == 0 {
		r.MaxHits = 100
	}
	search_index := r.Index
	// ===== construct search request =====
	query := r.Query
	log.Printf("search filters: %s\n", r.SearchFilters)
	var queryString string
	// construct date query
	fromDate := "*"
	toDate := "*"

	if r.SearchFilters.DateFrom != "" || r.SearchFilters.DateTo != "" {
		var err error
		if r.SearchFilters.DateFrom != "" {
			fromDate, err = convertToRFC3339(r.SearchFilters.DateFrom)
			if err != nil {
				return nil, fmt.Errorf("invalid date format for DateFrom: %v", err)
			}
			r.SearchFilters.DateFrom = ""
		}
		if r.SearchFilters.DateTo != "" {
			toDate, err = convertToRFC3339(r.SearchFilters.DateTo)
			if err != nil {
				return nil, fmt.Errorf("invalid date format for DateTo: %v", err)
			}
			r.SearchFilters.DateTo = ""
		}
	}
	dateQuery := fmt.Sprintf("date_filed:[%s TO %s]", fromDate, toDate)

	if len(r.Query) >= 0 {
		queryString = fmt.Sprintf("((text:(%s) OR name:(%s)) AND %s)", query, query, dateQuery)
	}

	filtersString := constructQuickwitMetadataQueryString(r.SearchFilters)

	queryString = queryString + filtersString
	log.Printf("full query string: %s\n", queryString)

	// construct sortby string
	sortbyStr := ""
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

	data, err := ExtractSearchData(searchResponse)
	if err != nil {
		log.Printf("Error creating response data: %s", err)
		errturn(err)
	}

	return data, nil
}

func constructQuickwitMetadataQueryString(filter FilterFields) string {
	var filterQuery string
	filters := []string{}
	fmt.Printf("constructing filter of:\n%s\n", filter)

	// ===== reflect the filter metadata =====
	values := reflect.ValueOf(filter)
	types := reflect.TypeOf(filter)

	// ===== iterate over metadata for filter =====
	for i := 0; i < types.NumField(); i++ {
		// get the field and value
		field := types.Field(i)
		value := values.Field(i)
		// Get the json tag value
		tag := field.Tag.Get("json")

		// skip all non-slice filters and empty slices
		if value.Kind() != reflect.String || value.Len() <= 0 {
			continue
		}

		// format each query equality
		s := fmt.Sprintf("metadata.%s:(%s)", tag, value)
		filters = append(filters, s)

		// TODO: allow for multiple distinct filters per filter segment
		// construct the filter specific string
		// filterString := field_queries[0]
		// for q := 1; q < len(field_queries); q++ {
		// 	filterString += fmt.Sprintf(" OR %s", field_queries[q])
		// }
	}
	// concat all filters with AND clauses
	for _, f := range filters {
		// WARN: This is potentially not safe. TBD if quickwit's query language is
		// vulnerable to injectable attacks
		// fmt.Printf("got filter: \n%s\n", f)
		filterQuery += fmt.Sprintf(" AND (%s)", f)
	}
	return filterQuery
}

func SearchMilvus(request SearchRequest) ([]SearchData, error) {
	return []SearchData{}, fmt.Errorf("not implemented")
}

type SearchReturn struct {
	Results []SearchData
	Error   error
}

func HybridSearch(request SearchRequest) ([]SearchData, error) {
	// This could absolutely be dryified
	chanMilvus := make(chan SearchReturn)
	chanQuickwit := make(chan SearchReturn)
	go func() {
		results, err := SearchMilvus(request)
		return_results := SearchReturn{results, err}
		chanMilvus <- return_results
	}()
	go func() {
		results, err := SearchQuickwit(request)
		return_results := SearchReturn{results, err}
		chanQuickwit <- return_results
	}()
	resultsMilvus := <-chanMilvus
	resultsQuickwit := <-chanQuickwit
	if resultsMilvus.Error == nil {
		fmt.Printf("Milvus returned error: %s", resultsMilvus.Error)
	}
	if resultsQuickwit.Error == nil {
		fmt.Printf("Quickwit returned error: %s", resultsQuickwit.Error)
	}
	if resultsMilvus.Error != nil && resultsQuickwit.Error != nil {
		return []SearchData{}, fmt.Errorf("both Milvus and Quickwit returned errors. milvus error: %s quickwit error: %s", resultsMilvus.Error, resultsQuickwit.Error)
	}
	unrankedResults := append(resultsMilvus.Results, resultsQuickwit.Results...)
	rerankedData, err := rerankSearchResults(unrankedResults, request.Query)
	// Fail semi silently and returns the regular unranked results
	if err != nil {
		log.Printf("Error reranking results: %s", err)
		return unrankedResults, nil
	}

	return rerankedData, nil
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
