package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
)

var quickwitURL = os.Getenv("QUICKWIT_ENDPOINT")

type Hit struct {
	CreatedAt     string   `json:"created_at"`
	Doctype       string   `json:"doctype"`
	Hash          string   `json:"hash"`
	Lang          string   `json:"lang"`
	Metadata      Metadata `json:"metadata"`
	Name          string   `json:"name"`
	OriginalText  *string  `json:"original_text"`
	SaOrmSentinel *string  `json:"sa_orm_sentinel"`
	ShortSummary  *string  `json:"short_summary"`
	Source        string   `json:"source"`
	SourceID      string   `json:"source_id"`
	Stage         string   `json:"stage"`
	Summary       *string  `json:"summary"`
	Text          string   `json:"text"`
	Timestamp     string   `json:"timestamp"`
	UpdatedAt     string   `json:"updated_at"`
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

type Snippet struct {
	Text []string `json:"text"`
}

type quickwitSearchResponse struct {
	Hits     []Hit     `json:"hits"`
	Snippets []Snippet `json:"snippets"`
}

type QuickwitSearchRequest struct {
	Query         string `json:"query"`
	SnippetFields string `json:"snippet_fields"`
	MaxHits       int    `json:"max_hits"`
}

func createQWRequest(query string) QuickwitSearchRequest {
	queryString := fmt.Sprintf("text:(%s) OR name:(%s)", query, query)

	log.Printf("Query String: %s\n", queryString)

	return QuickwitSearchRequest{
		Query:         queryString,
		SnippetFields: "text",
		MaxHits:       20,
	}
}

type SearchData struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	DocID    string `json:"docID"`
	SourceID string `json:"sourceID"`
}

// Function to create search data array
func ExtractSearchData(data quickwitSearchResponse) ([]SearchData, error) {
	var result []SearchData

	// Map snippets text to hit names
	for i, hit := range data.Hits {
		sdata := SearchData{
			Name:     hit.Name,
			Text:     data.Snippets[i].Text[0],
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

func constructQuickwitMetadataQueryString(filter Metadata) string {
	var filterQuery string
	filters := []string{}

	values := reflect.ValueOf(filter)
	types := reflect.TypeOf(filter)

	for i := 0; i < types.NumField(); i++ {
		// get the field and value
		field := types.Field(i)
		value := values.Field(i)
		// Get the json tag value
		tag := field.Tag.Get("json")

		// skip all non-slice filters and empty slices
		if value.Kind() != reflect.Slice || value.Len() <= 0 {
			continue
		}

		// TODO: handle excluding values from the query

		field_queries := []string{}
		// format each query equality
		for j := 0; j < value.Len(); j++ {
			s := fmt.Sprintf("metadata.%s:(%s)", tag, value.Index(j))
			field_queries = append(field_queries, s)
		}

		// construct the filter specific string
		filterString := field_queries[0]
		for q := 1; q < len(field_queries); q++ {
			filterString += fmt.Sprintf(" OR %s", field_queries[q])
		}
	}
	// concat all filters with AND clauses
	for _, f := range filters {
		// WARN: This is potentially not safe. TBD if quickwit's query language is
		// vulnerable to injectable attacks
		filterQuery += fmt.Sprintf(" AND (%s)", f)
	}
	return filterQuery
}

func SearchQuickwit(r SearchRequest) ([]SearchData, error) {
	if len(r.Query) <= 0 {
		return []SearchData{}, nil
	}

	queryString := r.Query
	filtersString := constructQuickwitMetadataQueryString(r.SearchFilters)

	queryString += filtersString

	request := createQWRequest(queryString)
	jsonData, err := json.Marshal(request)
	log.Printf("jsondata: \n%s", jsonData)
	if err != nil {
		log.Printf("Error Marshalling quickwit request: %s", err)
		return nil, err
	}

	request_url := fmt.Sprintf("%s/api/v1/dockets/search", quickwitURL)
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

func SearchMilvus(request SearchRequest) ([]SearchData, error) {
	return []SearchData{}, fmt.Errorf("Not implemented")
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
		searchResultsString += fmt.Sprintf("Text: %s\n", result.Text)
		// searchResultsString += fmt.Sprintf("DocID: %s\n", result.DocID)
	}
	return searchResultsString
}
