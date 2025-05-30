package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/internal/database"
	"kessler/internal/dbstore"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/networking"
	"kessler/pkg/util"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

func DefineConversationsRoutes(r *mux.Router) {
	// ------- conversations subroute -------
	conversationsRouter := r.PathPrefix("/conversations").Subrouter()

	conversationsRouter.HandleFunc(
		"/list",
		ConversationSemiCompletePaginatedList,
	).Methods(http.MethodGet)

	conversationsRouter.HandleFunc(
		"/named-lookup/{name}",
		ConversationGetByUnknownHandler,
	).Methods(http.MethodGet)

	conversationsRouter.HandleFunc(
		"/verify",
		ConversationVerifyHandler,
	).Methods(http.MethodPost)

	conversationsRouter.HandleFunc(
		"/list/semi-complete",
		ConversationSemiCompleteListAll,
	).Methods(http.MethodGet)
}

func ConversationGetByUnknownHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(fmt.Sprintf("Getting file with metadata"))
	q := database.GetTx()

	params := mux.Vars(r)
	docketIdStr := params["name"]
	ctx := r.Context()

	conv_info, err := ConversationGetByUnknown(ctx, q, docketIdStr)
	if err != nil {
		log.Info(fmt.Sprintf("Error reading organization: %v", err))
		if err.Error() == "no rows in result set" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(conv_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ConversationGetByUnknown(ctx context.Context, q *dbstore.Queries, query_string string) (dbstore.DocketConversation, error) {
	if query_string == "" {
		return dbstore.DocketConversation{}, fmt.Errorf("No proceeding with empty query string")
	}
	conv_uuid, err := uuid.Parse(query_string)
	if err == nil {
		return ConversationGetByUUID(ctx, q, conv_uuid)
	}
	return ConversationGetByName(ctx, q, query_string)
}

func ConversationGetListByUUID(ctx context.Context, q *dbstore.Queries, conversation_uuids []uuid.UUID) ([]dbstore.DocketConversation, error) {
	conversations := []dbstore.DocketConversation{}
	for _, conv_uuid := range conversation_uuids {
		convo, err := ConversationGetByUUID(ctx, q, conv_uuid)
		conversations = append(conversations, convo)
		if err != nil {
			return []dbstore.DocketConversation{}, err
		}
	}
	return conversations, nil
}

func ConversationGetByUUID(ctx context.Context, q *dbstore.Queries, conversation_uuid uuid.UUID) (dbstore.DocketConversation, error) {
	conv_infos, err := q.DocketConversationRead(ctx, conversation_uuid)
	if err != nil {
		return dbstore.DocketConversation{}, err
	}
	return conv_infos, nil
}

func ConversationGetByName(ctx context.Context, q *dbstore.Queries, docketIdStr string) (dbstore.DocketConversation, error) {
	conv_infos, err := q.DocketConversationFetchByDocketIdMatch(ctx, docketIdStr)
	if err != nil {
		return dbstore.DocketConversation{}, err
	}
	if len(conv_infos) == 0 {
		errorstr := fmt.Sprintf("No proceeding with name %s", docketIdStr)
		return dbstore.DocketConversation{}, fmt.Errorf(errorstr)
	}
	return conv_infos[0], nil
}

func ConversationVerifyHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	var req conversations.ConversationInformation
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		errorstring := fmt.Sprintf("Error decoding JSON: %v\n Offending json looked like: %v", err, string(bodyBytes))
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	// log.Info(fmt.Sprintf("Unmarshaled request: %+v\n", req))

	ctx := r.Context()
	// ctx := context.Background()

	q := *database.GetTx()

	// log.Info(fmt.Sprintf("Calling verifyConversationUUID with req: %+v\n", req))
	conversation_info, err := verifyConversationUUID(ctx, q, &req, true)
	if err != nil {
		errorstring := fmt.Sprintf("Error verifying conversation %v: %v\n", req.DocketGovID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	// log.Info(fmt.Sprintf("verifyConversationUUID returned: %+v\n", conversation_info))

	// No error handling since we always want it to retun a 200 at this point.
	response, _ := json.Marshal(conversation_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func verifyConversationUUID(ctx context.Context, q dbstore.Queries, conv_info *conversations.ConversationInformation, update bool) (conversations.ConversationInformation, error) {
	// log.Info(fmt.Sprintf("Starting verifyConversationUUID with conv_info: %+v, update: %v\n", conv_info, update))

	if conv_info.ID != uuid.Nil && !update {
		log.Info("Existing UUID found and no update requested, returning early")
		return *conv_info, nil
	}

	// Try to find existing conversation for this docket
	// TODO: Change query to also match state if state exists
	results, err := q.DocketConversationFetchByDocketIdMatch(ctx, conv_info.DocketGovID)
	if err != nil {
		log.Info(fmt.Sprintf("Error fetching conversation by docket ID: %v\n", err))
		return *conv_info, err
	}

	// If conversation exists, return it
	if len(results) > 0 {
		log.Info(fmt.Sprintf("Found existing conversation with %d results\n", len(results)))
		conv := results[0]
		conv_info.ID = conv.ID
		if update {
			// log.Info("Updating existing conversation with data %v", conv_info)
			args := dbstore.DocketConversationUpdateParams{
				ID:            conv_info.ID,
				DocketGovID:   conv_info.DocketGovID,
				State:         conv_info.State,
				Name:          conv_info.Name,
				Description:   conv_info.Description,
				MatterType:    conv_info.MatterType,
				IndustryType:  conv_info.IndustryType,
				Metadata:      conv_info.Metadata,
				Extra:         conv_info.Extra,
				DatePublished: pgtype.Timestamptz{Time: time.Time(conv_info.DatePublished), Valid: true},
				// conv_info.DatePublished
			}
			_, err := q.DocketConversationUpdate(ctx, args)
			if err != nil {
				log.Info(fmt.Sprintf("Error updating conversation: %v\n", err))
				return *conv_info, err
			}
			return *conv_info, nil
		}
		conv_info.State = conv.State
		conv_info.Name = conv.Name
		conv_info.Description = conv.Description
		// log.Info("Returning existing conversation without update")
		return *conv_info, nil

	}

	// Create new conversation if none exists
	create_params := dbstore.DocketConversationCreateParams{
		DocketGovID:   conv_info.DocketGovID,
		State:         conv_info.State,
		Name:          conv_info.Name,
		Description:   conv_info.Description,
		MatterType:    conv_info.MatterType,
		IndustryType:  conv_info.IndustryType,
		Metadata:      conv_info.Metadata,
		Extra:         conv_info.Extra,
		DatePublished: pgtype.Timestamptz{Time: time.Time(conv_info.DatePublished), Valid: true},
	}

	conv_id, err := q.DocketConversationCreate(ctx, create_params)
	if err != nil {
		return *conv_info, err
	}

	conv_info.ID = conv_id
	return *conv_info, nil
}

func FileConversationUpsert(ctx context.Context, q dbstore.Queries, file_id uuid.UUID, conv_info conversations.ConversationInformation, insert bool) error {
	// Sometimes this is getting called with an insert when the metadata already exists in the table, this causes a PGERROR, since it violates uniqueness. However, setting it up so it tries to update will fall back to insert if the file doesnt exist. Its probably a good idea to remove this and debug what is causing the new file thing at some point.
	insert = false
	shouldnt_process := conv_info.ID == uuid.Nil && conv_info.DocketGovID == ""
	if shouldnt_process {
		return nil
	}
	new_conv_info, err := verifyConversationUUID(ctx, q, &conv_info, false)
	if err != nil {
		return err
	}
	if new_conv_info.ID == uuid.Nil {
		return fmt.Errorf("ASSERT FAILURE: verifyConversationUUID should never return a null uuid")
	}
	if !insert {
		args := dbstore.DocketDocumentUpdateParams{
			ConversationUuid: new_conv_info.ID,
			FileID:           file_id,
		}
		_, err = q.DocketDocumentUpdate(ctx, args)
		// If encounter a not found error, break error handling control flow and inserting object
		if err == nil {
			return nil
		}
		if err.Error() != "no rows in result set" {
			// If the error is nil, this still returns the error
			return err
		}
	}

	insert_params := dbstore.DocketDocumentInsertParams{
		ConversationUuid: new_conv_info.ID,
		FileID:           file_id,
	}
	_, err = q.DocketDocumentInsert(ctx, insert_params)
	return err
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
	q := database.GetTx()

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
