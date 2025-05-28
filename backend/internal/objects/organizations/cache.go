package organizations

import (
	"encoding/json"
	"fmt"
	cache "kessler/internal/cache"
	"kessler/pkg/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

const (
	CacheKeyRoot = "public:organization"
	LoggerName   = "cache:organization"
)

func getLogger() *otelzap.Logger {
	return logger.GetLogger(LoggerName)
}

func Cached(id string) (OrganizationSchemaComplete, error) {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error(" Memcached was not connected checking the org cache")
	}
	item, err := client.Get(cache.PrepareKey(CacheKeyRoot, id))
	if err != nil {
		return OrganizationSchemaComplete{}, err
	}
	var org OrganizationSchemaComplete
	if err := json.Unmarshal(item.Value, &org); err != nil {
		return OrganizationSchemaComplete{}, err
	}

	return org, nil
}

func AddOrgToCache(org OrganizationSchemaComplete) error {
	log := getLogger()
	client := cache.MemcachedClient
	if cache.MemecachedIsConnected() != nil {
		log.Error(" Memcached was not connected adding org to cache")
		return fmt.Errorf("memcached not connected")
	}

	value, err := json.Marshal(org)
	if err != nil {
		return fmt.Errorf("failed to marshal organization: %v", err)
	}

	key := cache.PrepareKey(CacheKeyRoot, org.ID.String())
	err = client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to set cache: %v", err)
	}

	return nil
}
