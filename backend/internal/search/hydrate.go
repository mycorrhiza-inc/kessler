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
	ObjectUUID  uuid.UUID `json:"object_uuid"`
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
	ObjectUUID  uuid.UUID `json:"object_uuid"`
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
	ConvoName string    `json:"convo_name"`
	ConvoID   uuid.UUID `json:"convo_id"`
}

type DocumentCardData struct {
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Timestamp    time.Time            `json:"timestamp"`
	ExtraInfo    string               `json:"extraInfo,omitempty"`
	Index        int                  `json:"index"`
	Type         string               `json:"type"`
	ObjectUUID   uuid.UUID            `json:"object_uuid"`
	FragmentID   string               `json:"fragment_id"`
	Authors      []DocumentAuthor     `json:"authors"`
	Conversation DocumentConversation `json:"conversation"`
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
func (s *SearchService) HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int, full_fetch bool) (CardData, error) {
	log := logger.FromContext(ctx)
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

	if len(result.ID) < 36 {
		return DocketCardData{}, fmt.Errorf("Document does not have long enough uuid.")
	}
	parsedUUID, err := uuid.Parse(result.ID[:36])
	if err != nil {
		return DocketCardData{}, fmt.Errorf("Could not parse uuid for object")
	}
	if full_fetch {
		// Fallback to database if metadata not provided
		needFetch := name == "" || description == result.Text
		if needFetch {
			q := dbstore.New(s.db)
			// Fetch file basic info
			if fileRec, err := q.ReadFile(ctx, parsedUUID); err == nil {
				if name == "" {
					name = fileRec.Name
				}
				timestamp = fileRec.DatePublished.Time
			} else {
				log.Warn("Failed to read file record", zap.String("file_id", parsedUUID.String()), zap.Error(err))
			}
			// Fetch metadata record
			if metaRec, err := q.FetchMetadata(ctx, parsedUUID); err == nil {
				var m map[string]interface{}
				if err := json.Unmarshal(metaRec.Mdata, &m); err == nil {
					if fn, ok := m["file_name"].(string); ok {
						name = fn
					}
					if desc2, ok := m["description"].(string); ok {
						description = desc2
					}
					if ca, ok := m["created_at"].(string); ok {
						if t2, err := time.Parse(time.RFC3339, ca); err == nil {
							timestamp = t2
						}
					}
					if cn, ok := m["case_number"].(string); ok {
						extraInfo = fmt.Sprintf("Case: %s", cn)
					}
				} else {
					log.Warn("Failed to unmarshal metadata JSON", zap.String("file_id", parsedUUID.String()), zap.Error(err))
				}
			} else {
				log.Warn("Failed to fetch metadata record", zap.String("file_id", parsedUUID.String()), zap.Error(err))
			}
		}
	}
	// If no name, use ID
	if name == "" {
		name = fmt.Sprintf("Document %s", result.ID)
	}
	card := DocumentCardData{
		Name:         name,
		Description:  description,
		Timestamp:    timestamp,
		ExtraInfo:    extraInfo,
		Index:        index,
		Type:         "document",
		ObjectUUID:   parsedUUID,
		FragmentID:   result.ID[36:],
		Authors:      []DocumentAuthor{}, // Would need to query authorship table
		Conversation: DocumentConversation{},
	}
	log.Info("Successfully Created Initial Card Data", zap.String("file_id", parsedUUID.String()))

	// Try to get authors if this is a file document
	queries := dbstore.New(s.db)
	authorships, err := queries.AuthorshipDocumentListOrganizations(ctx, parsedUUID)
	if err != nil {
		log.Warn("Failed to list authorships", zap.String("file_id", parsedUUID.String()), zap.Error(err))
	} else if len(authorships) == 0 {
		log.Warn("No authorships found", zap.String("file_id", parsedUUID.String()))
	} else {
		for _, authorship := range authorships {
			// Get organization details
			org, err := queries.OrganizationRead(ctx, authorship.OrganizationID)
			if err != nil {
				log.Info("Failed to read organization for authorship", zap.String("org_id", authorship.OrganizationID.String()), zap.Error(err))
				continue
			}
			author := DocumentAuthor{
				AuthorName:      org.Name,
				IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
				IsPrimaryAuthor: authorship.IsPrimaryAuthor.Valid && authorship.IsPrimaryAuthor.Bool,
				AuthorID:        org.ID,
			}
			card.Authors = append(card.Authors, author)
		}
	}

	// Try to get conversation info if this is a file document
	// conversation_uuid is stored in public.docket_documents
	conv_info, err := queries.ConversationIDFetchFromFileID(ctx, parsedUUID)
	if err != nil {
		log.Warn("Failed to fetch conversation ID from file ID", zap.String("file_id", parsedUUID.String()), zap.Error(err))
	} else if len(conv_info) == 0 {
		log.Warn("No conversation info found for file", zap.String("file_id", parsedUUID.String()))
	} else {
		// Fetch conversation details
		conv, err := queries.DocketConversationRead(ctx, conv_info[0].ConversationUuid)
		if err != nil {
			log.Info("Failed to read conversation details", zap.String("conversation_id", conv_info[0].ConversationUuid.String()), zap.Error(err))
		} else {
			card.Conversation = DocumentConversation{
				ConvoName: conv.Name,
				ConvoID:   conv.ID,
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
