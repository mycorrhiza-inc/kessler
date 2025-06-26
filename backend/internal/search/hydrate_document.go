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
	docid := uuid.New()
	if segmentIndex := strings.Index(documentID, "-segment-"); segmentIndex != -1 {
		documentID = documentID[:segmentIndex]
		docid = uuid.MustParse(documentID)
	}

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
	attachmentID, err := uuid.Parse(documentID)
	if err != nil {
	}
	queries := dbstore.New(s.db)

	attachmentResult, err := queries.GetAttachmentWithAuthors(ctx, attachmentID)
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
