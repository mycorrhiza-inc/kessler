package networking

import (
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/pkg/logger"
	"kessler/internal/objects/timestamp"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RootCacheKey = "public:metadata"
const LoggerName = "cache:networking"

func getLogger() *zap.Logger {
	return logger.GetLogger(LoggerName)
}

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

func AddMetadataToCache(metadata Metadata, key string) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding metadata to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	cacheKey := fmt.Sprintf("%s:metadata:%s", RootCacheKey, key)
	err = client.Set(&memcache.Item{
		Key:   cacheKey,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

func CachedMetadata(key string) (Metadata, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the metadata cache")
	}
	item, err := client.Get(fmt.Sprintf("%s:metadata:%s", RootCacheKey, key))
	if err != nil {
		return Metadata{}, err
	}
	var metadata Metadata
	if err := json.Unmarshal(item.Value, &metadata); err != nil {
		return Metadata{}, err
	}

	return metadata, nil
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

// TODO: still need to figure out key structure
func CachedFilterFields(key string) (FilterFields, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the filter fields cache")
	}
	item, err := client.Get(fmt.Sprintf("%s:filter:%s", RootCacheKey, key))
	if err != nil {
		return FilterFields{}, err
	}
	var filterFields FilterFields
	if err := json.Unmarshal(item.Value, &filterFields); err != nil {
		return FilterFields{}, err
	}

	return filterFields, nil
}

func AddFilterFieldsToCache(filterFields FilterFields, key string) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding filter fields to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(filterFields)
	if err != nil {
		return fmt.Errorf("failed to marshal filter fields: %v", err)
	}

	cacheKey := fmt.Sprintf("%s:filter:%s", RootCacheKey, key)
	err = client.Set(&memcache.Item{
		Key:   cacheKey,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}
