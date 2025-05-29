package files

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/pkg/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const LoggerName = "cache:file"

const (
	CacheKeyRoot           = "public:file"
	CacheKeyCompleteFile   = "complete"
	CacheKeyFileText       = "text"
	CacheKeyFileAttachment = "attachment"
)

func getLogger() *otelzap.Logger {
	return logger.GetLogger(LoggerName)
}

func CachedFileText(fileID uuid.UUID, language string) (FileTextSchema, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the file text cache")
	}
	key := cache.PrepareKey(
		CacheKeyRoot,
		CacheKeyFileText,
		fileID.String(),
		language,
	)
	item, err := client.Get(key)
	if err != nil {
		log.Info("cache miss", zap.String("key", key))
		return FileTextSchema{}, err
	}
	var fileText FileTextSchema
	if err := json.Unmarshal(item.Value, &fileText); err != nil {
		return FileTextSchema{}, err
	}

	log.Info("cache hit", zap.String("key", key))
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
	log.Info("adding file text to cache", zap.String("key", key))
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
	key := cache.PrepareKey(CacheKeyRoot, CacheKeyFileAttachment, id.String())
	item, err := client.Get(key)
	if err != nil {
		log.Info("cache miss", zap.String("key", key))
		return CompleteAttachmentSchema{}, err
	}
	var attachment CompleteAttachmentSchema
	if err := json.Unmarshal(item.Value, &attachment); err != nil {
		return CompleteAttachmentSchema{}, err
	}

	log.Info("cache hit", zap.String("key", key))
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
	log.Info("adding attachment to cache", zap.String("key", key))
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
	key := cache.PrepareKey(CacheKeyRoot, id.String())
	item, err := client.Get(key)
	if err != nil {
		log.Info("cache miss", zap.String("key", key))
		return FileSchema{}, err
	}
	var file FileSchema
	if err := json.Unmarshal(item.Value, &file); err != nil {
		return FileSchema{}, err
	}

	log.Info("cache hit", zap.String("key", key))
	return file, nil
}

func CachedCompleteFile(id uuid.UUID) (CompleteFileSchema, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the complete file cache")
	}
	key := cache.PrepareKey(CacheKeyRoot, CacheKeyCompleteFile, id.String())
	item, err := client.Get(key)
	if err != nil {
		log.Info("cache miss", zap.String("key", key))
		return CompleteFileSchema{}, err
	}
	var file CompleteFileSchema
	if err := json.Unmarshal(item.Value, &file); err != nil {
		return CompleteFileSchema{}, err
	}

	log.Info("cache hit", zap.String("key", key))
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
	log.Info("adding file to cache", zap.String("key", key))
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
	log.Info("adding complete file to cache", zap.String("key", key))
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
