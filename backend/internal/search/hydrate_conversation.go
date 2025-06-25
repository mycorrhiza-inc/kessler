package search

import (
	"context"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"

	"github.com/google/uuid"
	//"go.opentelemetry.io/otel"
)

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
