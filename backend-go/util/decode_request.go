package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func DecodeRequest(body io.ReadCloser, target interface{}, w http.ResponseWriter) {
	err := json.NewDecoder(body).Decode(target)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
	}
}

func ParseStringSliceUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(ids))
	for i, v := range ids {
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			return []uuid.UUID{}, err
		}
		uuids[i] = parsedUUID
	}
	return uuids, nil
}
