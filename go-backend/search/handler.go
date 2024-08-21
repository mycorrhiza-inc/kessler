package search

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
}

func HandleSearchRequest(w http.ResponseWriter, r *http.Request) {
	// Create an instance of RequestData
	var RequestData SearchRequest

	// Decode the JSON body into the struct
	err := json.NewDecoder(r.Body).Decode(&RequestData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body when done

	data, err := searchQuickwit(string(RequestData.Query))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respString, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(respString))
}
