package quickwit

import "github.com/google/uuid"

type QuickwitFileUploadData struct {
	Text      string                 `json:"text"`
	Name      string                 `json:"name"`
	Metadata  map[string]interface{} `json:"metadata"`
	SourceID  uuid.UUID              `json:"source_id"`
	DateFiled string                 `json:"date_filed"`
	Verified  bool                   `json:"verified"`
	Timestamp int64                  `json:"timestamp"`
}

var (
	NYPUCIndex          = "NY_PUC"
	NYConversationIndex = "NY_Conversations"
	NYOrganizationIndex = "NY_Organizations"
)

var (
	TestNYPUCIndex          = "TEST_NY_PUC"
	TestNYConversationIndex = "TEST_NY_Conversations"
	TestNYOrganizationIndex = "TEST_NY_Organizations"
)
