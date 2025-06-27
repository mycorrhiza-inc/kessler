package search

import (
	"context"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	//"go.opentelemetry.io/otel"
)

func (s *SearchService) HydrateConversation(ctx context.Context, id string, score float32, index int) (CardData, error) {
	log := logger.FromContext(ctx)
	// Check cache first
	cacheKey := cache.PrepareKey("search", "conversation", id)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if docket, ok := cached.(DocketCardData); ok {
			docket.Index = index // Update index for current search
			return docket, nil
		}
	}

	// Parse UUID
	if len(id) > 36 {
		log.Warn("Raw conversation id was longer then 36 the length of a uuid", zap.String("convo_id", id))
		id = id[:36]
	}
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
