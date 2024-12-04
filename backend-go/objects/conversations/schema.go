package conversations

import "github.com/google/uuid"

type ConversationInformation struct {
	ID          uuid.UUID `json:"id"`
	DocketID    string    `json:"docket_id"`
	State       string    `json:"state"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
