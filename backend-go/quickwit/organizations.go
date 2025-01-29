package quickwit

func SearchOrganizations() {
	// Search for conversations
}

func IndexOrganizations() {
	// Index conversations
}

func ReindexOrganizations(ids []string) {
	// Reindex conversations
}

func CreateQuickwitOrganizationsIndex(indexName string) error {
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
					Name: "uuid",
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
