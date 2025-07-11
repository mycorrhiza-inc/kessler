// service.go
package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var serviceTracer = otel.Tracer("filter-service")

// Cache keys and TTLs
const (
	CacheKeyAllFilters       = "filters:all"
	CacheKeyNamespaceFilters = "filters:namespace:%s"
	CacheKeyFilterValues     = "filters:values:%s"
	FilterCacheTTL           = int32(300)
)

// Service handles filter operations
type Service struct {
	fuguServerURL string
	cacheCtrl     cache.CacheController
	cacheEnabled  bool
}

// NewService creates a new filter service
func NewService(fuguServerURL string) *Service {
	cacheCtrl, err := cache.NewCacheController()
	cacheEnabled := err == nil

	if !cacheEnabled {
		logger.Warn(context.Background(), "failed to initialize filter cache controller", zap.Error(err))
	}

	return &Service{
		fuguServerURL: fuguServerURL,
		cacheCtrl:     cacheCtrl,
		cacheEnabled:  cacheEnabled,
	}
}

// GetAllFilters retrieves all available filters with caching
func (s *Service) GetAllFilters(ctx context.Context) ([]fugusdk.Filter, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:get-all-filters")
	defer span.End()

	// Try cache first
	if s.cacheEnabled {
		if cached, err := s.cacheCtrl.Get(CacheKeyAllFilters); err == nil {
			var filters []fugusdk.Filter
			if err := json.Unmarshal(cached, &filters); err == nil {
				logger.Info(ctx, "all filters served from cache")
				return filters, nil
			}
		}
	}

	// Fetch from Fugu
	logger.Info(ctx, "fetching all filters from fugu backend")

	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	filters, err := client.GetAllFilters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get filters from fugu: %w", err)
	}

	// Cache the result
	if s.cacheEnabled {
		if data, err := json.Marshal(filters); err == nil {
			if err := s.cacheCtrl.Set(CacheKeyAllFilters, data, FilterCacheTTL); err != nil {
				logger.Warn(ctx, "failed to cache all filters", zap.Error(err))
			} else {
				logger.Info(ctx, "all filters cached successfully", zap.Int("filter_count", len(filters)))
			}
		}
	}

	return filters, nil
}

// GetNamespaceFilters retrieves filters for a specific namespace with caching
func (s *Service) GetNamespaceFilters(ctx context.Context, namespace string) (map[string][]string, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:get-namespace-filters")
	defer span.End()

	cacheKey := fmt.Sprintf(CacheKeyNamespaceFilters, namespace)

	// Try cache first
	if s.cacheEnabled {
		if cached, err := s.cacheCtrl.Get(cacheKey); err == nil {
			var filters map[string][]string
			if err := json.Unmarshal(cached, &filters); err == nil {
				logger.Info(ctx, "namespace filters served from cache", zap.String("namespace", namespace))
				return filters, nil
			}
		}
	}

	// Fetch from Fugu
	logger.Info(ctx, "fetching namespace filters from fugu backend", zap.String("namespace", namespace))

	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	response, err := client.GetNamespaceFilters(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace filters from fugu: %w", err)
	}

	// Extract filter paths from response
	filters := make(map[string][]string)
	if response.Data != nil {
		if data, ok := response.Data.(map[string]interface{}); ok {
			if filterPaths, ok := data["filter_paths"].(map[string]interface{}); ok {
				for path, valuesInterface := range filterPaths {
					if valuesList, ok := valuesInterface.([]interface{}); ok {
						values := make([]string, 0, len(valuesList))
						for _, v := range valuesList {
							if str, ok := v.(string); ok {
								values = append(values, str)
							}
						}
						filters[path] = values
					}
				}
			}
		}
	}

	// Cache the result
	if s.cacheEnabled {
		if data, err := json.Marshal(filters); err == nil {
			if err := s.cacheCtrl.Set(cacheKey, data, FilterCacheTTL); err != nil {
				logger.Warn(ctx, "failed to cache namespace filters", zap.Error(err))
			} else {
				logger.Info(ctx, "namespace filters cached successfully",
					zap.String("namespace", namespace),
					zap.Int("filter_count", len(filters)))
			}
		}
	}

	return filters, nil
}

// GetFilterValues retrieves values for a specific filter path with caching
func (s *Service) GetFilterValues(ctx context.Context, filterPath string) ([]string, error) {
	ctx, span := serviceTracer.Start(ctx, "filter-service:get-filter-values")
	defer span.End()

	cacheKey := fmt.Sprintf(CacheKeyFilterValues, filterPath)

	// Try cache first
	if s.cacheEnabled {
		if cached, err := s.cacheCtrl.Get(cacheKey); err == nil {
			var values []string
			if err := json.Unmarshal(cached, &values); err == nil {
				logger.Info(ctx, "filter values served from cache", zap.String("filter_path", filterPath))
				return values, nil
			}
		}
	}

	// Fetch from Fugu
	logger.Info(ctx, "fetching filter values from fugu backend", zap.String("filter_path", filterPath))

	client, err := fugusdk.NewClient(ctx, s.fuguServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	response, err := client.GetFilterValues(ctx, filterPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get filter values from fugu: %w", err)
	}

	// Extract values from response
	var values []string
	if response.Data != nil {
		if data, ok := response.Data.(map[string]interface{}); ok {
			if valuesList, ok := data["values"].([]interface{}); ok {
				values = make([]string, 0, len(valuesList))
				for _, v := range valuesList {
					if str, ok := v.(string); ok {
						values = append(values, str)
					}
				}
			}
		}
	}

	// Cache the result
	if s.cacheEnabled {
		if data, err := json.Marshal(values); err == nil {
			if err := s.cacheCtrl.Set(cacheKey, data, FilterCacheTTL); err != nil {
				logger.Warn(ctx, "failed to cache filter values", zap.Error(err))
			} else {
				logger.Info(ctx, "filter values cached successfully",
					zap.String("filter_path", filterPath),
					zap.Int("value_count", len(values)))
			}
		}
	}

	return values, nil
}

// InvalidateCache clears all filter caches
func (s *Service) InvalidateCache(ctx context.Context) {
	if !s.cacheEnabled {
		return
	}

	// For a more sophisticated cache controller, you might implement pattern-based deletion
	// For now, just delete the main cache key
	s.cacheCtrl.Delete(CacheKeyAllFilters)
	logger.Info(ctx, "filter cache invalidated")
}
