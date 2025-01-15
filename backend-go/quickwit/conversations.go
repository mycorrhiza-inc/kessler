package quickwit

func SearchConversations() {
	// Search for conversations
}

func IndexConversations() {
	// Index conversations
}

func ReindexConversations() {
	// Reindex conversations
}

func CreateQuickwitIndexConversations() error {
	conversationIndex := QuickwitIndex{
		Version: "0.7",
		IndexID: "NY_Conversations",
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
				{Name: "docket_id", Type: "text", Fast: true},
				{Name: "title", Type: "text", Fast: true},
				{Name: "conversation_id", Type: "boo;", Fast: true},
			},
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
