I need to refactor this FilingInfoPayload in `ingest/tasks/case_schemas.go` from its current schema into its new type. Could you go ahead and replace the filing info payload with the new schema (components defined in case_schemas.go) and change everything that depends on filing info payload and make it work with the new data structure
```go
type FillingInfoPayload struct {
	Text                  string                `json:"text"`
	FileType              string                `json:"file_type"`
	DocketID              string                `json:"docket_id"`
	PublishedDate         timestamp.RFC3339Time `json:"published_date" example:"2024-02-27T12:34:56Z"`
	Name                  string                `json:"name"`
	InternalSourceName    string                `json:"internal_source_name"`
	State                 string                `json:"state"`
	AuthorIndividual      string                `json:"author_individual"`
	AuthorIndividualEmail string                `json:"author_individual_email"`
	AuthorOrganisation    string                `json:"author_organisation"`
	FileClass             string                `json:"file_class"`
	Lang                  string                `json:"lang"`
	ItemNumber            string                `json:"item_number"`
	ExtraMetadata         map[string]any        `json:"extra_metadata"`
	Attachments           []AttachmentChildInfo `json:"attachments"`
}

type NewFillingInfoPayload struct {
	Filing   FilingChildInfo `json:"filling"`
	CaseInfo CaseInfoMinimal `json:"case_info"`
}
```
