package conversations

import (
	"time"

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
	DatePublished time.Time
	ID            uuid.UUID
}
