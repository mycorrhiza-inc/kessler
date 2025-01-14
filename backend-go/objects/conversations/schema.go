package conversations

import (
	"kessler/objects/timestamp"

	"github.com/google/uuid"
)

type ConversationInformation struct {
	DocketGovID   string
	State         string
	Name          string
	Description   string
	MatterType    string
	IndustryType  string
	Metadata      []byte
	Extra         []byte
	DatePublished timestamp.KesslerTime
	ID            uuid.UUID
}
