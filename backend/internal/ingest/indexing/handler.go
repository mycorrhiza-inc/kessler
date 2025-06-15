package indexing

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kessler/pkg/logger"
)

// IndexHandler provides HTTP endpoints for admin indexing.
// @Tags Indexing
// @Description Admin endpoints for indexing conversations and organizations into FuguDB
// @Produce json
// @Router /admin/indexing [basePath]
type IndexHandler struct {
	svc *IndexService
}

// NewIndexHandler creates a new handler backed by IndexService.
func NewIndexHandler(svc *IndexService) *IndexHandler {
	return &IndexHandler{svc: svc}
}

// RegisterAdminIndexingRoutes mounts indexing endpoints under /admin/indexing.
func RegisterAdminIndexingRoutes(r *mux.Router) {
	sr := r.PathPrefix("/indexing").Subrouter()
	h := NewIndexHandler(NewIndexService("http://fugudb:3301"))

	sr.HandleFunc("/conversations", h.IndexAllConversations).Methods(http.MethodPost)
	sr.HandleFunc("/conversations/{id}", h.IndexConversationByID).Methods(http.MethodPost)
	sr.HandleFunc("/organizations", h.IndexAllOrganizations).Methods(http.MethodPost)
	sr.HandleFunc("/organizations/{id}", h.IndexOrganizationByID).Methods(http.MethodPost)
}

// IndexAllConversations godoc
// @Summary Batch index all conversations
// @Description Retrieves all conversations from the database and indexes them in FuguDB
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/conversations [post]
func (h *IndexHandler) IndexAllConversations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	count, err := h.svc.IndexAllConversations(ctx)
	if err != nil {
		logger.Error(ctx, "index all conversations failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// IndexConversationByID godoc
// @Summary Index a conversation by ID
// @Description Retrieves a single conversation by UUID and indexes it in FuguDB
// @Param id path string true "Conversation UUID"
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/conversations/{id} [post]
func (h *IndexHandler) IndexConversationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	count, err := h.svc.IndexConversationByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "index conversation failed", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// IndexAllOrganizations godoc
// @Summary Batch index all organizations
// @Description Retrieves all organizations from the database and indexes them in FuguDB
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 500 {object} map[string]string{"error":string}
// @Router /admin/indexing/organizations [post]
func (h *IndexHandler) IndexAllOrganizations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	count, err := h.svc.IndexAllOrganizations(ctx)
	if err != nil {
		logger.Error(ctx, "index all organizations failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// IndexOrganizationByID godoc
// @Summary Index an organization by ID
// @Description Retrieves a single organization by UUID and indexes it in FuguDB
// @Param id path string true "Organization UUID"
// @Success 200 {object} map[string]int{"indexed":int}
// @Failure 400 {object} map[string]string{"error":string}
// @Router /admin/indexing/organizations/{id} [post]
func (h *IndexHandler) IndexOrganizationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	count, err := h.svc.IndexOrganizationByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "index organization failed", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]int{"indexed": count})
}

// respondJSON writes a JSON response.
func (h *IndexHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondError writes an error string as JSON.
func (h *IndexHandler) respondError(w http.ResponseWriter, status int, msg string) {
	h.respondJSON(w, status, map[string]string{"error": msg})
}
