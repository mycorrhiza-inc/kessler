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
func InitMemcached(servers ...string) error {
	// check for environment servers
	serverList := servers
	env_servers := os.Getenv("MEMCACHED_SERVERS")

	if env_servers == "" && len(servers) == 0 {
		env_servers = "localhost:11211" // Default server
	}

	// Split the servers string into individual server addresses
	envlist := strings.Split(env_servers, ",")

	// Trim any whitespace from server addresses
	for _, server := range envlist {
		serverList = append(serverList, strings.TrimSpace(server))
	}

	if len(servers) == 0 {
		return fmt.Errorf("no memcached servers specified")
	}

	serverListString := strings.Join(serverList, ",")

	// Initialize the global client
	MemcachedClient = memcache.New(serverListString)

	// Test the connection
	err := MemcachedClient.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to memcached: %w", err)
	}

	return nil
}
