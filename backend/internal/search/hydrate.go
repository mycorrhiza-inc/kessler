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

// type DocumentAuthor struct {
// 	AuthorName      string    `json:"author_name"`
// 	IsPerson        bool      `json:"is_person"`
// 	IsPrimaryAuthor bool      `json:"is_primary_author"`
// 	AuthorID        uuid.UUID `json:"author_id"`
// }

type DocumentAuthor struct {
	AuthorID        string `json:"author_id"`
	AuthorName      string `json:"author_name"`
	IsPerson        bool   `json:"is_person"`
	IsPrimaryAuthor bool   `json:"is_primary_author"`
	Description     string `json:"description"`
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

// Updated HydrateDocument function with proper error handling
func (s *SearchService) HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int) (CardData, error) {
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
			parsed, err := time.Parse(time.RFC3339, createdAt)
			if err == nil {
				timestamp = parsed
			}
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
	docid := uuid.MustParse(documentID)

	card := DocumentCardData{
		Name:        name,
		Description: description,
		Timestamp:   timestamp,
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "document",
		ObjectUUID:  docid,
		Authors:     []DocumentAuthor{}, // Initialize empty slice
	}

	// Get attachment authors with proper error handling
	queries := dbstore.New(s.db)

	attachmentResult, err := queries.GetAttachmentWithAuthors(ctx, docid)
	if err == nil {
		logger.Error(ctx, "Error parsing document ID as UUID", zap.Error(err))
	}
	// Check if we have valid JSON authors data
	if attachmentResult.AuthorsJson != "" && attachmentResult.AuthorsJson != "[]" {
		logger.Debug(ctx, "Raw authors JSON", zap.Dict(attachmentResult.AuthorsJson))

		// Parse the JSON into authors slice
		var authors []DocumentAuthor
		if err := json.Unmarshal([]byte(attachmentResult.AuthorsJson), &authors); err != nil {
			// Log the error but don't fail the entire operation
			logger.Error(ctx, "Error parsing authors JSON", zap.Error(err))
			logger.Error(ctx, "JSON content", zap.Dict(attachmentResult.AuthorsJson))
		} else {
			card.Authors = authors
			logger.Debug(ctx, "Successfully parsed authors", zap.Any("authors", authors))
		}
	} else {
		fmt.Printf("Error getting attachment with authors: %v\n", err)
	}

	// Cache and return
	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}

// Hydrate conversation/docket data
func (s *SearchService) HydrateConversation(ctx context.Context, id string, score float32, index int) (CardData, error) {
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
		Timestamp:   conversation.CreatedAt.Time,
		Index:       index,
		Type:        "docket",
		ObjectUUID:  conversation.ID,
	}

	// Cache and return
	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}

// Hydrate organization/author data
func (s *SearchService) HydrateOrganization(ctx context.Context, id string, score float32, index int) (CardData, error) {
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
		Timestamp:   org.CreatedAt.Time,
		ExtraInfo:   extraInfo,
		Index:       index,
		Type:        "author",
		ObjectUUID:  orgID,
	}

	// Cache and return
	s.cacheCard(ctx, cacheKey, card)

	return card, nil
}

// // Hydrate document data
// func (s *SearchService) _hydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int, full_fetch bool) (CardData, error) {
// 	log := logger.FromContext(ctx)
// 	// Check cache first
// 	cacheKey := cache.PrepareKey("search", "document", result.ID)
// 	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
// 		if doc, ok := cached.(DocumentCardData); ok {
// 			doc.Index = index // Update index for current search
// 			return doc, nil
// 		}
// 	}

// 	// Extract metadata
// 	name := ""
// 	description := result.Text
// 	timestamp := time.Now()
// 	extraInfo := ""

// 	if result.Metadata != nil {
// 		if fileName, ok := result.Metadata["file_name"].(string); ok {
// 			name = fileName
// 		}
// 		if desc, ok := result.Metadata["description"].(string); ok {
// 			description = desc
// 		}
// 		if createdAt, ok := result.Metadata["created_at"].(string); ok {
// 			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
// 				timestamp = t
// 			}
// 		}
// 		if caseNumber, ok := result.Metadata["case_number"].(string); ok {
// 			extraInfo = fmt.Sprintf("Case: %s", caseNumber)
// 		}
// 	}

