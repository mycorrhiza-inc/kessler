// indexing/handler.go
package indexing

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"kessler/internal/dbstore"
	"kessler/pkg/logger"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// IndexHandler provides HTTP endpoints for admin indexing.
// @Tags Indexing
// @Description Admin endpoints for indexing conversations, organizations, and batch data into FuguDB
// @Produce json
// @Router /admin/indexing [basePath]
type IndexHandler struct {
	svc *IndexService
}

var tracer = otel.Tracer("http-server")

// NewIndexHandler creates a new handler backed by IndexService.
func NewIndexHandler(svc *IndexService) *IndexHandler {
	return &IndexHandler{svc: svc}
}

// RegisterAdminIndexingRoutes mounts indexing endpoints under /admin/indexing.
func RegisterIndexingRoutes(r *mux.Router, db dbstore.DBTX) {
	sr := r.PathPrefix("/indexing").Subrouter()

	// TODO: have this be configurable
	h := NewIndexHandler(NewIndexService("http://fugudb:3301", db))

	// Conversation endpoints
	sr.HandleFunc("/conversations", h.IndexAllConversations).Methods(http.MethodPost)
	sr.HandleFunc("/conversations/{id}", h.IndexConversationByID).Methods(http.MethodPost)
	sr.HandleFunc("/conversations/{id}", h.DeleteConversation).Methods(http.MethodDelete)

	// Organization endpoints
	sr.HandleFunc("/organizations", h.IndexAllOrganizations).Methods(http.MethodPost)
	sr.HandleFunc("/organizations/{id}", h.IndexOrganizationByID).Methods(http.MethodPost)
	sr.HandleFunc("/organizations/{id}", h.DeleteOrganization).Methods(http.MethodDelete)

	// Data endpoints - full CRUD operations with namespace facet support
	sr.HandleFunc("/data", h.IndexAllData).Methods(http.MethodPost)
	sr.HandleFunc("/data/batch-ingest", h.BatchIngestData).Methods(http.MethodPost)
	sr.HandleFunc("/data/{id}", h.IndexDataRecordByID).Methods(http.MethodPost)
	sr.HandleFunc("/data/{id}", h.DeleteDataRecord).Methods(http.MethodDelete)
	sr.HandleFunc("/data/health", h.DataHealthCheck).Methods(http.MethodGet)

	// Bulk operations
	sr.HandleFunc("/all", h.IndexAllData).Methods(http.MethodPost)
}

