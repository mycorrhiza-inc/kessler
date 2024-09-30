package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	TopN      int      `json:"top_n"`
	Documents []string `json:"documents"`
}

type RerankResponse struct {
	Results []struct {
		Index          int     `json:"index"`
		RelevanceScore float64 `json:"relevance_score"`
	} `json:"results"`
}

var CO_API_KEY = os.Getenv("CO_API_KEY")

func rerankStringsAndQueryPermutation(ctx context.Context, query string, documents []string) ([]int, error) {
	const url = "https://api.cohere.com/v1/rerank"

	reqBody := RerankRequest{
		Model:     "rerank-english-v3.0",
		Query:     query,
		TopN:      len(documents),
		Documents: documents,
	}

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", CO_API_KEY))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rerankResp RerankResponse
	err = json.Unmarshal(body, &rerankResp)
	if err != nil {
		return nil, err
	}

	permutation := make([]int, len(documents))
	for i, result := range rerankResp.Results {
		permutation[i] = result.Index
	}
	return permutation, nil
}

func rerankSearchResults(searchResults []SearchData, query string) ([]SearchData, error) {
	var documents []string
	for _, result := range searchResults {
		documents = append(documents, result.Text)
	}
	permutation, err := rerankStringsAndQueryPermutation(context.Background(), query, documents)
	if err != nil {
		return nil, err
	}
	rerankedResults := make([]SearchData, len(searchResults))
	for i, permutation := range permutation {
		rerankedResults[i] = searchResults[permutation]
	}
	return rerankedResults, nil
}
