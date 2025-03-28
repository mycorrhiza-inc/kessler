package conversations

import (
	"kessler/internal/objects/timestamp"

	"github.com/google/uuid"
)

type ConversationInformation struct {
	DocketGovID    string                `json:"docket_gov_id"`
	State          string                `json:"state"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	MatterType     string                `json:"matter_type"`
	IndustryType   string                `json:"industry_type"`
	Metadata       string                `json:"metadata"`
	Extra          string                `json:"extra"`
	DocumentsCount int                   `json:"documents_count"`
	DatePublished  timestamp.KesslerTime `json:"date_published"`
	ID             uuid.UUID             `json:"id"`
}
