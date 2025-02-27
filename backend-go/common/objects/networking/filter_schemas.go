package networking

import (
	"encoding/json"
	"kessler/common/objects/timestamp"
	"kessler/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var log = logger.GetLogger("util")

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
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Error("failed to marshal Metadata", zap.Error(err))
		return ""
	}
	return string(jsonData)
}

func (m SearchMetadata) String() string {
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Error("failed to marshal SearchMetadata", zap.Error(err))
		return ""
	}
	return string(jsonData)
}

type MetadataFilterFields struct {
	SearchMetadata
	DateFrom timestamp.KesslerTime `json:"date_from"`
	DateTo   timestamp.KesslerTime `json:"date_to"`
}

func (f MetadataFilterFields) String() string {
	jsonData, err := json.Marshal(f)
	if err != nil {
		log.Error("failed to marshal MetadataFilterFields", zap.Error(err))
		return ""
	}
	return string(jsonData)
}

type UUIDFilterFields struct {
	AuthorUUID       uuid.UUID `json:"author_uuids,omitempty"`
	ConversationUUID uuid.UUID `json:"conversation_uuid,omitempty"`
	FileUUID         uuid.UUID `json:"file_uuid,omitempty"`
}

func (f UUIDFilterFields) String() string {
	jsonData, err := json.Marshal(f)
	if err != nil {
		log.Error("failed to marshal UUIDFilterFields", zap.Error(err))
		return ""
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

func (f FilterFields) String() string {
	jsonData, err := json.Marshal(f)
	if err != nil {
		log.Error("failed to marshal FilterFields", zap.Error(err))
		return ""
	}
	return string(jsonData)
}
