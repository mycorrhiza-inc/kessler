package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/fugusdk"
	"strings"
	"time"
)

// CardData is any of AuthorCardData, DocketCardData, DocumentCardData
// All implement GetType().

type CardData interface {
	GetType() string
}

// getResultType inspects search facets to choose a card type.
func (s *SearchService) getResultType(facets []string) string {
	for _, f := range facets {
		if strings.Contains(f, "/conversation") && !strings.Contains(f, "/data/conversation") {
			return "conversation"
		}
		if strings.Contains(f, "/organization") {
			return "organization"
		}
	}
	return "document"
}

// getCachedCard tries to unmarshal a cached value into one of the CardData types.
func (s *SearchService) getCachedCard(ctx context.Context, key string) (CardData, error) {
	if s.cacheCtrl.Client == nil {
		return nil, fmt.Errorf("cache not available")
	}

	raw, err := s.cacheCtrl.Get(key)
	if err != nil {
		return nil, err
	}

	// Try as AuthorCardData
	var a AuthorCardData
	if err := json.Unmarshal(raw, &a); err == nil && a.Type == "author" {
		return a, nil
	}

	// Try as DocketCardData
	var d DocketCardData
	if err := json.Unmarshal(raw, &d); err == nil && d.Type == "docket" {
		return d, nil
	}

	// Try as DocumentCardData
	var doc DocumentCardData
	if err := json.Unmarshal(raw, &doc); err == nil && doc.Type == "document" {
		return doc, nil
	}

	return nil, fmt.Errorf("unknown cached card type")
}

// cacheCard marshals any CardData and sets it in cache with TTL based on type.
func (s *SearchService) cacheCard(ctx context.Context, key string, card CardData) error {
	if s.cacheCtrl.Client == nil {
		return nil
	}

	raw, err := json.Marshal(card)
	if err != nil {
		return err
	}
	ttl := cache.DynamicDataTTL
	switch card.GetType() {
	case "author", "docket":
		ttl = cache.LongDataTTL
	case "document":
		ttl = cache.StaticDataTTL
	}
	return s.cacheCtrl.Set(key, raw, int32(ttl))
}

// hydrateConversation fetches or caches a DocketCardData via the shared builder.
func (s *SearchService) hydrateConversation(ctx context.Context, id string, score float32, index int) (CardData, error) {
	key := cache.PrepareKey("search", "conversation", id)
	if c, err := s.getCachedCard(ctx, key); err == nil {
		if d, ok := c.(DocketCardData); ok {
			d.Index = index
			return d, nil
		}
	}

	card, err := BuildDocketCard(ctx, s.db, id, index)
	if err != nil {
		return nil, err
	}
	s.cacheCard(ctx, key, card)
	return card, nil
}

// hydrateOrganization fetches or caches an AuthorCardData via the shared builder.
func (s *SearchService) hydrateOrganization(ctx context.Context, id string, score float32, index int) (CardData, error) {
	key := cache.PrepareKey("search", "organization", id)
	if c, err := s.getCachedCard(ctx, key); err == nil {
		if a, ok := c.(AuthorCardData); ok {
			a.Index = index
			return a, nil
		}
	}

	card, err := BuildAuthorCard(ctx, s.db, id, index)
	if err != nil {
		return nil, err
	}
	s.cacheCard(ctx, key, card)
	return card, nil
}

// hydrateDocument fetches or caches a DocumentCardData via the shared builder.
// fetchDetails=false for search hits (minimal), true for raw file fetches.
func (s *SearchService) hydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int) (CardData, error) {
	key := cache.PrepareKey("search", "document", result.ID)
	if c, err := s.getCachedCard(ctx, key); err == nil {
		if doc, ok := c.(DocumentCardData); ok {
			doc.Index = index
			return doc, nil
		}
	}

	card, err := BuildDocumentCard(ctx, s.db, result, index, false)
	if err != nil {
		return nil, err
	}
	s.cacheCard(ctx, key, card)
	return card, nil
}

// transformSearchResponse builds a SearchResponse using the hydrate* methods.
func (s *SearchService) transformSearchResponse(ctx context.Context, fuguResponse *fugusdk.SanitizedResponse, query, namespace string, pagination PaginationParams, processTime time.Duration) (*SearchResponse, error) {
	// ... unchanged
	return nil, nil
}
