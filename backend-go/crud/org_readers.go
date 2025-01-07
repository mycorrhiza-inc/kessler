package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/files"
	"kessler/objects/organizations"
	"kessler/routing"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func OrgGetWithFilesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting file with metadata")
	q := *routing.DBQueriesFromRequest(r)

	params := mux.Vars(r)
	orgID := params["uuid"]
	parsedUUID, err := uuid.Parse(orgID)
	if err != nil {
		http.Error(w, "Invalid File ID format", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	complete_org_info, err := OrgWithFilesGetByID(ctx, &q, parsedUUID)
	if err != nil {
		log.Printf("Error reading organization: %v", err)
		if err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(complete_org_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func OrgWithFilesGetByUnknownQuery(ctx context.Context, q *dbstore.Queries, orgQuery string) (organizations.OrganizationSchemaComplete, error) {
	if orgQuery == "" {
		return organizations.OrganizationSchemaComplete{}, fmt.Errorf("No query string provided.")
	}
	org_id, err := uuid.Parse(orgQuery)
	if err == nil {
		return OrgWithFilesGetByID(ctx, q, org_id)
	}
	arg := dbstore.OrganizationAliasIdNameGetParams{

	}
	test, err := q.OrganizationAliasIdNameGet(ctx, arg dbstore.OrganizationAliasIdNameGetParams)
}

func OrgWithFilesGetByID(ctx context.Context, q *dbstore.Queries, orgID uuid.UUID) (organizations.OrganizationSchemaComplete, error) {
	org_info, err := q.OrganizationRead(ctx, orgID)
	if err != nil {
		return organizations.OrganizationSchemaComplete{}, err
	}

	org_files_raw, err := q.AuthorshipOrganizationListDocuments(ctx, orgID)
	if err != nil {
		return organizations.OrganizationSchemaComplete{}, err
	}

	org_files := make([]files.FileSchema, len(org_files_raw))
	org_file_ids := make([]uuid.UUID, len(org_files_raw))
	for i, f := range org_files_raw {
		file_uuid := f.DocumentID
		org_files[i] = files.FileSchema{
			ID: file_uuid,
		}
		org_file_ids[i] = file_uuid
	}

	return organizations.OrganizationSchemaComplete{
		ID:               orgID,
		Name:             org_info.Name,
		FilesAuthored:    org_files,
		FilesAuthoredIDs: org_file_ids,
	}, nil
}
