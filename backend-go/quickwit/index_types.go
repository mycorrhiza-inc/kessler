package quickwit

type MergePolicy struct {
	Type           string `json:"type"`
	MaxMergeOps    int    `json:"max_merge_ops"`
	MergeFactor    int    `json:"merge_factor"`
	MaxMergeFactor int    `json:"max_merge_factor"`
}

type Resources struct {
	MaxMergeWriteThroughput string `json:"max_merge_write_throughput"`
}

type IndexingSettings struct {
	MergePolicy MergePolicy `json:"merge_policy"`
	Resources   Resources   `json:"resources"`
}

type SearchSettings struct {
	DefaultSearchFields []string `json:"default_search_fields"`
}

type FieldMapping struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	InputFormats  []string `json:"input_formats,omitempty"`
	FastPrecision string   `json:"fast_precision,omitempty"`
	Fast          bool     `json:"fast"`
}

type DynamicMapping struct {
	Indexed    bool   `json:"indexed"`
	Stored     bool   `json:"stored"`
	Tokenizer  string `json:"tokenizer"`
	Record     string `json:"record"`
	ExpandDots bool   `json:"expand_dots"`
	Fast       bool   `json:"fast"`
}

type DocMapping struct {
	Mode           string         `json:"mode"`
	DynamicMapping DynamicMapping `json:"dynamic_mapping"`
	FieldMappings  []FieldMapping `json:"field_mappings"`
	TimestampField string         `json:"timestamp_field,omitempty"`
}

type Retention struct {
	Period   string `json:"period,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

type QuickwitIndex struct {
	Version          string           `json:"version"`
	IndexID          string           `json:"index_id"`
	DocMapping       DocMapping       `json:"doc_mapping,omitempty"`
	SearchSettings   SearchSettings   `json:"search_settings"`
	IndexingSettings IndexingSettings `json:"indexing_settings,omitempty"`
	Retention        Retention        `json:"retention"`
}
