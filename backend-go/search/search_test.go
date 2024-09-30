package search

import (
	"fmt"
	"testing"
)

var example_search_request = SearchRequest{
	"Marshall Fire",
	Metadata{},
}

func TestQuickwit(t *testing.T) {
	results, err := searchQuickwit(example_search_request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	resultsString := fmt.Sprintf("%v", results)
	fmt.Println("Results:", resultsString)
}
