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

func (s *SearchService) hydrateDocumentConvos(ctx context.Context, card *DocumentCardData, convoID uuid.UUID) error {
	queries := dbstore.New(s.db)

	conv, err := queries.DocketConversationRead(ctx, convoID)
	if err != nil {
		log.Warn("Failed to read conversation details", zap.String("conversation_id", convoID.String()), zap.Error(err))
		return err
	} else {
		card.Conversation = DocumentConversation{
			ConvoName: conv.Name,
			ConvoID:   conv.ID,
		}
		return nil
	}

}
func (s *SearchService) hydrateDocumentAuthors(ctx context.Context, card *DocumentCardData, authorIDs []uuid.UUID) error {
	queries := dbstore.New(s.db)

	attachmentResult, err := queries.GetAttachmentWithAuthors(ctx, card.ObjectUUID)
	if err != nil {
		logger.Error(ctx, "Error getting attachment with authors", zap.Error(err))
		return err // CRITICAL: Must return here
	}

	// Check if we have valid JSON authors data
	if attachmentResult.AuthorsJson != "" && attachmentResult.AuthorsJson != "[]" {
		logger.Debug(ctx, "Raw authors JSON", zap.String("json", attachmentResult.AuthorsJson))

		// Parse the JSON into authors slice
		var authors []DocumentAuthor
		if err := json.Unmarshal([]byte(attachmentResult.AuthorsJson), &authors); err != nil {
			logger.Error(ctx, "Error parsing authors JSON", zap.Error(err))
			logger.Error(ctx, "JSON content", zap.String("json", attachmentResult.AuthorsJson))
			return err // Return error instead of continuing silently
		}

		card.Authors = authors
		logger.Debug(ctx, "Successfully parsed authors", zap.Any("authors", authors))
	}

	return nil
}

func (s *SearchService) HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "document", result.ID)

	if cached, err := s.getCachedCard(ctx, cacheKey); err == nil {
		if doc, ok := cached.(DocumentCardData); ok {
			doc.Index = index
			return doc, nil
		}
	}

	// validate and hydrate card data
	card := DocumentCardData{
		Index: index,
	}

	// Parse AttachmentUUID and FragmentID from result.ID first
	segmentIndex := strings.Index(result.ID, "-segment-")
	if segmentIndex == -1 {
		// No segment found - this might be a full document reference
		var err error
		card.AttachmentUUID, err = uuid.Parse(result.ID)
		if err != nil {
			return DocumentCardData{}, fmt.Errorf("Could not parse document ID as UUID: %w", err)
		}
		card.FragmentID = ""
	} else {
		// Validate ID length before parsing
		if len(result.ID) < 36 { // UUID_LEN
			return DocumentCardData{}, fmt.Errorf("Document ID too short for UUID: %s", result.ID)
		}

		var err error
		card.AttachmentUUID, err = uuid.Parse(result.ID[:segmentIndex])
		if err != nil {
			return DocumentCardData{}, fmt.Errorf("Could not parse attachment UUID: %w", err)
		}
		card.FragmentID = result.ID[segmentIndex+len("-segment-"):]
	}

	if result.Metadata != nil {
		// File name
		if fileName, ok := result.Metadata["file_name"].(string); ok {
			card.Name = fileName
		}

		// ObjectUUID
		if fileIDString, ok := result.Metadata["file_id"].(string); ok {
			fileID, err := uuid.Parse(fileIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse file_id in metadata: %w", err)
			}
			card.ObjectUUID = fileID
		}

		// Conversation
		if convoIDString, ok := result.Metadata["conversation_id"].(string); ok {
			convoID, err := uuid.Parse(convoIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse conversation_id in metadata: %w", err)
			}

			err = s.hydrateDocumentConvos(ctx, &card, convoID)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to hydrate conversation: %w", err)
			}
		}

		// Authors
		if authorIDsRaw, ok := result.Metadata["author_ids"].([]string); ok {
			parseUUID := func(val string) (uuid.UUID, error) { return uuid.Parse(val) }
			authorIDs, err := util.MapErrorBubble(authorIDsRaw, parseUUID)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse author_ids in metadata: %w", err)
			}
			err = s.hydrateDocumentAuthors(ctx, &card, authorIDs)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to hydrate authors: %w", err)
			}
		}

		// Description
		if desc, ok := result.Metadata["description"].(string); ok {
			card.Description = desc
		}

		// Timestamp
		if createdAt, ok := result.Metadata["created_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				card.Timestamp = t
			} else {
				logger.Warn(ctx, "Failed to parse created_at timestamp", zap.Error(err))
			}
		}

		// Case number
		if caseNumber, ok := result.Metadata["case_number"].(string); ok {
			card.ExtraInfo = fmt.Sprintf("Case: %s", caseNumber)
		}
	}

	// Set default name if not provided
	if card.Name == "" {
		card.Name = fmt.Sprintf("Document %s", result.ID)
	}

	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}
