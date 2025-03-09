package filters

import (
	"encoding/json"
	"kessler/cache"
	"kessler/common/objects/networking"
	"kessler/database"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type FilterServiceHandler struct {
	service *FilterService
}

func NewFilterHandler(service *FilterService) *FilterServiceHandler {
	return &FilterServiceHandler{
		service: service,
	}
}

func RegisterFilterRoutes(r *mux.Router) error {
	service := NewFilterService(database.ConnPool, cache.MemcachedClient)
	fsh := &FilterServiceHandler{
		service: service,
	}
	filtersRoute := r.PathPrefix("/filters").Subrouter()
	filtersRoute.HandleFunc(
		"",
		fsh.GetFilters,
	).Methods(http.MethodGet)
	return nil
}

func (h *FilterServiceHandler) GetFilters(w http.ResponseWriter, r *http.Request) {
	log.Info("get filters called")
	state := r.URL.Query().Get("state")
	pagination := networking.PaginationFromUrlParams(r)

	f, err := h.service.GetFiltersByState(r.Context(), state)

	if err != nil {
		log.Error("There was an error listing the filters", zap.Error(err))
		switch err {
		case ErrInvalidFilterState:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case ErrDatabaseOperation:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}

	filters := make([]Filter, len(f))
	for i, fn := range f {
		filters[i] = Filter{fn}
	}

	// Apply pagination
	start := int(pagination.Offset)
	end := start + int(pagination.Limit)
	if end > len(filters) {
		end = len(filters)
	}
	if start > len(filters) {
		filters = []Filter{}
	} else {
		filters = filters[start:end]
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filters); err != nil {
		log.Error("failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
