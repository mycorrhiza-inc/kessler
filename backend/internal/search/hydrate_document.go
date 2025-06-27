package search

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/cache"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
	"kessler/pkg/util"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"

	//"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// Updated HydrateDocument function with proper error handling
func (s *SearchService) HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int, full_fetch bool) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "document", result.ID)
	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if doc, ok := cached.(DocumentCardData); ok {
			doc.Index = index
			return doc, nil
		}
	}

	// Extract metadata
	name := ""
	description := result.Text
	timestamp := time.Now()
	var err error
	extraInfo := ""
	fileID := uuid.Nil
	convoID := uuid.Nil
	authorIDs := []uuid.UUID{}

	if result.Metadata != nil {
		if fileName, ok := result.Metadata["file_name"].(string); ok {
			name = fileName
		}

		if fileIDString, ok := result.Metadata["file_id"].(string); ok {
			fileID, err = uuid.Parse(fileIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse file_id in metadata")
			}
		}

		if convoIDString, ok := result.Metadata["conversation_id"].(string); ok {
			convoID, err = uuid.Parse(convoIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse conversation_id in metadata")
			}
		}

		if authorIDsRaw, ok := result.Metadata["author_ids"].([]string); ok {
			parseUUID := func(val string) (uuid.UUID, error) { return uuid.Parse(val) }
			authorIDs, err = util.MapErrorBubble(authorIDsRaw, parseUUID)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse an author_id in author_ids in metadata")
			}
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

	if name == "" {
		name = fmt.Sprintf("Document %s", result.ID)
	}

	attachmentIDRaw := result.ID
	attachID := uuid.New()
	fragmentID := ""
	if segmentIndex := strings.Index(attachmentIDRaw, "-segment-"); segmentIndex != -1 {
		fragmentID = attachmentIDRaw[segmentIndex:]
		attachmentIDRaw = attachmentIDRaw[:segmentIndex]
	}
	attachID, err = uuid.Parse(attachmentIDRaw)
	if err != nil {
		log.Error("Could not parse attachmentID", zap.String("attach_id", attachmentIDRaw))
	}

	card := DocumentCardData{
		Name:           name,
		Description:    description,
		Timestamp:      timestamp,
		ExtraInfo:      extraInfo,
		Index:          index,
		Type:           "document",
		ObjectUUID:     fileID,
		AttachmentUUID: attachID,
		FragmentID:     fragmentID,
		Authors:        []DocumentAuthor{}, // Initialize empty slice
		Conversation:   DocumentConversation{},
	}

	// Get attachment authors with proper error handling
	queries := dbstore.New(s.db)

	if full_fetch {
		attachmentResult, err := queries.GetAttachmentWithAuthors(ctx, attachID)
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
	}

	if !full_fetch {
		// TODO: Run both of these concurrently with a goroutine or something.

		// Lookup organization details
		for _, orgID := range authorIDs {
			if orgID == uuid.Nil {
				log.Error("Encountered null org_id when hydrating without a full fetch", zap.String("attach_id", attachID.String()), zap.String("file_id", fileID.String()))
				return DocumentCardData{}, fmt.Errorf("Encountered null org_id when hydrating without a full fetch")
			}
			org, err := queries.OrganizationRead(ctx, orgID)
			if err != nil {
				log.Warn("Failed to read organization for authorship", zap.String("org_id", orgID.String()), zap.Error(err))
				continue
			}
			card.Authors = append(card.Authors, DocumentAuthor{
				AuthorName:      org.Name,
				IsPerson:        org.IsPerson.Valid && org.IsPerson.Bool,
				IsPrimaryAuthor: true,
				AuthorID:        org.ID,
			})
		}

		// Lookup conversation details
		if convoID == uuid.Nil {
			log.Error("Encountered null convo_id when hydrating without a full fetch", zap.String("attach_id", attachID.String()), zap.String("file_id", fileID.String()))
			return DocumentCardData{}, fmt.Errorf("Encountered null convo_id when hydrating without a full fetch")
		}
		conv, err := queries.DocketConversationRead(ctx, convoID)
		if err != nil {
			log.Warn("Failed to read conversation details", zap.String("conversation_id", convoID.String()), zap.Error(err))
		} else {
			card.Conversation = DocumentConversation{
				ConvoName: conv.Name,
				ConvoID:   conv.ID,
			}
		}
	}
	// Cache and return
	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}
