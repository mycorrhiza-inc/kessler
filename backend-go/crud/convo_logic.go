package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/common/objects/conversations"
	"kessler/db"
	"kessler/gen/dbstore"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

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

	q := *db.GetTx()

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
			q.DocketConversationUpdate(ctx, args)
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

func fileConversationUpsert(ctx context.Context, q dbstore.Queries, file_id uuid.UUID, conv_info conversations.ConversationInformation, insert bool) error {
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
