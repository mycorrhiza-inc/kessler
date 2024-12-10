package crud

import (
	"encoding/json"
	"fmt"
	"kessler/objects/files"
	"kessler/objects/organizations"
	"kessler/routing"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetOrgWithFiles(w http.ResponseWriter, r *http.Request) {
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
	// TODO: Get these 2 requests to happen in the same query, and or run concurrently
	org_info, err := q.OrganizationRead(ctx, parsedUUID)
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
		parsedUUID)
	if err != nil {
		log.Printf("Error reading organization: %v", err)
		if err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	complete_org_info := organizations.OrganizationSchemaComplete{
		ID:               parsedUUID,
		Name:             org_info.Name,
		FilesAuthored:    org_files,
		FilesAuthoredIDs: org_file_ids,
	}

	response, _ := json.Marshal(complete_org_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ConversationGetByName(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting file with metadata")
	q := *routing.DBQueriesFromRequest(r)

	params := mux.Vars(r)
	docketIdStr := params["name"]
	ctx := r.Context()
	conv_infos, err := q.DocketConversationFetchByDocketIdMatch(ctx, docketIdStr)
	if err != nil {
		log.Printf("Error reading organization: %v", err)
		if err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(conv_infos) == 0 {
		errorstr := fmt.Sprintf("No proceeding with name %s", docketIdStr)
		http.Error(w, errorstr, http.StatusNotFound)
		fmt.Println(errorstr)
		return
	}
	conv_info := conv_infos[0]
	response, _ := json.Marshal(conv_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func OrgListAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting all organizations")
	q := *routing.DBQueriesFromRequest(r)

	ctx := r.Context()
	organizations, err := q.OrganizationList(ctx)
	if err != nil {
		log.Printf("Error reading organization: %v", err)
		if err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(organizations)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ConversationListAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting all proceedings")
	q := *routing.DBQueriesFromRequest(r)

	ctx := r.Context()
	proceedings, err := q.DocketConversationList(ctx)
	if err != nil {
		log.Printf("Error reading organization: %v", err)
		if err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(proceedings)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

type ConversationSemiCompleteInfo struct {
	ID            uuid.UUID              `json:"id"`
	DocketId      string                 `json:"docket_id"`
	Name          string                 `json:"name"`
	Description   map[string]interface{} `json:"description"`
	DocumentCount int                    `json:"document_count"`
	DateCreated   time.Time              `json:"date_created"`
	DateUpdated   time.Time              `json:"date_updated"`
}

func ConversationSemiCompleteListAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting all proceedings")
	q := *routing.DBQueriesFromRequest(r)

	ctx := r.Context()
	proceedings_raw, err := q.ConversationSemiCompleteInfoList(ctx)
	proceedings := make([]ConversationSemiCompleteInfo, len(proceedings_raw))
	for i, proceeding_raw := range proceedings_raw {
		parsed_description := map[string]interface{}{}
		err := json.Unmarshal([]byte(proceeding_raw.Description), &parsed_description)
		if err != nil {
			continue
		}
		proceedings[i] = ConversationSemiCompleteInfo{
			ID:            proceeding_raw.ID,
			DocketId:      proceeding_raw.DocketID,
			Name:          proceeding_raw.Name,
			Description:   parsed_description,
			DocumentCount: int(proceeding_raw.DocumentCount),
			DateCreated:   proceeding_raw.CreatedAt.Time,
			DateUpdated:   proceeding_raw.UpdatedAt.Time,
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(proceedings)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
