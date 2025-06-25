package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"time"

	"github.com/google/uuid"
	//"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// Updated search response to use card data

// NewSearchService creates a new search service with database access

// Helper to determine result type from facets

// Get or create cached card data

// Cache card data
func (s *SearchService) cacheCard(ctx context.Context, key string, card CardData) error {
	if s.cacheCtrl.Client == nil {
		return nil // Skip if cache not available
	}

	data, err := json.Marshal(card)
	if err != nil {
		return err
	}

	// Use appropriate TTL based on card type
	ttl := cache.DynamicDataTTL
	switch card.GetType() {
	case "author", "docket":
		ttl = cache.LongDataTTL
	case "document":
		ttl = cache.StaticDataTTL
	}

	return s.cacheCtrl.Set(key, data, int32(ttl))
}

// Hydrate conversation/docket data
func (s *SearchService) HydrateConversation(ctx context.Context, id string, score float32, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "conversation", id)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if docket, ok := cached.(DocketCardData); ok {
			docket.Index = index // Update index for current search
			return docket, nil
		}
	}

	// Parse UUID
	conversationID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Query database
	queries := dbstore.New(s.db)
	conversation, err := queries.DocketConversationRead(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to read conversation: %w", err)
	}

	// Create card data
	card := DocketCardData{
		Name:        conversation.Name,
		Description: conversation.Description,
		Timestamp:   conversation.CreatedAt.Time,
		Index:       index,
		Type:        "docket",
		ObjectUUID:  conversation.ID,
	}

	// Cache the result
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}

// Hydrate organization/author data
func (s *SearchService) HydrateOrganization(ctx context.Context, id string, score float32, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "organization", id)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if author, ok := cached.(AuthorCardData); ok {
			author.Index = index // Update index for current search
			return author, nil
		}
	}

	// Parse UUID
	orgID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	// Query database
	queries := dbstore.New(s.db)
	org, err := queries.OrganizationRead(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to read organization: %w", err)
	}

	// Create card data
	extraInfo := ""
	if org.IsPerson.Valid && org.IsPerson.Bool {
		extraInfo = "Individual contributor"
	} else {
		extraInfo = "Organization"
	}

	card := AuthorCardData{
		Name:        org.Name,
		Description: org.Description,
		Timestamp:   org.CreatedAt.Time,
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "author",
		ObjectUUID:  org.ID,
	}

	// Cache the result
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}

// Hydrate document data

// Updated transformSearchResponse to return card data
func (s *SearchService) transformSearchResponse(ctx context.Context, fuguResponse *fugusdk.SanitizedResponse, query, namespace string, pagination PaginationParams, processTime time.Duration) (*SearchResponse, error) {
	if fuguResponse == nil || len(fuguResponse.Results) == 0 {
		return &SearchResponse{
			Data:        []CardData{},
			Total:       0,
			Page:        pagination.Page,
			PerPage:     pagination.Limit,
			Query:       query,
			Namespace:   namespace,
			ProcessTime: processTime.String(),
		}, nil
	}

	var cards []CardData

	for i, result := range fuguResponse.Results {
		resultType := s.getResultType(result.Facets)

		var card CardData
		var err error

		switch resultType {
		case "conversation":
			card, err = s.HydrateConversation(ctx, result.ID, result.Score, i)
		case "organization":
			card, err = s.HydrateOrganization(ctx, result.ID, result.Score, i)
		default:
			card, err = s.HydrateDocument(ctx, result, i, false)
		}

		if err != nil {
			logger.Warn(ctx, "failed to hydrate result",
				zap.String("id", result.ID),
				zap.String("type", resultType),
				zap.Error(err))
			continue // Skip this result
		}

		cards = append(cards, card)
	}

	return &SearchResponse{
		Data:        cards,
		Total:       fuguResponse.Total,
		Page:        pagination.Page,
		PerPage:     pagination.Limit,
		Query:       query,
		Namespace:   namespace,
		ProcessTime: processTime.String(),
	}, nil
}

// Update ProcessSearch to use the new transformer
