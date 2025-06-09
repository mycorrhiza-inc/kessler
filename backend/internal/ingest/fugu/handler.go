package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("fugu-migration-endpoint")

// IngestRequest represents the incoming batch ingest request
type IngestRequest struct {
	Records []Record `json:"records"`
}

// Record represents a single record to be ingested
type Record struct {
	ID       string                 `json:"id"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IngestResponse represents the response from batch ingestion
type IngestResponse struct {
	Successful  int      `json:"successful"`
	Failed      int      `json:"failed"`
	FailedIDs   []string `json:"failed_ids,omitempty"`
	Message     string   `json:"message"`
	ProcessedAt string   `json:"processed_at"`
}

// MigrationService handles the business logic for fugu migration
type MigrationService struct {
	fuguServerURL string
}

// NewMigrationService creates a new migration service
func NewMigrationService(fuguServerURL string) *MigrationService {
	return &MigrationService{
		fuguServerURL: fuguServerURL,
	}
}

// MigrationServiceHandler handles HTTP requests for migration
type MigrationServiceHandler struct {
	service *MigrationService
}

// NewMigrationHandler creates a new migration handler
func NewMigrationHandler(service *MigrationService) *MigrationServiceHandler {
	return &MigrationServiceHandler{
		service: service,
	}
}

// RegisterMigrationRoutes registers migration routes with the router
func RegisterMigrationRoutes(r *mux.Router) error {
	fuguServerURL := "http://fugudb"
	service := NewMigrationService(fuguServerURL)
	msh := NewMigrationHandler(service)

	migrationRoute := r.PathPrefix("/migration").Subrouter()
	migrationRoute.HandleFunc(
		"/batch-ingest",
		msh.BatchIngest,
	).Methods(http.MethodPost)

	migrationRoute.HandleFunc(
		"/health",
		msh.Health,
	).Methods(http.MethodGet)

	return nil
}

// Health handles health check requests
func (h *MigrationServiceHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "migration-api:health")
	defer span.End()

	logger.Info(ctx, "migration health check called")

	// Create fugu client and check health
	client, err := fugusdk.NewClient(ctx, h.service.fuguServerURL)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client", zap.Error(err))
		http.Error(w, "Fugu server unavailable", http.StatusServiceUnavailable)
		return
	}

	if err := client.Health(ctx); err != nil {
		logger.Error(ctx, "fugu server health check failed", zap.Error(err))
		http.Error(w, "Fugu server unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "fugu-migration",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode health response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// BatchIngest handles batch ingestion requests
func (h *MigrationServiceHandler) BatchIngest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "migration-api:batch-ingest")
	defer span.End()

	logger.Info(ctx, "batch ingest called")

	// Parse request body
	var ingestReq IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&ingestReq); err != nil {
		logger.Error(ctx, "failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(ingestReq.Records) == 0 {
		logger.Error(ctx, "no records provided in request")
		http.Error(w, "No records provided", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "processing batch ingest request", zap.Int("record_count", len(ingestReq.Records)))

	// Process the ingestion
	response, err := h.service.ProcessBatchIngest(ctx, ingestReq.Records)
	if err != nil {
		logger.Error(ctx, "batch ingest processing failed", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response status based on results
	w.Header().Set("Content-Type", "application/json")
	if response.Failed > 0 && response.Successful == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else if response.Failed > 0 {
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error(ctx, "failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "batch ingest completed",
		zap.Int("successful", response.Successful),
		zap.Int("failed", response.Failed))
}

// ProcessBatchIngest handles the business logic for batch ingestion
func (s *MigrationService) ProcessBatchIngest(ctx context.Context, records []Record) (*IngestResponse, error) {
	ctx, span := tracer.Start(ctx, "migration-service:process-batch-ingest")
	defer span.End()

	// Create fugu client
	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Convert records to fugu format
	fuguObjects, conversionErrors := s.convertToFuguObjects(ctx, records)

	// Initialize response
	response := &IngestResponse{
		ProcessedAt: time.Now().Format(time.RFC3339),
		FailedIDs:   make([]string, 0),
	}

	// Add conversion errors to failed count
	for recordID, err := range conversionErrors {
		logger.Error(ctx, "record conversion failed",
			zap.String("record_id", recordID),
			zap.Error(err))
		response.Failed++
		response.FailedIDs = append(response.FailedIDs, recordID)
	}

	// Ingest valid objects if any
	if len(fuguObjects) > 0 {
		if err := client.IngestObjects(ctx, fuguObjects); err != nil {
			logger.Error(ctx, "fugu ingestion failed", zap.Error(err))

			// Mark all remaining objects as failed
			for _, obj := range fuguObjects {
				response.Failed++
				response.FailedIDs = append(response.FailedIDs, obj.ID)
			}

			response.Message = "Ingestion failed"
			return response, nil
		}

		// Ingestion successful
		response.Successful = len(fuguObjects)
		logger.Info(ctx, "successfully ingested records", zap.Int("count", response.Successful))
	}

	// Set final message
	if response.Failed == 0 {
		response.Message = "All records ingested successfully"
	} else if response.Successful == 0 {
		response.Message = "All records failed to ingest"
	} else {
		response.Message = fmt.Sprintf("Partial success: %d succeeded, %d failed", response.Successful, response.Failed)
	}

	return response, nil
}

// convertToFuguObjects converts incoming records to fugu ObjectRecord format
func (s *MigrationService) convertToFuguObjects(ctx context.Context, records []Record) ([]fugusdk.ObjectRecord, map[string]error) {
	ctx, span := tracer.Start(ctx, "migration-service:convert-to-fugu-objects")
	defer span.End()

	fuguObjects := make([]fugusdk.ObjectRecord, 0, len(records))
	conversionErrors := make(map[string]error)

	for _, record := range records {
		// Validate required fields
		if record.ID == "" {
			conversionErrors["unknown"] = fmt.Errorf("record ID cannot be empty")
			continue
		}

		if record.Text == "" {
			conversionErrors[record.ID] = fmt.Errorf("record text cannot be empty")
			continue
		}

		// Create fugu object
		fuguObj := fugusdk.ObjectRecord{
			ID:       record.ID,
			Text:     record.Text,
			Metadata: record.Metadata,
		}

		// Add migration metadata
		if fuguObj.Metadata == nil {
			fuguObj.Metadata = make(map[string]interface{})
		}

		fuguObj.Metadata["migrated_at"] = time.Now().Format(time.RFC3339)
		fuguObj.Metadata["migration_source"] = "kessler-batch-ingest"

		// Parse date from metadata if available (matching Python script logic)
		if dateStr, ok := fuguObj.Metadata["date"].(string); ok {
			if parsedTime, err := s.parseDate(dateStr); err == nil {
				fuguObj.Metadata["date_iso"] = parsedTime.Format(time.RFC3339)
			} else {
				logger.Warn(ctx, "could not parse date from metadata",
					zap.String("record_id", record.ID),
					zap.String("date", dateStr),
					zap.Error(err))
			}
		}

		fuguObjects = append(fuguObjects, fuguObj)
	}

	return fuguObjects, conversionErrors
}

// parseDate attempts to parse various date formats (similar to Python script)
func (s *MigrationService) parseDate(dateStr string) (time.Time, error) {
	// Try M/D/Y format first (as used in Python script)
	layouts := []string{
		"1/2/2006",   // M/D/Y
		"01/02/2006", // MM/DD/YYYY
		"2006-01-02", // YYYY-MM-DD
		"2006/01/02", // YYYY/MM/DD
		time.RFC3339, // ISO format
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
