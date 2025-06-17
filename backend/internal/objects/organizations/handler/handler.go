package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/internal/dbstore"
	"kessler/internal/objects/authors"
	"kessler/internal/objects/files"
	"kessler/internal/objects/networking"
	"kessler/internal/objects/organizations"
	"kessler/pkg/database"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

// OrgHandler holds dependencies for organization operations
type OrgHandler struct {
	db dbstore.DBTX
}

// NewFileHandler creates a new file handler with the given database connection
func NewOrgHandler(db dbstore.DBTX) *OrgHandler {
	return &OrgHandler{
		db: db,
	}
}

func DefineOrganizationRoutes(r *mux.Router, db dbstore.DBTX) {
	organizationsRouter := r.PathPrefix("/organizations").Subrouter()
	handler := NewOrgHandler(db)

	r.HandleFunc(
		"/list",
		handler.OrgSemiCompletePaginated,
	).Methods(http.MethodGet)

	organizationsRouter.HandleFunc(
		"/{uuid}",
		handler.OrgGetWithFilesHandler,
	).Methods(http.MethodGet)

	organizationsRouter.HandleFunc(
		"/verify",
		handler.OrganizationVerifyHandler,
	).Methods(http.MethodPost)
}

type OrganizationRequest struct {
	OrganizationName string `json:"organization_name"`
	IsPerson         bool   `json:"is_person"`
}

func (h *OrgHandler) OrganizationVerifyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "organizations:OrganizationVerifyHandler")
	defer span.End()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	var req OrganizationRequest
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		errorstring := fmt.Sprintf("Error decoding JSON: %v\n", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	q := *database.GetTx()

	author_info := authors.AuthorInformation{AuthorName: req.OrganizationName, IsPerson: req.IsPerson}
	author_info, err = VerifyAuthorOrganizationUUID(ctx, q, &author_info)
	if err != nil {
		errorstring := fmt.Sprintf("Error verifying author %v: %v\n", req.OrganizationName, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	// No error handling since we always want it to retun a 200 at this point.
	response, _ := json.Marshal(author_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func VerifyAuthorOrganizationUUID(ctx context.Context, q dbstore.Queries, author_info *authors.AuthorInformation) (authors.AuthorInformation, error) {
	if author_info.AuthorID != uuid.Nil {
		return *author_info, nil
	}
	// TODO: Change the sql so that this also matches IsPerson, but for now it shouldnt matter.
	org_return_info, err := q.OrganizationFetchByAliasMatchSingle(ctx, author_info.AuthorName)
	if err == nil {
		org_uuid := org_return_info.ID.Bytes
		author_info.AuthorID = org_uuid
		author_info.IsPerson = org_return_info.IsPerson.Bool
		return *author_info, nil
	}
	org_create_params := dbstore.CreateOrganizationParams{
		OrganizationAlias: author_info.AuthorName,
		Description:       "", // I should make this fixable at some point, but for now it will kinda work (tm)
		IsPerson:          pgtype.Bool{Bool: author_info.IsPerson, Valid: true},
	}
	pg_org_id, err := q.CreateOrganization(ctx, org_create_params)
	if err != nil {
		log.Info(err)
		return *author_info, err
	}
	org_uuid := pg_org_id
	author_info.AuthorID = org_uuid
	return *author_info, nil
}

func (h *OrgHandler) OrgGetWithFilesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "organizations:OrgGetWithFilesHandler")
	defer span.End()
	log.Info(fmt.Sprintf("Getting file with metadata"))
	q := database.GetTx()

	params := mux.Vars(r)
	orgID := params["uuid"]
	parsedUUID, err := uuid.Parse(orgID)
	if err != nil {
		http.Error(w, "Invalid File ID format", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	complete_org_info, err := OrgWithFilesGetByID(ctx, q, parsedUUID)
	if err != nil {
		log.Info(fmt.Sprintf("Error reading organization: %v", err))
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

func OrgWithFilesGetByUnknown(ctx context.Context, q *dbstore.Queries, orgQuery string) (organizations.OrganizationSchemaComplete, error) {
	if orgQuery == "" {
		return organizations.OrganizationSchemaComplete{}, fmt.Errorf("No query string provided.")
	}
	org_id, err := uuid.Parse(orgQuery)
	if err == nil {
		return OrgWithFilesGetByID(ctx, q, org_id)
	}
	orgs_pgobject, err := q.OrganizationFetchByAliasMatchSingle(ctx, orgQuery)
	if err != nil {
		return organizations.OrganizationSchemaComplete{}, err
	}
	org_id = orgs_pgobject.ID.Bytes
	return OrgWithFilesGetByID(ctx, q, org_id)
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

	org_aliases_pg, err := q.OrganizationGetAllAliases(ctx, orgID)
	if err != nil {
		return organizations.OrganizationSchemaComplete{}, err
	}
	org_aliases := make([]string, len(org_aliases_pg))
	for index, alias_type := range org_aliases_pg {
		org_aliases[index] = alias_type.OrganizationAlias
	}
	if len(org_aliases) == 0 {
		org_aliases = []string{org_info.Name}
	}

	return organizations.OrganizationSchemaComplete{
		ID:               orgID,
		Name:             org_info.Name,
		Aliases:          org_aliases,
		FilesAuthored:    org_files,
		FilesAuthoredIDs: org_file_ids,
	}, nil
}

func (h *OrgHandler) OrgSemiCompletePaginated(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "organizations:OrgSemiCompletePaginated")
	defer span.End()
	paginationData := networking.PaginationFromUrlParams(r)
	q := database.GetQueries(h.db)
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
