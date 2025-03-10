package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/internal/database"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
