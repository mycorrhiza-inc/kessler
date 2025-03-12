package authors

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/internal/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const CacheKeyRoot = "public:author"
const LoggerName = "cache:author"

func getLogger() *zap.Logger {
	return logger.GetLogger(LoggerName)
}

func Cached(id uuid.UUID) (AuthorInformation, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the author cache")
	}
	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, id.String()))
	if err != nil {
		return AuthorInformation{}, err
	}
	var author AuthorInformation
	if err := json.Unmarshal(item.Value, &author); err != nil {
		return AuthorInformation{}, err
	}

	return author, nil
}

func AddAuthorToCache(author AuthorInformation) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding author to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(author)
	if err != nil {
		return fmt.Errorf("failed to marshal author: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, author.AuthorID.String())
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}
