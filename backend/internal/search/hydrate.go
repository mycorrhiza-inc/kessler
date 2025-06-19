package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	//"go.opentelemetry.io/otel"
	"go.uber.org/zap"
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
	ObjectUUID  string    `json:"object_uuid"`
}

func (a AuthorCardData) GetType() string {
	return "author"
}

type DocketCardData struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Index       int       `json:"index"`
	Type        string    `json:"type"`
	ObjectUUID  string    `json:"object_uuid"`
}

func (d DocketCardData) GetType() string {
	return "docket"
}

type DocumentAuthor struct {
	AuthorName      string `json:"author_name"`
	IsPerson        bool   `json:"is_person"`
	IsPrimaryAuthor bool   `json:"is_primary_author"`
	AuthorID        string `json:"author_id"`
}

type DocumentCardData struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Timestamp   time.Time        `json:"timestamp"`
	ExtraInfo   string           `json:"extraInfo,omitempty"`
	Index       int              `json:"index"`
	Type        string           `json:"type"`
	ObjectUUID  string           `json:"object_uuid"`
	Authors     []DocumentAuthor `json:"authors"`
}

func (d DocumentCardData) GetType() string {
	return "document"
}

// Updated search response to use card data

// NewSearchService creates a new search service with database access

// Helper to determine result type from facets
func (s *SearchService) getResultType(facets []string) string {
	for _, facet := range facets {
		if strings.Contains(facet, "/conversation") && !strings.Contains(facet, "/data/conversation") {
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

// Hydrate conversation/docket data
func (s *SearchService) hydrateConversation(ctx context.Context, id string, score float32, index int) (CardData, error) {
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
		ObjectUUID:  conversation.ID.String(),
	}

	// Cache the result
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}

// Hydrate organization/author data
func (s *SearchService) hydrateOrganization(ctx context.Context, id string, score float32, index int) (CardData, error) {
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
		ObjectUUID:  org.ID.String(),
	}

	// Cache the result
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}

// Hydrate document data
func (s *SearchService) hydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "document", result.ID)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if doc, ok := cached.(DocumentCardData); ok {
			doc.Index = index // Update index for current search
			return doc, nil
		}
	}

	// Extract metadata
	name := ""
	description := result.Text
	timestamp := time.Now()
	extraInfo := ""

	if result.Metadata != nil {
		if fileName, ok := result.Metadata["file_name"].(string); ok {
			name = fileName
		}
		if desc, ok := result.Metadata["description"].(string); ok {
			description = desc
		}
		if createdAt, ok := result.Metadata["created_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				timestamp = t
			}
		}
		if caseNumber, ok := result.Metadata["case_number"].(string); ok {
			extraInfo = fmt.Sprintf("Case: %s", caseNumber)
		}
	}

	// If no name, use ID
	if name == "" {
		name = fmt.Sprintf("Document %s", result.ID)
	}

	card := DocumentCardData{
		Name:        name,
		Description: description,
		Timestamp:   timestamp,
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "document",
		ObjectUUID:  result.ID,
		Authors:     []DocumentAuthor{}, // Would need to query authorship table
	}

	// Try to get authors if this is a file document
	if fileID, err := uuid.Parse(result.ID); err == nil {
		queries := dbstore.New(s.db)
		authorships, err := queries.AuthorshipDocumentListOrganizations(ctx, fileID)
		if err == nil && len(authorships) > 0 {
			for _, authorship := range authorships {
				// Get organization details
				org, err := queries.OrganizationRead(ctx, authorship.OrganizationID)
				if err == nil {
					author := DocumentAuthor{
						AuthorName:      org.Name,
						IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
						IsPrimaryAuthor: authorship.IsPrimaryAuthor.Valid && authorship.IsPrimaryAuthor.Bool,
						AuthorID:        org.ID.String(),
					}
					card.Authors = append(card.Authors, author)
				}
			}
		}
	}

	// Cache the result
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}

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
			card, err = s.hydrateConversation(ctx, result.ID, result.Score, i)
		case "organization":
			card, err = s.hydrateOrganization(ctx, result.ID, result.Score, i)
		default:
			card, err = s.hydrateDocument(ctx, result, i)
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
