package juristictions

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/internal/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

const RootCacheKey = "public:jurisdiction"
const LoggerName = "cache:jurisdiction"

func GetLogger() *zap.Logger {
	return logger.GetLogger(LoggerName)
}

func Cached(key string) (JuristictionInformation, error) {
	log := GetLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the jurisdiction cache")
	}
	item, err := client.Get(fmt.Sprintf("%s:%s", RootCacheKey, key))
	if err != nil {
		return JuristictionInformation{}, err
	}
	var jurisdiction JuristictionInformation
	if err := json.Unmarshal(item.Value, &jurisdiction); err != nil {
		return JuristictionInformation{}, err
	}

	return jurisdiction, nil
}

func AddJurisdictionToCache(jurisdiction JuristictionInformation, key string) error {
	log := GetLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding jurisdiction to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(jurisdiction)
	if err != nil {
		return fmt.Errorf("failed to marshal jurisdiction: %v", err)
	}

	cacheKey := fmt.Sprintf("%s:%s", RootCacheKey, key)
	err = client.Set(&memcache.Item{
		Key:   cacheKey,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}
