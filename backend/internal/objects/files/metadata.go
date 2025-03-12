package files

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Metadata struct {
	ItemNumber       string      `json:"item_number"`
	Author           string      `json:"author"`
	Date             string      `json:"date"`
	DocketID         string      `json:"docket_id"`
	FileClass        string      `json:"file_class"`
	Extension        string      `json:"extension"`
	Lang             string      `json:"lang"`
	Language         string      `json:"language"`
	Source           string      `json:"source"`
	Title            string      `json:"title"`
	ConversationUUID uuid.UUID   `json:"conversation_uuid"`
	Authors          []string    `json:"authors"`
	AuthorUUIDs      []uuid.UUID `json:"author_uuids"`
}

func (m Metadata) String() string {
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Error("failed to marshal Metadata", zap.Error(err))
		return ""
	}
	return string(jsonData)
}
