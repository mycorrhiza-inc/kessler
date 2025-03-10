package filters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"kessler/internal/database"
	"kessler/internal/dbstore"
	"kessler/internal/logger"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Common errors that can be returned by the FilterService
var (
	ErrFilterNotFound     = fmt.Errorf("filter not found")
	ErrInvalidFilterState = fmt.Errorf("invalid filter state")
	ErrDatabaseOperation  = fmt.Errorf("database operation failed")
)

// FilterService handles business logic for filters
type FilterService struct {
	db          *pgxpool.Pool
	registry    *FilterRegistry
	queryEngine *dbstore.Queries // This will be generated by sqlc
	logger      *zap.Logger
}

var log = logger.GetLogger("main")

// NewFilterService creates a new instance of FilterService
func NewFilterService(db *pgxpool.Pool, cache *memcache.Client) *FilterService {
	registry, err := NewFilterRegistry(cache)
	if err != nil {
		log.Fatal("unable to load filter field", zap.Error(err))
	}
	filter_log := logger.GetLogger("filter_service")

	qe := database.GetTx()

	return &FilterService{
		db:          db,
		registry:    registry,
		queryEngine: qe,
		logger:      filter_log,
	}
}

func (s *FilterService) TempInitializeFilters() error {
	s.logger.Info("initializing dummy filter data")
	return nil
}

func (s *FilterService) RegisterDefaultFilters() error {
	// daterange filter
	err := s.registry.Register("daterange", NewDateRangeFilter(s.logger))
	if err != nil {
		s.logger.Error("failed to register date range filter",
			zap.Error(err))
		return fmt.Errorf("failed to register date range filter: %w", err)
	}

	err = s.registry.Register("text", NewTextFilter(s.logger))
	if err != nil {
		s.logger.Error("failed to register text filter",
			zap.Error(err))
		return fmt.Errorf("failed to register text filter: %w", err)
	}

	return nil
}

func (s *FilterService) AddFilterToState(ctx context.Context, stateAbbrev string) error {
	state := strings.ToUpper(stateAbbrev)
	cacheKey := fmt.Sprintf("filters:state:%s", state)
	log.Debug(fmt.Sprintf("adding filter to %s", stateAbbrev), zap.String("filter key", cacheKey))

	s.registry.mcClient.Set(&memcache.Item{
		Key:   cacheKey,
		Value: []byte{},
	})

	return nil
}

// GetFiltersByState retrieves filters by their state
func (s *FilterService) GetFiltersByState(ctx context.Context, stateAbbrev string) ([]dbstore.Filter, error) {
	if stateAbbrev == "" {
		return nil, fmt.Errorf("%w: empty state", ErrInvalidFilterState)
	}
	state := strings.ToUpper(stateAbbrev)

	// Try cache first
	cacheKey := fmt.Sprintf("filters:state:%s", state)
	if item, err := s.registry.mcClient.Get(cacheKey); err == nil {
		var filters []dbstore.Filter
		if err := json.Unmarshal(item.Value, &filters); err == nil {
			s.logger.Debug("cache hit for filters", zap.String("state", state))
			return filters, nil
		}
	}

	// If cache-miss, preform the query
	filters, err := s.queryEngine.GetFiltersByState(ctx, state)
	if err != nil {
		if err == sql.ErrNoRows {
			return []dbstore.Filter{}, nil
		}
		s.logger.Error("failed to fetch filters",
			zap.String("state", state),
			zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	// Cache the results
	if cacheData, err := json.Marshal(filters); err == nil {
		s.registry.mcClient.Set(&memcache.Item{
			Key:        cacheKey,
			Value:      cacheData,
			Expiration: 300, // 5 minute TTL
		})
	}

	return filters, nil
}
