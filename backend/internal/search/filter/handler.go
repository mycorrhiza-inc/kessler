// handler.go
package filter

import (
	"encoding/json"
	"kessler/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var handlerTracer = otel.Tracer("filter-handler")

// Handler handles HTTP requests for filters
type Handler struct {
	service *Service
}

// NewHandler creates a new filter handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetAllFilters handles requests for all available filters
func (h *Handler) GetAllFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:get-all-filters")
	defer span.End()

	logger.Info(ctx, "all filters requested")

	filters, err := h.service.GetAllFilters(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get all filters", zap.Error(err))
		http.Error(w, "Failed to get filters", http.StatusInternalServerError)
		return
	}

	response := FilterResponse{
		Status:  "success",
		Filters: filters,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode all filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "all filters served successfully", zap.Int("filter_count", len(filters)))
}

// GetNamespaceFilters handles requests for namespace-specific filters
func (h *Handler) GetNamespaceFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:get-namespace-filters")
	defer span.End()

	// Get namespace from URL
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	if namespace == "" {
		logger.Error(ctx, "namespace not provided")
		http.Error(w, "Namespace is required", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "namespace filters requested", zap.String("namespace", namespace))

	filters, err := h.service.GetNamespaceFilters(ctx, namespace)
	if err != nil {
		logger.Error(ctx, "failed to get namespace filters", zap.Error(err))
		http.Error(w, "Failed to get namespace filters", http.StatusInternalServerError)
		return
	}

	response := NamespaceFilterResponse{
		Status:    "success",
		Namespace: namespace,
		Filters:   filters,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode namespace filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "namespace filters served successfully",
		zap.String("namespace", namespace),
		zap.Int("filter_count", len(filters)))
}

// GetFilterValues handles requests for values of a specific filter path
func (h *Handler) GetFilterValues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:get-filter-values")
	defer span.End()

	// Get filter path from URL
	vars := mux.Vars(r)
	filterPath := vars["filterPath"]
	if filterPath == "" {
		logger.Error(ctx, "filter path not provided")
		http.Error(w, "Filter path is required", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "filter values requested", zap.String("filter_path", filterPath))

	values, err := h.service.GetFilterValues(ctx, filterPath)
	if err != nil {
		logger.Error(ctx, "failed to get filter values", zap.Error(err))
		http.Error(w, "Failed to get filter values", http.StatusInternalServerError)
		return
	}

	response := FilterValuesResponse{
		Status:     "success",
		FilterPath: filterPath,
		Values:     values,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode filter values response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "filter values served successfully",
		zap.String("filter_path", filterPath),
		zap.Int("value_count", len(values)))
}

// InvalidateCache handles requests to clear filter caches
func (h *Handler) InvalidateCache(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:invalidate-cache")
	defer span.End()

	logger.Info(ctx, "cache invalidation requested")

	h.service.InvalidateCache(ctx)

	response := map[string]string{
		"status":  "success",
		"message": "Filter cache invalidated",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode cache invalidation response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "cache invalidation completed")
}
