package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

func rerankDocuments(ctx context.Context, query string, documents []string) ([]string, error) {
	const url = "https://api.cohere.com/v1/rerank"
	apiKey := "your_api_key_here" // Replace with your actual API key

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
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", apiKey))

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

	rerankedDocuments := make([]string, len(documents))
	for _, result := range rerankResp.Results {
		rerankedDocuments[result.Index] = documents[result.Index]
	}
	return rerankedDocuments, nil
}

func test_reranker() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "What is the capital of the United States?"
	documents := []string{
		"Carson City is the capital city of the American state of Nevada.",
		"The Commonwealth of the Northern Mariana Islands is a group of islands in the Pacific Ocean. Its capital is Saipan.",
		"Washington, D.C. (also known as simply Washington or D.C., and officially as the District of Columbia) is the capital of the United States. It is a federal district.",
		"Capitalization or capitalisation in English grammar is the use of a capital letter at the start of a word. English usage varies from capitalization in other languages.",
		"Capital punishment (the death penalty) has existed in the United States since before the United States was a country. As of 2017, capital punishment is legal in 30 of the 50 states.",
	}

	rerankedDocs, err := rerankDocuments(ctx, query, documents)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Reranked Documents:")
	for _, doc := range rerankedDocs {
		fmt.Println(doc)
	}
}
