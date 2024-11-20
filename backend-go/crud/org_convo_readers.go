package crud

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func GetOrgWithFiles(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Getting file with metadata")
		q := *dbstore.New(dbtx_val)
		params := mux.Vars(r)
		orgID := params["uuid"]
		parsedUUID, err := uuid.Parse(orgID)
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}
		pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}
		ctx := r.Context()
		// TODO: Get these 2 requests to happen in the same query, and or run concurrently
		org_info, err := q.OrganizationRead(ctx, pgUUID)
		if err != nil {
			http.Error(w, "Organization not found", http.StatusNotFound)
			return
		}
		response, _ := json.Marshal(org_info)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
