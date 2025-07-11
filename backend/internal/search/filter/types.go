// internal/search/filter/types.go
package filter

import (
	"kessler/internal/fugusdk"
)

// Filter represents a filter path with its available values
type Filter struct {
	FilterPath string   `json:"filter_path"`
	Values     []string `json:"values"`
}

// FilterResponse represents the response containing all filters
type FilterResponse struct {
	Status  string           `json:"status"`
	Filters []fugusdk.Filter `json:"filters"`
}

// NamespaceFilterResponse represents filters for a specific namespace
type NamespaceFilterResponse struct {
	Status    string              `json:"status"`
	Namespace string              `json:"namespace"`
	Filters   map[string][]string `json:"filters"`
}

// FilterValuesResponse represents values for a specific filter path
type FilterValuesResponse struct {
	Status     string   `json:"status"`
	FilterPath string   `json:"filter_path"`
	Values     []string `json:"values"`
}
