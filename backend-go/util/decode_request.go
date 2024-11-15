package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeRequest(body io.ReadCloser, target interface{}, w http.ResponseWriter) {
	err := json.NewDecoder(body).Decode(target)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
	}
}
