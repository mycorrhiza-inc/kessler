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
		return fmt.Errorf("Error getting attachment with authors: %v\n", err)
	}
	return nil
}

// Updated HydrateDocument function with proper error handling
func (s *SearchService) HydrateDocument(ctx context.Context, result fugusdk.FuguSearchResult, index int) (CardData, error) {
	// Check cache first
	cacheKey := cache.PrepareKey("search", "document", result.ID)

	cached, err := s.getCachedCard(ctx, cacheKey)
	if err != nil {
		logger.Error(ctx, "Error getting cached card card ", zap.Error(err))
	}

	doc, ok := cached.(DocumentCardData)
	if ok {
		doc.Index = index
		return doc, nil
	}

	// validate and hydrate card data
	card := DocumentCardData{}

	if result.Metadata != nil {
		fileName, ok := result.Metadata["file_name"].(string)
		if ok {
			card.Name = fileName
		}

		fileIDString, ok := result.Metadata["file_id"].(string)
		if ok {
			fileID, err := uuid.Parse(fileIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse file_id in metadata", err)
			}
			card.ObjectUUID = fileID

		}

		convoIDString, ok := result.Metadata["conversation_id"].(string)
		if ok {
			convoID, err := uuid.Parse(convoIDString)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse conversation_id in metadata", err)
			}

			err = s.hydrateDocumentConvos(ctx, &card, convoID)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to hydrate conversaton", err)
			}
		}

		authorIDsRaw, ok := result.Metadata["author_ids"].([]string)
		if ok {
			parseUUID := func(val string) (uuid.UUID, error) { return uuid.Parse(val) }
			authorIDs, err := util.MapErrorBubble(authorIDsRaw, parseUUID)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to parse an author_id in author_ids in metadata")
			}
			err = s.hydrateDocumentAuthors(ctx, &card, authorIDs)
			if err != nil {
				return DocumentCardData{}, fmt.Errorf("Failed to hydrate authors in metadata", err)
			}
		}

		desc, ok := result.Metadata["description"].(string)
		if ok {
			card.Description = desc
		}

		createdAt, ok := result.Metadata["created_at"].(string)
		if ok {
			t, err := time.Parse(time.RFC3339, createdAt)
			if err == nil {
				card.Timestamp = t
			}
		}

		caseNumber, ok := result.Metadata["case_number"].(string)
		if ok {
			card.ExtraInfo = fmt.Sprintf("Case: %s", caseNumber)
		}
	}

	const UUID_LEN = 36
	segmentIndex := strings.Index(result.ID, "-segment-")
	if len(result.ID) < UUID_LEN {
		return DocumentCardData{}, fmt.Errorf("Document does not have long enough uuid.")
	}
	// WARN: assumes there is no prefixing
	card.AttachmentUUID, err = uuid.Parse(result.ID[:segmentIndex])
	if err != nil {
		return DocumentCardData{}, fmt.Errorf("Could not parse uuid for object")
	}

	if segmentIndex > 0 {
		card.FragmentID = result.ID[segmentIndex:]
	}

	if card.Name == "" {
		card.Name = fmt.Sprintf("Document %s", result.ID)
	}

	s.cacheCard(ctx, cacheKey, card)
	return card, nil
}
