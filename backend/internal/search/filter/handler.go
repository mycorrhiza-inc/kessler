package filter

import (
	"encoding/json"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var handlerTracer = otel.Tracer("filter-handler")

// Handler handles HTTP requests for filter configuration
type Handler struct {
	service *Service
}

// NewHandler creates a new filter configuration handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetConfiguration handles requests for filter configuration
func (h *Handler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:get-configuration")
	defer span.End()

	logger.Info(ctx, "filter configuration requested")

	config, err := h.service.BuildFilterConfiguration(ctx)
	if err != nil {
		logger.Error(ctx, "failed to build filter configuration", zap.Error(err))
		http.Error(w, "Failed to load filter configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		logger.Error(ctx, "failed to encode configuration response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "filter configuration served", zap.Int("field_count", len(config.Fields)))
}

// ConvertFilters handles requests to convert frontend filters to backend format
func (h *Handler) ConvertFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:convert-filters")
	defer span.End()

	logger.Info(ctx, "filter conversion requested")

	var req ConvertFiltersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode convert filters request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	backendFilters, err := h.service.ConvertFiltersToBackend(ctx, req.Filters)
	if err != nil {
		logger.Error(ctx, "failed to convert filters", zap.Error(err))
		http.Error(w, "Failed to convert filters", http.StatusInternalServerError)
		return
	}

	response := ConvertFiltersResponse{
		BackendFilters: backendFilters,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode convert filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "filters converted", zap.Int("filter_count", len(req.Filters)))
}

// ValidateFilters handles requests to validate filter values
func (h *Handler) ValidateFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:validate-filters")
	defer span.End()

	logger.Info(ctx, "filter validation requested")

	var req ValidateFiltersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode validate filters request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validation, err := h.service.ValidateFilters(ctx, req.Filters)
	if err != nil {
		logger.Error(ctx, "failed to validate filters", zap.Error(err))
		http.Error(w, "Failed to validate filters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(validation); err != nil {
		logger.Error(ctx, "failed to encode validate filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "filters validated",
		zap.Bool("is_valid", validation.IsValid),
		zap.Int("error_count", len(validation.Errors)))
}

// GetOptions handles requests for dynamic filter options
func (h *Handler) GetOptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := handlerTracer.Start(ctx, "filter-handler:get-options")
	defer span.End()

	logger.Info(ctx, "filter options requested")

	var req GetOptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode get options request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	options, err := h.service.GetDynamicOptions(ctx, req.FieldID, req.Context, req.Namespace)
	if err != nil {
		logger.Error(ctx, "failed to get dynamic options", zap.Error(err))
		http.Error(w, "Failed to get options", http.StatusInternalServerError)
		return
	}

	response := GetOptionsResponse{
		Options: options,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode get options response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "filter options served",
		zap.String("field_id", req.FieldID),
		zap.Int("option_count", len(options)))
}

// GetNamespaceFilters retrieves filters available for a specific namespace
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

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get namespace filters from fugu
	response, err := client.GetNamespaceFilters(ctx, namespace)
	if err != nil {
		logger.Error(ctx, "failed to get namespace filters", zap.Error(err))
		http.Error(w, "Failed to get namespace filters", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode namespace filters response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "namespace filters served successfully", zap.String("namespace", namespace))
}
