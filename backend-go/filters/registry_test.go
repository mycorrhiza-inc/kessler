package filters

import (
	"fmt"
	"sync"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func TestFilterRegistry(t *testing.T) {
	// Initialize registry with test memcached server
	client := memcache.New("localhost:11211")
	registry, err := NewFilterRegistry(client)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	t.Run("Registration validation", func(t *testing.T) {
		// Test empty name
		if err := registry.Register("", func(i interface{}) (interface{}, error) { return i, nil }); err == nil {
			t.Error("Expected error for empty name, got nil")
		}

		// Test nil implementation
		if err := registry.Register("test", nil); err == nil {
			t.Error("Expected error for nil implementation, got nil")
		}

		// Test valid registration
		if err := registry.Register("valid", func(i interface{}) (interface{}, error) { return i, nil }); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Test duplicate registration
		if err := registry.Register("valid", func(i interface{}) (interface{}, error) { return i, nil }); err == nil {
			t.Error("Expected error for duplicate registration, got nil")
		}
	})

	t.Run("Concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup
		const goroutines = 100

		// Register test filter
		testFilter := func(i interface{}) (interface{}, error) {
			return fmt.Sprintf("processed_%v", i), nil
		}
		registry.Register("concurrent_test", testFilter)

		// Concurrent execution
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Add(-1)
				result, err := registry.Execute("concurrent_test", id)
				if err != nil {
					t.Errorf("Concurrent execution failed: %v", err)
				}
				expected := fmt.Sprintf("processed_%d", id)
				if result != expected {
					t.Errorf("Expected %s, got %v", expected, result)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Type safety", func(t *testing.T) {
		registry.Register("type_test", func(i interface{}) (interface{}, error) {
			num, ok := i.(int)
			if !ok {
				return nil, fmt.Errorf("expected int, got %T", i)
			}
			return num * 2, nil
		})

		// Test with correct type
		result, err := registry.Execute("type_test", 5)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != 10 {
			t.Errorf("Expected 10, got %v", result)
		}

		// Test with incorrect type
		_, err = registry.Execute("type_test", "string")
		if err == nil {
			t.Error("Expected type assertion error, got nil")
		}
	})
}