// 	if len(result.ID) < 36 {
// 		return DocketCardData{}, fmt.Errorf("Document does not have long enough uuid.")
// 	}
// 	parsedUUID, err := uuid.Parse(result.ID[:36])
// 	if err != nil {
// 		return DocketCardData{}, fmt.Errorf("Could not parse uuid for object")
// 	}
// 	if full_fetch {
// 		// Fallback to database if metadata not provided
// 		needFetch := name == "" || description == result.Text
// 		if needFetch {
// 			q := dbstore.New(s.db)
// 			// Fetch file basic info
// 			if fileRec, err := q.ReadFile(ctx, parsedUUID); err == nil {
// 				if name == "" {
// 					name = fileRec.Name
// 				}
// 				timestamp = fileRec.DatePublished.Time
// 			} else {
// 				log.Warn("Failed to read file record", zap.String("file_id", parsedUUID.String()), zap.Error(err))
// 			}
// 			// Fetch metadata record
// 			if metaRec, err := q.FetchMetadata(ctx, parsedUUID); err == nil {
// 				var m map[string]interface{}
// 				if err := json.Unmarshal(metaRec.Mdata, &m); err == nil {
// 					if fn, ok := m["file_name"].(string); ok {
// 						name = fn
// 					}
// 					if desc2, ok := m["description"].(string); ok {
// 						description = desc2
// 					}
// 					if ca, ok := m["created_at"].(string); ok {
// 						if t2, err := time.Parse(time.RFC3339, ca); err == nil {
// 							timestamp = t2
// 						}
// 					}
// 					if cn, ok := m["case_number"].(string); ok {
// 						extraInfo = fmt.Sprintf("Case: %s", cn)
// 					}
// 				} else {
// 					log.Warn("Failed to unmarshal metadata JSON", zap.String("file_id", parsedUUID.String()), zap.Error(err))
// 				}
// 			} else {
// 				log.Warn("Failed to fetch metadata record", zap.String("file_id", parsedUUID.String()), zap.Error(err))
// 			}
// 		}
// 	}
// 	// If no name, use ID
// 	if name == "" {
// 		name = fmt.Sprintf("Document %s", result.ID)
// 	}
// 	card := DocumentCardData{
// 		Name:         name,
// 		Description:  description,
// 		Timestamp:    timestamp,
// 		ExtraInfo:    extraInfo,
// 		Index:        index,
// 		Type:         "document",
// 		ObjectUUID:   parsedUUID,
// 		FragmentID:   result.ID[36:],
// 		Authors:      []DocumentAuthor{}, // Would need to query authorship table
// 		Conversation: DocumentConversation{},
// 	}
// 	log.Debug("Successfully Created Initial Card Data", zap.String("file_id", parsedUUID.String()))

// 	// Try to get authors if this is a file document
// 	queries := dbstore.New(s.db)
// 	authorships, err := queries.AuthorshipDocumentListOrganizations(ctx, parsedUUID)
// 	if err != nil {
// 		log.Warn("Failed to list authorships", zap.String("file_id", parsedUUID.String()), zap.Error(err))
// 	} else if len(authorships) == 0 {
// 		log.Warn("No authorships found", zap.String("file_id", parsedUUID.String()))
// 	} else {
// 		for _, authorship := range authorships {
// 			// Get organization details
// 			org, err := queries.OrganizationRead(ctx, authorship.OrganizationID)
// 			if err != nil {
// 				log.Debug("Failed to read organization for authorship", zap.String("org_id", authorship.OrganizationID.String()), zap.Error(err))
// 				continue
// 			}
// 			author := DocumentAuthor{
// 				AuthorName:      org.Name,
// 				IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
// 				IsPrimaryAuthor: authorship.IsPrimaryAuthor.Valid && authorship.IsPrimaryAuthor.Bool,
// 				AuthorID:        org.ID,
// 			}
// 			card.Authors = append(card.Authors, author)
// 		}
// 	}

// 	// Try to get conversation info if this is a file document
// 	// conversation_uuid is stored in public.docket_documents
// 	conv_info, err := queries.ConversationIDFetchFromFileID(ctx, parsedUUID)
// 	if err != nil {
// 		log.Warn("Failed to fetch conversation ID from file ID", zap.String("file_id", parsedUUID.String()), zap.Error(err))
// 	} else if len(conv_info) == 0 {
// 		log.Warn("No conversation info found for file", zap.String("file_id", parsedUUID.String()))
// 	} else {
// 		// Fetch conversation details
// 		conv, err := queries.DocketConversationRead(ctx, conv_info[0].ConversationUuid)
// 		if err != nil {
// 			log.Debug("Failed to read conversation details", zap.String("conversation_id", conv_info[0].ConversationUuid.String()), zap.Error(err))
// 		} else {
// 			card.Conversation = DocumentConversation{
// 				ConvoName: conv.Name,
// 				ConvoID:   conv.ID,
// 			}
// 		}
// 	}

// 	// Cache the result
// 	s.cacheCard(ctx, cacheKey, card)

// 	return card, nil
// }

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
			card, err = s.HydrateDocument(ctx, result, i)
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