// IndexAllConversations godoc
// @Summary Batch index all conversations
// @Description Retrieves all conversations from the database and indexes them in FuguDB with proper namespace facets
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/conversations [post]
func (h *IndexHandler) IndexAllConversations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:IndexAllConversations")
	defer span.End()

	logger.Info(ctx, "indexing all conversations requested")

	count, err := h.svc.IndexAllConversations(ctx)
	if err != nil {
		logger.Error(ctx, "index all conversations failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info(ctx, "successfully indexed all conversations", zap.Int("count", count))
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// IndexConversationByID godoc
// @Summary Index a conversation by ID
// @Description Retrieves a single conversation by UUID and indexes it in FuguDB with proper namespace facets
// @Param id path string true "Conversation UUID"
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/conversations/{id} [post]
func (h *IndexHandler) IndexConversationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:IndexConversationByID")
	defer span.End()

	id := mux.Vars(r)["id"]

	logger.Info(ctx, "indexing conversation by ID", zap.String("conversation_id", id))

	count, err := h.svc.IndexConversationByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "index conversation failed", zap.String("conversation_id", id), zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "successfully indexed conversation", zap.String("conversation_id", id))
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// DeleteConversation godoc
// @Summary Delete a conversation from the search index
// @Description Removes a conversation from the FuguDB search index
// @Param id path string true "Conversation UUID"
// @Success 200 {object} map[string]string{"message":string}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/conversations/{id} [delete]
func (h *IndexHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:DeleteConversation")
	defer span.End()

	id := mux.Vars(r)["id"]

	logger.Info(ctx, "deleting conversation from index", zap.String("conversation_id", id))

	err := h.svc.DeleteConversationFromIndex(ctx, id)
	if err != nil {
		logger.Error(ctx, "delete conversation failed", zap.String("conversation_id", id), zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "successfully deleted conversation from index", zap.String("conversation_id", id))
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Conversation deleted from index"})
}

// IndexAllOrganizations godoc
// @Summary Batch index all organizations
// @Description Retrieves all organizations from the database and indexes them in FuguDB with proper namespace facets
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/organizations [post]
func (h *IndexHandler) IndexAllOrganizations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:IndexAllOrganizations")
	defer span.End()

	logger.Info(ctx, "indexing all organizations requested")

	count, err := h.svc.IndexAllOrganizations(ctx)
	if err != nil {
		logger.Error(ctx, "index all organizations failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info(ctx, "successfully indexed all organizations", zap.Int("count", count))
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// IndexOrganizationByID godoc
// @Summary Index an organization by ID
// @Description Retrieves a single organization by UUID and indexes it in FuguDB with proper namespace facets
// @Param id path string true "Organization UUID"
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/organizations/{id} [post]
func (h *IndexHandler) IndexOrganizationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:IndexOrganizationByID")
	defer span.End()

	id := mux.Vars(r)["id"]

	logger.Info(ctx, "indexing organization by ID", zap.String("organization_id", id))

	count, err := h.svc.IndexOrganizationByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "index organization failed", zap.String("organization_id", id), zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "successfully indexed organization", zap.String("organization_id", id))
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// DeleteOrganization godoc
// @Summary Delete an organization from the search index
// @Description Removes an organization from the FuguDB search index
// @Param id path string true "Organization UUID"
// @Success 200 {object} map[string]string{"message":string}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/organizations/{id} [delete]
func (h *IndexHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:DeleteOrganization")
	defer span.End()

	id := mux.Vars(r)["id"]

	logger.Info(ctx, "deleting organization from index", zap.String("organization_id", id))

	err := h.svc.DeleteOrganizationFromIndex(ctx, id)
	if err != nil {
		logger.Error(ctx, "delete organization failed", zap.String("organization_id", id), zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "successfully deleted organization from index", zap.String("organization_id", id))
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Organization deleted from index"})
}

// IndexAllData godoc
// @Summary Index all conversations and organizations
// @Description Bulk operation to index all conversations and organizations from the database with proper namespace facets
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/all [post]
func (h *IndexHandler) IndexAllData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:IndexAllData")
	defer span.End()

	logger.Info(ctx, "indexing all data requested")

	convCount, orgCount, err := h.svc.IndexAllData(ctx)
	if err != nil {
		logger.Error(ctx, "index all data failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"conversations_indexed": convCount,
		"organizations_indexed": orgCount,
		"total_indexed":         convCount + orgCount,
		"message":               "Successfully indexed all data with namespace facets",
	}

	logger.Info(ctx, "successfully indexed all data",
		zap.Int("conversations", convCount),
		zap.Int("organizations", orgCount))
	h.respondJSON(w, http.StatusOK, response)
}

// BatchIngestData godoc
// @Summary Batch ingest arbitrary data records
// @Description Accepts a batch of data records and indexes them in FuguDB with namespace facet support
// @Accept json
// @Param request body DataIngestRequest true "Batch ingest request with namespace facet fields"
// @Success 200 {object} DataIngestResponse
// @Success 206 {object} DataIngestResponse "Partial success"
// @Failure 400 {object} map[string]string{"error":string}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/data/batch-ingest [post]
func (h *IndexHandler) BatchIngestData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:BatchIngestData")
	defer span.End()

	logger.Info(ctx, "batch data ingest request received")

	// Parse request body
	var ingestReq DataIngestRequest
	if err := json.NewDecoder(r.Body).Decode(&ingestReq); err != nil {
		logger.Error(ctx, "failed to decode batch ingest request body", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(ingestReq.Records) == 0 {
		logger.Error(ctx, "no records provided in batch ingest request")
		h.respondError(w, http.StatusBadRequest, "No records provided")
		return
	}

	logger.Info(ctx, "processing batch data ingest request",
		zap.Int("record_count", len(ingestReq.Records)))

	// Process the ingestion
	response, err := h.svc.ProcessBatchDataIngest(ctx, ingestReq.Records)
	if err != nil {
		logger.Error(ctx, "batch data ingest processing failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "Internal server error")
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
		logger.Error(ctx, "failed to encode batch ingest response", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	logger.Info(ctx, "batch data ingest completed",
		zap.Int("successful", response.Successful),
		zap.Int("failed", response.Failed))
}

// IndexDataRecordByID godoc
// @Summary Index a single data record by ID
// @Description Accepts a single data record and indexes it in FuguDB with namespace facet support
// @Accept json
// @Param id path string true "Data Record ID"
// @Param request body DataRecord true "Data record to index with namespace facet fields"
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 400 {object} map[string]string{"error":string}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/data/{id} [post]
func (h *IndexHandler) IndexDataRecordByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:IndexDataRecordByID")
	defer span.End()

	id := mux.Vars(r)["id"]

	logger.Info(ctx, "indexing data record by ID", zap.String("data_record_id", id))

	// Parse request body
	var record DataRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		logger.Error(ctx, "failed to decode data record request body", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Ensure the ID in the URL matches the record ID
	if record.ID == "" {
		record.ID = id
	} else if record.ID != id {
		logger.Error(ctx, "data record ID mismatch",
			zap.String("url_id", id),
			zap.String("body_id", record.ID))
		h.respondError(w, http.StatusBadRequest, "ID in URL must match ID in request body")
		return
	}

	// Process single record as a batch of one
	count, err := h.svc.IndexDataRecordByID(ctx, record)
	if err != nil {
		logger.Error(ctx, "index data record failed", zap.String("data_record_id", id), zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "successfully indexed data record", zap.String("data_record_id", id))
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// DeleteDataRecord godoc
// @Summary Delete a data record from the search index
// @Description Removes a data record from the FuguDB search index
// @Param id path string true "Data Record ID"
// @Success 200 {object} map[string]string{"message":string}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/data/{id} [delete]
func (h *IndexHandler) DeleteDataRecord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:DeleteDataRecord")
	defer span.End()

	id := mux.Vars(r)["id"]

	logger.Info(ctx, "deleting data record from index", zap.String("data_record_id", id))

	err := h.svc.DeleteDataRecordFromIndex(ctx, id)
	if err != nil {
		logger.Error(ctx, "delete data record failed", zap.String("data_record_id", id), zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "successfully deleted data record from index", zap.String("data_record_id", id))
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Data record deleted from index"})
}

// DataHealthCheck godoc
// @Summary Health check for data ingestion
// @Description Checks the health of the FuguDB connection for data ingestion
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]string{"error":string}
// @Router /admin/indexing/data/health [get]
func (h *IndexHandler) DataHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "indexing:DataHealthCheck")
	defer span.End()

	logger.Info(ctx, "data health check requested")

	// Check FuguDB health through the service
	client, err := h.svc.createFuguClient(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create fugu client for health check", zap.Error(err))
		h.respondError(w, http.StatusServiceUnavailable, "FuguDB unavailable")
		return
	}

	if err := client.Health(ctx); err != nil {
		logger.Error(ctx, "fugu health check failed", zap.Error(err))
		h.respondError(w, http.StatusServiceUnavailable, "FuguDB health check failed")
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "kessler-indexing",
		"fugudb":    "healthy",
	}

	logger.Info(ctx, "health check successful")
	h.respondJSON(w, http.StatusOK, response)
}

// parsePageParams extracts and validates pagination parameters from the request
func (h *IndexHandler) parsePageParams(r *http.Request) (page, perPage int) {
	page = 1
	perPage = 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	return page, perPage
}

// respondJSON writes a JSON response.
func (h *IndexHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Log encoding error but don't try to respond again
		logger.Error(nil, "failed to encode JSON response", zap.Error(err))
	}
}

// respondError writes an error string as JSON.
func (h *IndexHandler) respondError(w http.ResponseWriter, status int, msg string) {
	h.respondJSON(w, status, map[string]string{"error": msg})
}
