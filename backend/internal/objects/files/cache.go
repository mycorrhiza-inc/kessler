package files

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/internal/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const LoggerName = "cache:file"

const CacheKeyRoot = "public:file"
const CacheKeyCompleteFile = "complete"
const CacheKeyFileText = "text"
const CacheKeyFileAttachment = "attachment"

func getLogger() *zap.Logger {
	return logger.GetLogger(LoggerName)
}

func CachedFileText(fileID uuid.UUID, language string) (FileTextSchema, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the file text cache")
	}
	item, err := client.Get(cache.PrepareKey(
		CacheKeyRoot,
		CacheKeyFileText,
		fileID.String(),
		language,
	))
	if err != nil {
		return FileTextSchema{}, err
	}
	var fileText FileTextSchema
	if err := json.Unmarshal(item.Value, &fileText); err != nil {
		return FileTextSchema{}, err
	}

	return fileText, nil
}

func AddFileTextToCache(fileText FileTextSchema) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding file text to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(fileText)
	if err != nil {
		return fmt.Errorf("failed to marshal file text: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, CacheKeyFileText, fileText.FileID.String(), fileText.Language)
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

func CachedAttachment(id uuid.UUID) (CompleteAttachmentSchema, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the attachment cache")
	}
	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, CacheKeyFileAttachment, id.String()))
	if err != nil {
		return CompleteAttachmentSchema{}, err
	}
	var attachment CompleteAttachmentSchema
	if err := json.Unmarshal(item.Value, &attachment); err != nil {
		return CompleteAttachmentSchema{}, err
	}

	return attachment, nil
}

func AddAttachmentToCache(attachment CompleteAttachmentSchema) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding attachment to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(attachment)
	if err != nil {
		return fmt.Errorf("failed to marshal attachment: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, CacheKeyFileAttachment, attachment.ID.String())
	err = client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: cache.StaticDataTTL,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

func CachedFile(id uuid.UUID) (FileSchema, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the file cache")
	}
	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, id.String()))
	if err != nil {
		return FileSchema{}, err
	}
	var file FileSchema
	if err := json.Unmarshal(item.Value, &file); err != nil {
		return FileSchema{}, err
	}

	return file, nil
}

func CachedCompleteFile(id uuid.UUID) (CompleteFileSchema, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the complete file cache")
	}
	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, CacheKeyCompleteFile, id.String()))
	if err != nil {
		return CompleteFileSchema{}, err
	}
	var file CompleteFileSchema
	if err := json.Unmarshal(item.Value, &file); err != nil {
		return CompleteFileSchema{}, err
	}

	return file, nil
}

func AddFileToCache(file FileSchema) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding file to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("failed to marshal file: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, file.ID.String())
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

func AddCompleteFileToCache(file CompleteFileSchema) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding complete file to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("failed to marshal complete file: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, CacheKeyCompleteFile, file.ID.String())
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	// Also cache the basic file schema
	basicFile := file.CompleteFileSchemaPrune()
	return AddFileToCache(basicFile)
}
