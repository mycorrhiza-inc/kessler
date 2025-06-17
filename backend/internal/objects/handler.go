package objects

import (
	"encoding/json"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("objects-endpoint")

// ObjectResponse represents a response containing a single object
type ObjectResponse struct {
	Data      interface{} `json:"data,omitempty"`      // The object data
	ID        string      `json:"id,omitempty"`        // The object ID
	Text      string      `json:"text,omitempty"`      // The object text content
	Metadata  interface{} `json:"metadata,omitempty"`  // The object metadata
	Namespace string      `json:"namespace,omitempty"` // The object namespace
}

// ObjectService handles business logic for objects
type ObjectService struct {
	fuguServerURL string
}

// NewObjectService creates a new object service
func NewObjectService(fuguServerURL string) *ObjectService {
	return &ObjectService{
		fuguServerURL: fuguServerURL,
	}
}

// ObjectHandler handles HTTP requests for objects
type ObjectHandler struct {
	service *ObjectService
}

// NewObjectHandler creates a new object handler
func NewObjectHandler(service *ObjectService) *ObjectHandler {
	return &ObjectHandler{
		service: service,
	}
}

// RegisterObjectRoutes registers object routes with the router
func RegisterObjectRoutes(r *mux.Router) error {
	fuguServerURL := "http://fugudb:3301"
	service := NewObjectService(fuguServerURL)
	handler := NewObjectHandler(service)

	// Create objects subrouter
	objectsRoute := r.PathPrefix("").Subrouter()

	// Get object by ID endpoint
	objectsRoute.HandleFunc("/{id}", handler.GetObjectByID).Methods(http.MethodGet)

	return nil
}

// @Summary     Get Object by ID
// @Description Retrieves a specific object by its ID
// @Tags        objects
// @Produce     json
// @Param       id path string true "Object ID"
// @Success     200 {object} ObjectResponse
// @Failure     400 {string} string "Invalid object ID"
// @Failure     404 {string} string "Object not found"
// @Failure     500 {string} string "Internal server error"
// @Router      /objects/{id} [get]
func (h *ObjectHandler) GetObjectByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "objects:GetObjectByID")
	defer span.End()

	// Get object ID from URL
	vars := mux.Vars(r)
	objectID := vars["id"]
	if objectID == "" {
		logger.Error(ctx, "object ID not provided")
		http.Error(w, "Object ID is required", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "getting object by ID", zap.String("object_id", objectID))

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get object from fugu
	response, err := client.GetObjectByID(ctx, objectID)
	if err != nil {
		logger.Error(ctx, "failed to get object from fugu", zap.Error(err))
		http.Error(w, "Object not found or server error", http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "successfully retrieved object", zap.String("object_id", objectID))
}
