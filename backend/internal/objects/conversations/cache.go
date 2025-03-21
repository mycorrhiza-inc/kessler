package conversations

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/internal/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const CacheKeyRoot = "public:conversation"
const LoggerName = "cache:conversation"

func getLogger() *zap.Logger {
	return logger.GetLogger(LoggerName)
}

func Cached(id uuid.UUID) (ConversationInformation, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the conversation cache")
	}
	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, id.String()))
	if err != nil {
		return ConversationInformation{}, err
	}
	var conversation ConversationInformation
	if err := json.Unmarshal(item.Value, &conversation); err != nil {
		return ConversationInformation{}, err
	}

	return conversation, nil
}

func AddConversationToCache(conversation ConversationInformation) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding conversation to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, conversation.ID.String())
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}
