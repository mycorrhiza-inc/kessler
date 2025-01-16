package quickwit

import (
	"encoding/json"
	"fmt"
	"net/http"

	"kessler/crud"
	"kessler/routing"
	util "kessler/util"

	"github.com/gorilla/mux"
)

type SearchIngestRequest struct {
	Index string                   `json:"index"`
	Data  []QuickwitFileUploadData `json:"data"`
}

type ClearIndexRequest struct {
	Index string `json:"index"`
}

type IndexConversationsRequest struct {
	Index string   `json:"index"`
	Ids   []string `json:"ids"`
}

func DefineRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/v2/engine/search").Subrouter()
	subrouter.HandleFunc(
		"/ingest",
		HandleQuckwitIngest,
	).Methods(http.MethodPost)
	subrouter.HandleFunc(
		"/ingest",
		HandleQuickwitIndexClear,
	).Methods(http.MethodPost)
}

// Ingest Requests
func HandleQuckwitIngest(w http.ResponseWriter, r *http.Request) {
	var req SearchIngestRequest
	util.DecodeRequest(r.Body, &req, w)
	defer r.Body.Close()
	err := IngestIntoIndex(req.Index, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "injest job started")
}
func HandleQuckwitIngestConversations(w http.ResponseWriter, r *http.Request) {
	var indexReq IndexConversationsRequest
	util.DecodeRequest(r.Body, &indexReq, w)

	ids, err := util.ParseStringSliceUUIDs(indexReq.Ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := *routing.DBQueriesFromRequest(r)
	ctx := r.Context()

	conversations, err := crud.ConversationGetListByUUID(ctx, &q, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	IndexConversations(conversations)
}

func HandleQuckwitIngestOrganizations(w http.ResponseWriter, r *http.Request) {

}

// Clear Requests
func HandleQuickwitIndexClear(w http.ResponseWriter, r *http.Request) {
	var req ClearIndexRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
}
