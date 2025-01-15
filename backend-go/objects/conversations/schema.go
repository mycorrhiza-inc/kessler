package conversations

import (
	"kessler/objects/timestamp"

	"github.com/google/uuid"
)

type ConversationInformation struct {
	DocketGovID   string                `json:"docket_gov_id"`
	State         string                `json:"state"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	MatterType    string                `json:"matter_type"`
	IndustryType  string                `json:"industry_type"`
	Metadata      []byte                `json:"metadata"`
	Extra         []byte                `json:"extra"`
	DatePublished timestamp.KesslerTime `json:"date_published"`
	ID            uuid.UUID             `json:"id"`
}
