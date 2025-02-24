package filters

import (
	"encoding/json"
	"kessler/common/objects/networking"
	"kessler/gen/dbstore"
	"net/http"

	"github.com/charmbracelet/log"
)

type FilterHandler struct {
	service *FilterService
}

func NewFilterHandler(service *FilterService) *FilterHandler {
	return &FilterHandler{
		service: service,
	}
}

func (h *FilterHandler) GetFilters(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	pagination := networking.PaginationFromUrlParams(r)

	filters, err := h.service.GetFiltersByState(r.Context(), state)
	if err != nil {
		switch err {
		case ErrInvalidFilterState:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case ErrDatabaseOperation:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Error("failed to get filters", "error", err)
		return
	}

	// Apply pagination
	start := int(pagination.Offset)
	end := start + int(pagination.Limit)
	if end > len(filters) {
		end = len(filters)
	}
	if start > len(filters) {
		filters = []dbstore.Filter{}
	} else {
		filters = filters[start:end]
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filters); err != nil {
		log.Error("failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
