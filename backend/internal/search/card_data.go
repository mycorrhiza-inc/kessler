package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"strings"
	"time"

	"github.com/google/uuid"
	//"go.opentelemetry.io/otel"
)

// Card data types matching the frontend requirements
type CardData interface {
	GetType() string
}

type AuthorCardData struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	ExtraInfo   string    `json:"extraInfo,omitempty"`
	Index       int       `json:"index"`
	Type        string    `json:"type"`
	ObjectUUID  uuid.UUID `json:"object_uuid"`
}

func (a AuthorCardData) GetType() string {
	return "author"
}

type DocketCardData struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Timestamp    time.Time `json:"timestamp"`
	Index        int       `json:"index"`
	Type         string    `json:"type"`
	ObjectUUID   uuid.UUID `json:"object_uuid"`
	DocketNumber string    `json:"docket_number"`
}

func (d DocketCardData) GetType() string {
	return "docket"
}

type DocumentAuthor struct {
	AuthorName      string    `json:"author_name"`
	IsPerson        bool      `json:"is_person"`
	IsPrimaryAuthor bool      `json:"is_primary_author"`
	AuthorID        uuid.UUID `json:"author_id"`
}

type DocumentConversation struct {
	ConvoName   string    `json:"convo_name"`
	ConvoNumber string    `json:"convo_number"`
	ConvoID     uuid.UUID `json:"convo_id"`
}

type DocumentCardData struct {
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Timestamp      time.Time            `json:"timestamp"`
	ExtraInfo      string               `json:"extraInfo,omitempty"`
	Index          int                  `json:"index"`
	Type           string               `json:"type"`
	ObjectUUID     uuid.UUID            `json:"object_uuid"`
	FileUUID       uuid.UUID            `json:"file_uuid"`
	AttachmentUUID uuid.UUID            `json:"attachment_uuid"`
	FragmentID     string               `json:"fragment_id"`
	Authors        []DocumentAuthor     `json:"authors"`
	Conversation   DocumentConversation `json:"conversation"`
}

func (d DocumentCardData) GetType() string {
	return "document"
}

// Updated search response to use card data

// NewSearchService creates a new search service with database access

// Helper to determine result type from facets
func (s *SearchService) getResultType(facets []string) string {
	for _, facet := range facets {
		if strings.Contains(facet, "/conversation") && !strings.Contains(facet, "/data/conversation") &&
			!strings.Contains(facet, "metadata/conversation_id") {
			return "conversation"
		}
		if strings.Contains(facet, "/organization") {
			return "organization"
		}
	}
	return "document"
}

// Get or create cached card data
func (s *SearchService) getCachedCard(ctx context.Context, key string) (CardData, error) {
	if s.cacheCtrl.Client == nil {
		return nil, fmt.Errorf("cache not available")
	}

	data, err := s.cacheCtrl.Get(key)
	if err != nil {
		return nil, err
	}

	// Try to unmarshal as different card types
	var authorCard AuthorCardData
	if err := json.Unmarshal(data, &authorCard); err == nil && authorCard.Type == "author" {
		return authorCard, nil
	}

	var docketCard DocketCardData
	if err := json.Unmarshal(data, &docketCard); err == nil && docketCard.Type == "docket" {
		return docketCard, nil
	}

	var docCard DocumentCardData
	if err := json.Unmarshal(data, &docCard); err == nil && docCard.Type == "document" {
		return docCard, nil
	}

	return nil, fmt.Errorf("unknown card type in cache")
}

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
