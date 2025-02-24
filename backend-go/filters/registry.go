package filters

import (
	"fmt"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
)

// FilterFunc represents a filter implementation
type FilterFunc func(input interface{}) (interface{}, error)

// FilterRegistry provides a thread-safe registry for filter implementations
type FilterRegistry struct {
	mu       sync.RWMutex
	filters  map[string]FilterFunc
	mcClient *memcache.Client
}

// NewFilterRegistry creates and initializes a new FilterRegistry
func NewFilterRegistry(memcachedServers []string) (*FilterRegistry, error) {
	if len(memcachedServers) == 0 {
		return nil, fmt.Errorf("memcached servers list cannot be empty")
	}

	return &FilterRegistry{
		filters:  make(map[string]FilterFunc),
		mcClient: memcache.New(memcachedServers...),
	}, nil
}

// Register adds a new filter implementation to the registry
func (r *FilterRegistry) Register(name string, implementation FilterFunc) error {
	if name == "" {
		return fmt.Errorf("filter name cannot be empty")
	}
	if implementation == nil {
		return fmt.Errorf("filter implementation cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.filters[name]; exists {
		return fmt.Errorf("filter %s is already registered", name)
	}

	r.filters[name] = implementation
	return nil
}

// Get retrieves a filter implementation by name
func (r *FilterRegistry) Get(name string) (FilterFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter, exists := r.filters[name]
	if !exists {
		return nil, fmt.Errorf("filter %s not found", name)
	}

	return filter, nil
}

// Execute runs a named filter with caching
func (r *FilterRegistry) Execute(name string, input interface{}) (interface{}, error) {
	filter, err := r.Get(name)
	if err != nil {
		return nil, err
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("filter:%s:%v", name, input)
	if item, err := r.mcClient.Get(cacheKey); err == nil {
		return item.Value, nil
	}

	// Execute filter if not in cache
	result, err := filter(input)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if resultBytes, ok := result.([]byte); ok {
		r.mcClient.Set(&memcache.Item{
			Key:   cacheKey,
			Value: resultBytes,
		})
	}

	return result, nil
}
