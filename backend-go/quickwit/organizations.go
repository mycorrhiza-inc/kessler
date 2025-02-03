package quickwit

import (
	"context"
	"kessler/gen/dbstore"

	"github.com/google/uuid"
)

type OrganizationQuickwitSchema struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Aliases            []string  `json:"aliases"`
	FilesAuthoredCount int64     `json:"files_authored_count"`
}

func SearchOrganizations() {
	// Search for conversations
}

func IndexOrganizations(org_schemas OrganizationQuickwitSchema, ctx context.Context, index_name string) error {
	// Index conversations
}

func ReindexAllOrganizations(ctx context.Context, q dbstore.Queries, index_name string) error {
	// Reindex conversations
}

func CreateQuickwitOrganizationsIndex(indexName string) error {
	if indexName == "" {
		indexName = NYOrganizationIndex
	}
	requestData := QuickwitIndex{
		Version: "0.7",
		IndexID: indexName,
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
				{
					Name: "name",
					Type: "text",
					Fast: true,
				},
				{
					Name: "aliases",
					Type: "array<text>",
					Fast: true,
				},
				{
					Name: "id",
					Type: "text",
					Fast: true,
				},
				{
					Name:          "timestamp",
					Type:          "datetime",
					Fast:          true,
					InputFormats:  []string{"unix_timestamp"},
					FastPrecision: "seconds",
				},
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
	err := CreateIndex(requestData)
	return err
}
