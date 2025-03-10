package quickwit

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/networking"
	"kessler/internal/objects/timestamp"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/google/uuid"
)

type ConvoSearchRequestData struct {
	Search ConversationSearchSchema `json:"search"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}

type ConversationSearchSchema struct {
	Query        string                `json:"query"`
	IndustryType string                `json:"industry_type"`
	DateFrom     timestamp.KesslerTime `json:"date_from"`
	DateTo       timestamp.KesslerTime `json:"date_to"`
}

func HandleConvoSearch(w http.ResponseWriter, r *http.Request) {
	log.Info("Received a search request")

	// Create an instance of RequestData
	var convo_search ConversationSearchSchema

	// Decode the JSON body into the struct
	err := json.NewDecoder(r.Body).Decode(&convo_search)
	if err != nil {
		json_bytes := []byte{}
		r.Body.Read(json_bytes)
		errorstring := fmt.Sprintf("Error decoding JSON: %v\nWith this error: %v", json_bytes, err)
		log.Info(fmt.Sprintf(errorstring))
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body when done

	pagination := networking.PaginationFromUrlParams(r)
	convoRequestData := ConvoSearchRequestData{
		Search: convo_search,
		Limit:  int(pagination.Limit),
		Offset: int(pagination.Offset),
	}
	results, err := SearchConversations(convoRequestData, r.Context())
	if err != nil {
		errorstring := fmt.Sprintf("Error searching conversations: %v\n", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
	}

	respString, err := json.Marshal(results)
	if err != nil {
		log.Info("Error marshalling response data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(respString))
}

func convoQuickwitMakeQueryString(search_data ConvoSearchRequestData) string {
	search_values := search_data.Search
	dateQueryString := ConstructDateTextQuery(search_values.DateFrom, search_values.DateTo, search_values.Query)
	if search_values.IndustryType != "" {
		dateQueryString = dateQueryString + fmt.Sprintf(" AND industry_type:(%s)", search_values.IndustryType)
	}
	return dateQueryString
}

func SearchConversations(search_data ConvoSearchRequestData, ctx context.Context) ([]conversations.ConversationInformation, error) {
	queryString := convoQuickwitMakeQueryString(search_data)

	// Search for conversations
	search_request := QuickwitSearchRequest{
		Query:       queryString,
		MaxHits:     search_data.Limit,
		StartOffset: search_data.Offset,
		SortBy:      "documents_count",
	}
	search_index := NYConversationIndex
	var search_results []conversations.ConversationInformation

	err := SearchHitsQuickwitGeneric(&search_results, search_request, search_index)
	return search_results, err
}

func IndexAllConversations(q dbstore.Queries, ctx context.Context, index_name string) error {
	conversations, err := q.ConversationCompleteQuickwitListGet(ctx)
	if err != nil {
		return err
	}
	err = IndexConversations(conversations, index_name)
	return err
}

func IndexConversations(convos []dbstore.ConversationCompleteQuickwitListGetRow, index_name string) error {
	if index_name == "" {
		index_name = NYConversationIndex
	}
	// Index conversations
	quickwit_convos := make([]conversations.ConversationInformation, len(convos))
	for index, convo := range convos {
		quickwit_convo := conversations.ConversationInformation{
			DocketGovID:    convo.DocketGovID,
			State:          convo.State,
			Name:           convo.Name,
			Description:    convo.Description,
			MatterType:     convo.MatterType,
			IndustryType:   convo.IndustryType,
			Metadata:       convo.Metadata,
			Extra:          convo.Extra,
			DatePublished:  timestamp.KesslerTime(convo.DatePublished.Time),
			ID:             convo.ID,
			DocumentsCount: int(convo.TotalDocuments),
		}
		quickwit_convos[index] = quickwit_convo

	}
	IngestIntoIndex(index_name, quickwit_convos, true)
	return nil
}

func DeleteConversationsFromIndex(conversationUUIDs []uuid.UUID, index string) error {
	// Format the current time as an RFC3339 timestamp
	now := time.Now().Format(time.RFC3339)
	tempIDS := fmt.Sprintf("%s", uuid.UUIDs(conversationUUIDs).Strings())
	idlist := strings.ReplaceAll(tempIDS, ",", "")
	query := fmt.Sprintf(`uuid IN %s`, idlist)
	task := DeleteTask{
		Query:        query,
		EndTimestamp: now,
	}
	err := CreateDeleteTask(index, task)
	if err != nil {
		log.Info(err)
	}
	return nil
}

func CreateQuickwitIndexConversations() error {
	conversationIndex := QuickwitIndex{
		Version: "0.7",
		IndexID: NYConversationIndex,
		DocMapping: DocMapping{
			Mode: "dynamic",
			DynamicMapping: DynamicMapping{
				Indexed:    true,
				Stored:     true,
				Tokenizer:  "default",
				Record:     "basic",
				ExpandDots: true,
				Fast:       true,
			},
			FieldMappings: []FieldMapping{
				{Name: "id", Type: "text", Fast: true},
				{Name: "state", Type: "text", Fast: true},
				{Name: "docket_id", Type: "text", Fast: true},
				{Name: "title", Type: "text", Fast: true},
				{
					Name:          "date_published",
					Type:          "datetime",
					Fast:          true,
					InputFormats:  []string{"rfc3339"},
					FastPrecision: "seconds",
				},
				{
					Name:          "timestamp",
					Type:          "datetime",
					Fast:          true,
					InputFormats:  []string{"unix_timestamp"},
					FastPrecision: "seconds",
				},

				// {Name: "conversation_id", Type: "text", Fast: true},
			},
			TimestampField: "timestamp",
		},
		SearchSettings: SearchSettings{
			DefaultSearchFields: []string{"name"},
		},
		IndexingSettings: IndexingSettings{
			MergePolicy: MergePolicy{
				Type:           "limit_merge",
				MaxMergeOps:    3,
				MergeFactor:    10,
				MaxMergeFactor: 12,
			},
			Resources: Resources{
				MaxMergeWriteThroughput: "80mb",
			},
		},
		Retention: Retention{
			Period:   "10 years",
			Schedule: "yearly",
		},
	}
	err := CreateIndex(conversationIndex)
	return err
}

// func ReindexConversationsFromUUID(conversationUUIDs []uuid.UUID, ctx context.Context) error {
// 	// Reindex conversations
// 	err := DeleteConversationsFromIndex(conversationUUIDs, "NY_Conversations")
// 	if err != nil {
// 		log.Info(err)
// 		return err
// 	}
// 	err = IndexConversationsFromUUIDs(conversationUUIDs, ctx)
// 	if err != nil {
// 		log.Info(err)
// 		return err
// 	}
// 	return nil
// }
//
// func IndexConversationsFromUUIDs(ids []uuid.UUID, ctx context.Context) error {
// 	//
// 	// Index conversations from UUIDs
// 	q := *util.DBQueriesFromContext(ctx)
// 	conversations, err := crud.ConversationGetListByUUID(ctx, &q, ids)
// 	if err != nil {
// 		return err
// 	}
// 	IndexConversations(conversations)
// 	return nil
// }
