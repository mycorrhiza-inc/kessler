package quickwit

import (
	"context"
	"fmt"
	"kessler/crud"
	"kessler/gen/dbstore"
	"kessler/objects/conversations"
	"kessler/objects/timestamp"
	"kessler/util"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ConversationSearchSchema struct {
	Query        string `json:"query"`
	IndustryType string `json:"industry_type"`
}

func SearchConversations(search_schema ConversationSearchSchema, ctx context.Context) ([]conversations.ConversationInformation, error) {
	// Search for conversations
}

func IndexAllConversations(q dbstore.Queries, ctx context.Context) error {
	conversations, err := q.DocketConversationList(ctx)
	if err != nil {
		return err
	}
	err = IndexConversations(conversations)
	return err
}

func IndexConversations(convos []dbstore.DocketConversation) error {
	// Index conversations
	quickwit_convos := make([]conversations.ConversationInformation, len(convos))
	for index, convo := range convos {
		quickwit_convo := conversations.ConversationInformation{
			DocketGovID:   convo.DocketGovID,
			State:         convo.State,
			Name:          convo.Name,
			Description:   convo.Description,
			MatterType:    convo.MatterType,
			IndustryType:  convo.IndustryType,
			Metadata:      convo.Metadata,
			Extra:         convo.Extra,
			DatePublished: timestamp.KesslerTime(convo.DatePublished.Time),
			ID:            convo.ID,
		}
		quickwit_convos[index] = quickwit_convo

	}
	IngestIntoIndex(NYConversationIndex, quickwit_convos)
	return nil
}

func IndexConversationsFromUUIDs(ids []uuid.UUID, ctx context.Context) error {
	//
	// Index conversations from UUIDs
	q := *util.DBQueriesFromContext(ctx)
	conversations, err := crud.ConversationGetListByUUID(ctx, &q, ids)
	if err != nil {
		return err
	}
	IndexConversations(conversations)
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
		fmt.Println(err)
	}
	return nil
}

func ReindexConversationsFromUUID(conversationUUIDs []uuid.UUID, ctx context.Context) error {
	// Reindex conversations
	err := DeleteConversationsFromIndex(conversationUUIDs, "NY_Conversations")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = IndexConversationsFromUUIDs(conversationUUIDs, ctx)
	if err != nil {
		fmt.Println(err)
		return err
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
				{Name: "uuid", Type: "text", Fast: true},
				{Name: "state", Type: "text", Fast: true},
				{Name: "docket_id", Type: "text", Fast: true},
				{Name: "title", Type: "text", Fast: true},
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
