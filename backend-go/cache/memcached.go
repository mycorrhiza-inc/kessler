package cache

import (
	"fmt"
	"os"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

var MemcachedClient *memcache.Client

// InitMemcached initializes the global memcached client using environment variables.
// It expects MEMCACHED_SERVERS to be a comma-separated list of host:port pairs.
// Example: "localhost:11211,otherhost:11211"
// If MEMCACHED_SERVERS is not set, it defaults to "localhost:11211"
func InitMemcached() error {
	servers := os.Getenv("MEMCACHED_SERVERS")
	if servers == "" {
		servers = "localhost:11211" // Default server
	}

	// Split the servers string into individual server addresses
	serverList := strings.Split(servers, ",")

	// Trim any whitespace from server addresses
	for i, server := range serverList {
		serverList[i] = strings.TrimSpace(server)
	}

	if len(serverList) == 0 {
		return fmt.Errorf("no memcached servers specified")
	}

	// Initialize the global client
	MemcachedClient = memcache.New(serverList...)

	// Test the connection
	err := MemcachedClient.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to memcached: %w", err)
	}

	return nil
}
