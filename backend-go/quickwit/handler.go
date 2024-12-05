package quickwit

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbstore "kessler/gen/dbstore"
	util "kessler/util"

	"github.com/gorilla/mux"
)

func HandleQuickwitIngestFromPostgresFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
	}
}

type SearchIngestRequest struct {
	Index string                   `json:"index"`
	Data  []map[string]interface{} `json:"data"`
}

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

type ClearIndexRequest struct {
	Index string `json:"index"`
}

func HandleQuickwitIndexClear(w http.ResponseWriter, r *http.Request) {
	var req ClearIndexRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
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
