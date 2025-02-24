package networking

import (
	"encoding/json"
	"thaumaturgy/common/objects/timestamp"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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

type SearchMetadata struct {
	Author    string `json:"author"`
	Date      string `json:"date"`
	DocketID  string `json:"docket_id"`
	FileClass string `json:"file_class"`
	Extension string `json:"extension"`
	Lang      string `json:"lang"`
	Language  string `json:"language"`
	Source    string `json:"source"`
	Title     string `json:"title"`
}

func (m Metadata) String() string {
	// Marshal the struct to JSON format
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Info("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
}

func (m SearchMetadata) String() string {
	// Marshal the struct to JSON format
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Info("Error marshalling JSON:", err)
	}
	// Print the formatted JSON string
	return string(jsonData)
}

type MetadataFilterFields struct {
	SearchMetadata
	DateFrom timestamp.KesslerTime `json:"date_from"`
	DateTo   timestamp.KesslerTime `json:"date_to"`
}

// String method for FilterFields struct
func (f MetadataFilterFields) String() string {
	jsonData, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		log.Info("Error marshalling JSON:", err)
	}
	return string(jsonData)
}

type UUIDFilterFields struct {
	AuthorUUID       uuid.UUID `json:"author_uuids,omitempty"`
	ConversationUUID uuid.UUID `json:"conversation_uuid,omitempty"`
	FileUUID         uuid.UUID `json:"file_uuid,omitempty"`
}

func (f UUIDFilterFields) String() string {
	jsonData, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		log.Info("Error marshalling JSON:", err)
	}
	return string(jsonData)
}

func (u *UUIDFilterFields) UnmarshalJSON(data []byte) error {
	aux := &struct {
		AuthorUUID       string `json:"author_uuids,omitempty"`
		ConversationUUID string `json:"conversation_uuid,omitempty"`
		FileUUID         string `json:"file_uuid,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.AuthorUUID != "" {
		id, err := uuid.Parse(aux.AuthorUUID)
		if err != nil {
			return err
		}
		u.AuthorUUID = id
	}
	if aux.ConversationUUID != "" {
		id, err := uuid.Parse(aux.ConversationUUID)
		if err != nil {
			return err
		}
		u.ConversationUUID = id
	}
	if aux.FileUUID != "" {
		id, err := uuid.Parse(aux.FileUUID)
		if err != nil {
			return err
		}
		u.FileUUID = id
	}

	return nil
}

type FilterFields struct {
	MetadataFilters MetadataFilterFields `json:"metadata_filters"`
	UUIDFilters     UUIDFilterFields     `json:"uuid_filters"`
}
