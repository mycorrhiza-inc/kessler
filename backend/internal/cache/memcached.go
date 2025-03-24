package cache

import (
	"fmt"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

// TTL constants for different types of data
const (
	// LongData
	LongDataTTL = 432000 // 5 days

	// Static data can be cached longer
	StaticDataTTL = 3600 // 1 hour

	// Dynamic data should have shorter TTL
	DynamicDataTTL = 300 // 5 minutes

	// List data might change frequently
	ListDataTTL = 60 // 1 minute
)

var MemcachedClient *memcache.Client

// InitMemcached initializes the global memcached client using environment variables.
// It expects MEMCACHED_SERVERS to be a comma-separated list of host:port pairs.
// Example: "localhost:11211,otherhost:11211"
// If MEMCACHED_SERVERS is not set, it defaults to "localhost:11211"
func InitMemcached(servers ...string) error {
	// // check for environment servers
	// serverList := servers
	// env_servers := os.Getenv("MEMCACHED_SERVERS")

	// if env_servers == "" && len(servers) == 0 {
	// 	env_servers = "localhost:11211" // Default server
	// }

	// // Split the servers string into individual server addresses
	// envlist := strings.Split(env_servers, ",")

	// // Trim any whitespace from server addresses
	// for _, server := range envlist {
	// 	serverList = append(serverList, strings.TrimSpace(server))
	// }

	// if len(servers) == 0 {
	// 	return fmt.Errorf("no memcached servers specified")
	// }

	// serverListString := strings.Join(serverList, ",")

	// Initialize the global client
	MemcachedClient = memcache.New("cache:11211")

	// Test the connection
	err := MemecachedIsConnected()
	if err != nil {
		fmt.Println("MEMEMEMMEMEMEMMEM")
	}
	return err
}

func MemecachedIsConnected() error {
	err := MemcachedClient.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to memcached: %w", err)
	}
	return nil
}

// TODO: rebuild this in the /pkg directory to be able to have a more portable tool
// @bradfitz abandoned this library years ago and there needs to be some maintenance on it

type CacheController struct {
	Client *memcache.Client
}

func PrepareKey(root string, args ...string) string {
	// TODO: Add key validation to ensure no invalid characters in keys
	parts := append([]string{root}, args...)
	return strings.Join(parts, ":")
}

func NewCacheController() (CacheController, error) {
	err := MemecachedIsConnected()
	if err != nil {
		return CacheController{}, err
	}

	return CacheController{
		Client: MemcachedClient,
	}, nil
}

func (cc CacheController) PatternListKeys(pattern string) []string {
	// Note: Memcached does not natively support pattern matching for keys
	// This is a placeholder that returns an empty slice since pattern-based
	// key listing is not available with basic memcached
	keys := []string{}
	return keys
}

// Set stores a key-value pair in memcached with optional expiration
func (cc CacheController) Set(key string, value []byte, expiration int32) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
	return cc.Client.Set(item)
}

// Get retrieves a value from memcached by key
func (cc CacheController) Get(key string) ([]byte, error) {
	item, err := cc.Client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

// Touch updates the expiration time for a key
func (cc CacheController) Touch(key string, expiration int32) error {
	return cc.Client.Touch(key, expiration)
}
