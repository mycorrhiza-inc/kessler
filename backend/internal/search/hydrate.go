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

// Card data types matching the frontend requirements exactly
type CardData interface {
	GetType() string
}

type AuthorCardData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	ExtraInfo   string `json:"extraInfo,omitempty"`
	Index       int    `json:"index"`
	Type        string `json:"type"`
	ObjectUUID  string `json:"object_uuid"`
}

func (a AuthorCardData) GetType() string {
	return "author"
}

type DocketCardData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	Index       int    `json:"index"`
	Type        string `json:"type"`
	ObjectUUID  string `json:"object_uuid"`
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

type DocumentConversation struct {
	ConvoID   string `json:"convo_id"`
	ConvoName string `json:"convo_name"`
}

type DocumentCardData struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Timestamp    string                `json:"timestamp"`
	ExtraInfo    string                `json:"extraInfo,omitempty"`
	Index        int                   `json:"index"`
	Type         string                `json:"type"`
	ObjectUUID   string                `json:"object_uuid"`
	Authors      []DocumentAuthor      `json:"authors"`
	Conversation *DocumentConversation `json:"conversation,omitempty"`
}

func (d DocumentCardData) GetType() string {
	return "document"
}

// Hydrate document data from attachment
func (s *SearchService) hydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "document", result.ID)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if doc, ok := cached.(DocumentCardData); ok {
			doc.Index = index
			return doc, nil
		}
	}

	// Extract basic metadata
	name := ""
	description := result.Text
	timestamp := time.Now().Format(time.RFC3339)
	extraInfo := ""

	if result.Metadata != nil {
		if fileName, ok := result.Metadata["file_name"].(string); ok {
			name = fileName
		}
		if desc, ok := result.Metadata["description"].(string); ok {
			description = desc
		}
		if createdAt, ok := result.Metadata["created_at"].(string); ok {
			timestamp = createdAt
		}
		if caseNumber, ok := result.Metadata["case_number"].(string); ok {
			extraInfo = fmt.Sprintf("Case: %s", caseNumber)
		}
	}

	if name == "" {
		name = fmt.Sprintf("Document %s", result.ID)
	}

	documentID := result.ID
	if segmentIndex := strings.Index(documentID, "-segment-"); segmentIndex != -1 {
		documentID = documentID[:segmentIndex]
	}

	card := DocumentCardData{
		Name:        name,
		Description: description,
		Timestamp:   timestamp,
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "document",
		ObjectUUID:  documentID,
		Authors:     []DocumentAuthor{},
	}

	// Remove segment suffix from ID

	// Get attachment authors
	if attachmentID, err := uuid.Parse(documentID); err == nil {
		queries := dbstore.New(s.db)

		if attachmentResult, err := queries.GetAttachmentWithAuthors(ctx, attachmentID); err == nil {
			if attachmentResult.AuthorsJson != "" && attachmentResult.AuthorsJson != "[]" {
				json.Unmarshal([]byte(attachmentResult.AuthorsJson), &card.Authors)
			}
		}
	}

	// Cache and return
	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}

// Hydrate conversation/docket data
func (s *SearchService) hydrateConversation(ctx context.Context, id string, score float32, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "conversation", id)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if docket, ok := cached.(DocketCardData); ok {
			docket.Index = index
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
		Timestamp:   conversation.CreatedAt.Time.Format(time.RFC3339),
		Index:       index,
		Type:        "docket",
		ObjectUUID:  conversation.ID.String(),
	}

	// Cache and return
	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}

// Hydrate organization/author data
func (s *SearchService) hydrateOrganization(ctx context.Context, id string, score float32, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "organization", id)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if author, ok := cached.(AuthorCardData); ok {
			author.Index = index
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
		Timestamp:   org.CreatedAt.Time.Format(time.RFC3339),
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "author",
		ObjectUUID:  org.ID.String(),
	}

	// Cache and return
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
