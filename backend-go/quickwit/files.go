package quickwit

func CreateQuickwitNYFileIndex(index_name string) error {
	if index_name == "" {
		index_name = NYPUCIndex
	}
	conversationIndex := QuickwitIndex{
		Version: "0.7",
		IndexID: index_name,
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
				{Name: "text", Type: "text", Fast: true},
				{Name: "name", Type: "text", Fast: true},
				{Name: "metadata", Type: "json", Fast: true},
				{Name: "source_id", Type: "text", Fast: true},
				{Name: "date_filed", Type: "text", Fast: true},
				{Name: "verified", Type: "bool", Fast: true},
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
			DefaultSearchFields: []string{"text", "name"},
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
