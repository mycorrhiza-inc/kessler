package search

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/internal/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

const CacheKeyRoot = "public:search"
const LoggerName = "cache:search"

func getLogger() *zap.Logger {
	return logger.GetLogger(LoggerName)
}

// CacheSearch retrieves search results from cache using the provided request key
func CacheSearch(requestKey string) ([]SearchDataHydrated, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the search cache")
		return nil, fmt.Errorf("memcached not connected")
	}

	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, requestKey))
	if err != nil {
		log.Info("cache miss", zap.String("request_key", requestKey))
		return nil, err
	}

	var searchResults []SearchDataHydrated
	if err := json.Unmarshal(item.Value, &searchResults); err != nil {
		log.Error("failed to unmarshal cached results", zap.Error(err))
		return nil, err
	}

	log.Info("cache hit", zap.String("request_key", requestKey))
	return searchResults, nil
}

// AddSearchToCache stores search results in cache with the provided request key
func AddSearchToCache(searchResults []SearchDataHydrated, requestKey string) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding search results to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(searchResults)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, requestKey)
	log.Info("adding new search query to cache", zap.String("key", key))
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}

// CacheSearchPlain retrieves plain search results from cache using the provided request key
func CacheSearchPlain(requestKey string) ([]SearchData, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected checking the plain search cache")
		return nil, fmt.Errorf("memcached not connected")
	}

	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, "plain", requestKey))
	if err != nil {
		return nil, err
	}

	var searchResults []SearchData
	if err := json.Unmarshal(item.Value, &searchResults); err != nil {
		return nil, err
	}

	return searchResults, nil
}

// AddSearchPlainToCache stores plain search results in cache with the provided request key
func AddSearchPlainToCache(searchResults []SearchData, requestKey string) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error("Memcached was not connected adding plain search results to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(searchResults)
	if err != nil {
		return fmt.Errorf("failed to marshal plain search results: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, "plain", requestKey)
	log.Info("adding new plain search query to cache", zap.String("key", key))
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}
