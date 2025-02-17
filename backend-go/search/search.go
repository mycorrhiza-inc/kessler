package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kessler/objects/authors"
	"kessler/objects/conversations"
	"kessler/objects/files"
	"kessler/objects/networking"
	"kessler/objects/timestamp"
	"kessler/quickwit"
	"reflect"

	"github.com/charmbracelet/log"

	"github.com/google/uuid"
)

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
		log.Info("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
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
		log.Info("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
}

func SearchQuickwit(r SearchRequest) ([]SearchDataHydrated, error) {
	log.Info(fmt.Sprintf("searching quickwit:\n%s", r))
	r.Index = "NY_PUC"
	search_index := r.Index
	// ===== construct search request =====
	query := r.Query

	var queryFilters networking.FilterFields = r.SearchFilters
	var metadataFilters networking.MetadataFilterFields = queryFilters.MetadataFilters
	var uuidFilters networking.UUIDFilterFields = queryFilters.UUIDFilters
	log.Info(fmt.Sprintf("zzxxcc: %v\n", uuidFilters))
	dateQueryString := quickwit.ConstructDateTextQuery(metadataFilters.DateFrom, metadataFilters.DateTo, query)

	filtersString := constructQuickwitMetadataQueryString(metadataFilters.SearchMetadata)
	uuidFilterString := constructQuickwitUUIDMetadataQueryString(uuidFilters)

	log.Info(fmt.Sprintf()
		"!!!!!!!!!!\nquery: %s\nfilters: %s\nuuid filters: %s\n!!!!!!!!!!\n",
		dateQueryString,
		filtersString,
		uuidFilterString,
	)
	// queryString = queryString + filtersString
	queryString := dateQueryString + filtersString + uuidFilterString
	log.Info(fmt.Sprintf("full query string: %s\n", queryString))

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
	request := quickwit.QuickwitSearchRequest{
		Query:         queryString,
		MaxHits:       r.MaxHits,
		SnippetFields: snippetFields,
		StartOffset:   r.StartOffset,
		SortBy:        sortbyStr,
	}
	return_bytes, err := quickwit.PerformGenericQuickwitRequest(request, search_index)
	if err != nil {
		log.Info(fmt.Sprintf("Error with Quickwit Request: %v", err))
		return []SearchDataHydrated{}, err
	}
	var searchResponse quickwitSearchResponse
	err = json.NewDecoder(bytes.NewReader(return_bytes)).Decode(&searchResponse)
	if err != nil {
		log.Info(fmt.Sprintf("quickwit response: %v\n", return_bytes))
		errorstring := fmt.Sprintf("Error decoding JSON: %v\n Offending json looked like: %v", err, string(return_bytes))
		return []SearchDataHydrated{}, fmt.Errorf(errorstring)
	}

	data, err := ExtractSearchData(searchResponse)
	if err != nil {
		log.Info(fmt.Sprintf("Error creating response data: %s", err))
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
		file_timestamp, _ := timestamp.KessTimeFromString(hit.DateFiled)
		file_schema := files.CompleteFileSchema{
			ID:           file_id,
			Name:         hit.Name,
			Extension:    hit.Metadata.Extension,
			Conversation: convo_info,
			Mdata: files.FileMetadataSchema{
				"docket_id":   hit.Metadata.DocketID,
				"date":        hit.Metadata.Date,
				"file_class":  hit.Metadata.FileClass,
				"item_number": hit.Metadata.ItemNumber,
			},
			IsPrivate:     false,
			DatePublished: file_timestamp,
			Authors:       author_infos,
			DocTexts:      []files.FileChildTextSource{},
			Stage:         files.DocProcStage{},
			Extra:         files.FileGeneratedExtras{},
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

// THESE ARE THE IMPORTANT MAPPINGS
var QuickwitFilterMapping = map[string]string{
	"DateFrom": "date_filed",
}

func constructQuickwitMetadataQueryString(meta networking.SearchMetadata) string {
	values := reflect.ValueOf(meta)
	types := reflect.TypeOf(meta)
	return quickwit.ConstructGenericFilterQuery(values, types, true)
}

func constructQuickwitUUIDMetadataQueryString(meta networking.UUIDFilterFields) string {
	values := reflect.ValueOf(meta)
	types := reflect.TypeOf(meta)
	return quickwit.ConstructGenericFilterQuery(values, types, false)
}

func SearchMilvus(request SearchRequest) ([]SearchData, error) {
	return []SearchData{}, fmt.Errorf("not implemented")
}

func FormatSearchResults(searchResults []SearchDataHydrated, query string) string {
	searchResultsString := fmt.Sprintf("Query: %s\n", query)
	for _, result := range searchResults {
		searchResultsString += fmt.Sprintf("Name: %s\n", result.Name)
		// searchResultsString += fmt.Sprintf("Text: %s\n", result.Text)
		searchResultsString += fmt.Sprintf("DocID: %s\n", result.DocID)
		searchResultsString += fmt.Sprintf("SourceID: %s\n", result.SourceID)
	}
	return searchResultsString
}
