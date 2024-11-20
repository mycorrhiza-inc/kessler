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

func GetOrgWithFilesFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
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
			log.Printf("Error reading organization: %v", err)
			if err.Error() == "no rows in result set" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		org_files_raw, err := q.AuthorshipOrganizationListDocuments(
			ctx,
			pgUUID)
		if err != nil {
			log.Printf("Error reading organization: %v", err)
			if err.Error() == "no rows in result set" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		org_files := make([]FileSchema, len(org_files_raw))
		org_file_ids := make([]uuid.UUID, len(org_files_raw))
		for i, f := range org_files_raw {
			file_uuid := f.DocumentID.Bytes
			org_files[i] = FileSchema{
				ID: file_uuid,
			}
			org_file_ids[i] = file_uuid
		}
		complete_org_info := OrganizationSchemaComplete{
			ID:               parsedUUID,
			Name:             org_info.Name,
			FilesAuthored:    org_files,
			FilesAuthoredIDs: org_file_ids,
		}

		response, _ := json.Marshal(complete_org_info)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
