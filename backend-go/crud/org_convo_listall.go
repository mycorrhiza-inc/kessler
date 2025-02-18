package crud

import (
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/common/objects/networking"
	"kessler/util"
	"net/http"
	"time"

	"github.com/charmbracelet/log"

	"github.com/google/uuid"
)

func OrgSemiCompletePaginated(w http.ResponseWriter, r *http.Request) {
	paginationData := networking.PaginationFromUrlParams(r)
	q := util.DBQueriesFromRequest(r)
	ctx := r.Context()
	args := dbstore.OrganizationSemiCompleteInfoListPaginatedParams{
		Limit:  int32(paginationData.Limit),
		Offset: int32(paginationData.Offset),
	}
	organizations, err := q.OrganizationSemiCompleteInfoListPaginated(ctx, args)
	if err != nil {
		log.Info(fmt.Sprintf("Error reading organization: %v", err))
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

func ConversationSemiCompletePaginatedList(w http.ResponseWriter, r *http.Request) {
	paginationData := networking.PaginationFromUrlParams(r)
	q := util.DBQueriesFromRequest(r)
	ctx := r.Context()
	args := dbstore.ConversationSemiCompleteInfoListPaginatedParams{
		Limit:  int32(paginationData.Limit),
		Offset: int32(paginationData.Offset),
	}
	proceedings, err := q.ConversationSemiCompleteInfoListPaginated(ctx, args)
	if err != nil {
		log.Info(fmt.Sprintf("Error reading organization: %v", err))
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
	log.Info(fmt.Sprintf("Getting all proceedings"))
	q := *util.DBQueriesFromRequest(r)

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
			DocketId:      proceeding_raw.DocketGovID,
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
