package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
)

var quickwitURL = "http://100.86.5.114:7280"

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
	Name  string `json:"name"`
	Text  string `json:"text"`
	DocID string `json:"docID"`
}

// Function to create search data array
func ExtractSearchData(data quickwitSearchResponse) ([]SearchData, error) {
	var result []SearchData

	// Map snippets text to hit names
	for i, hit := range data.Hits {
		sdata := SearchData{
			Name:  hit.Name,
			Text:  data.Snippets[i].Text[0],
			DocID: hit.Metadata.DocketID, // Assuming title from metadata is used as docID
		}
		result = append(result, sdata)
	}

	return result, nil
}

func errturn(err error) ([]SearchData, error) {
	return nil, err
}

func searchQuickwit(query string) ([]SearchData, error) {

	request := createQWRequest(query)
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error Marshalling quickwit request: %s", err)
		errturn(err)
	}

	request_url := fmt.Sprintf("%s/api/v1/dockets/search", quickwitURL)
	resp, err := http.Post(
		request_url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Error sending request to quickwit: %s", err)
		errturn(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %v", resp.StatusCode)
		return errturn(fmt.Errorf("Error: received status code %v", resp.StatusCode))
	}
	var searchResponse quickwitSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling quickwit response: %s", err)
		errturn(err)
	}

	data, err := ExtractSearchData(searchResponse)

	if err != nil {
		log.Fatalf("Error creating response data: %s", err)
		errturn(err)
	}

	return data, nil
}
