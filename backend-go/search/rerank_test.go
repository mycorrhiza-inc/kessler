package search

import (
	"context"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

func Test_reranker(t *testing.T) {
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

	rerankedDocPermutation, err := rerankStringsAndQueryPermutation(ctx, query, documents)
	if err != nil {
		t.Fatal("Error:", err)
	}
	log.Info("Permutation:")
	for _, docperm := range rerankedDocPermutation {
		log.Info(docperm)
	}
	rerankedDocs := make([]string, len(rerankedDocPermutation))
	for i, permutation := range rerankedDocPermutation {
		rerankedDocs[i] = documents[permutation]
	}

	log.Info("Query:", query)
	log.Info("Reranked Documents:")
	for _, doc := range rerankedDocs {
		log.Info(doc)
	}
}
