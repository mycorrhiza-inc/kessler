// indexing/conversation_service.go
package indexing

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"kessler/internal/fugusdk"
	"kessler/pkg/database"
	"kessler/pkg/logger"
)

// ConversationIndexer handles conversation-specific indexing operations
type ConversationIndexer struct {
	svc *IndexService
}

// NewConversationIndexer creates a new conversation indexer
func NewConversationIndexer(svc *IndexService) *ConversationIndexer {
	return &ConversationIndexer{svc: svc}
}

// IndexAllConversations retrieves all conversations and batch indexes them in chunks.
func (ci *ConversationIndexer) IndexAllConversations(ctx context.Context) (int, error) {
	q := database.GetQueries(ci.svc.db)
	rows, err := q.DocketConversationList(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all conversations: %w", err)
	}

	var recs []fugusdk.ObjectRecord
	skippedCount := 0

	for _, c := range rows {
		// Create meaningful text field with fallbacks
		text := ci.createConversationText(c.Name, c.Description, c.DocketGovID, c.ID.String())

		if text == "" {
			log.Printf("Skipping conversation %s - no valid text content", c.ID.String())
			skippedCount++
			continue
		}

		recs = append(recs, fugusdk.ObjectRecord{
			ID:   c.ID.String(),
			Text: text,
			Metadata: map[string]interface{}{
				"docket_gov_id":   c.DocketGovID,
				"description":     c.Description,
				"state":           c.State,
				"matter_type":     c.MatterType,
				"industry_type":   c.IndustryType,
				"conversation_id": c.ID.String(), // Store specific ID in metadata
				"entity_type":     "conversation",
			},
			// Use proper namespace facet structure (categorical only)
			Namespace:      ci.svc.defaultNamespace,
			ConversationID: c.ID.String(),  // This triggers namespace/NYPUC/conversation facet
			DataType:       "conversation", // This triggers namespace/NYPUC/data facet
		})
	}

	if skippedCount > 0 {
		log.Printf("Skipped %d conversations with empty content", skippedCount)
	}

	if len(recs) == 0 {
		log.Printf("No valid conversations to index")
		return 0, nil
	}

	client, err := ci.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	return ci.svc.processBatchInChunks(ctx, client, recs, "conversations")
}

// IndexConversationByID retrieves one conversation by UUID and indexes it.
func (ci *ConversationIndexer) IndexConversationByID(ctx context.Context, idStr string) (int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid conversation id: %w", err)
	}

	q := database.GetQueries(ci.svc.db)
	c, err := q.DocketConversationRead(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("read conversation: %w", err)
	}

	// Create meaningful text field with fallbacks
	text := ci.createConversationText(c.Name, c.Description, c.DocketGovID, c.ID.String())
	if text == "" {
		return 0, fmt.Errorf("conversation %s has no valid text content and cannot be indexed", idStr)
	}

	rec := fugusdk.ObjectRecord{
		ID:   c.ID.String(),
		Text: text,
		Metadata: map[string]interface{}{
			"docket_gov_id":   c.DocketGovID,
			"description":     c.Description,
			"state":           c.State,
			"matter_type":     c.MatterType,
			"industry_type":   c.IndustryType,
			"conversation_id": c.ID.String(), // Store specific ID in metadata
			"entity_type":     "conversation",
		},
		// Use proper namespace facet structure (categorical only)
		Namespace:      ci.svc.defaultNamespace,
		ConversationID: c.ID.String(),
		DataType:       "conversation",
	}

	client, err := ci.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	response, err := client.AddOrUpdateObject(ctx, rec)
	if err != nil {
		return 0, fmt.Errorf("index conversation: %w", err)
	}

	log.Printf("Successfully indexed conversation %s: %s", idStr, response.Message)
	return 1, nil
}

// DeleteConversationFromIndex removes a conversation from the search index.
func (ci *ConversationIndexer) DeleteConversationFromIndex(ctx context.Context, idStr string) error {
	client, err := ci.svc.createFuguClient(ctx)
	if err != nil {
		return fmt.Errorf("new fugu client: %w", err)
	}

	response, err := client.DeleteObject(ctx, idStr)
	if err != nil {
		return fmt.Errorf("delete conversation from index: %w", err)
	}

	log.Printf("Successfully deleted conversation %s from index: %s", idStr, response.Message)
	return nil
}

// createConversationText creates a meaningful text field for conversation indexing
// with multiple fallback options to ensure we always have searchable content
func (ci *ConversationIndexer) createConversationText(name, description, docketGovID, id string) string {
	// Try name first (most common case)
	if text := strings.TrimSpace(name); text != "" {
		return text
	}

	// Fall back to description
	if text := strings.TrimSpace(description); text != "" {
		return text
	}

	// Fall back to docket gov ID with meaningful prefix
	if text := strings.TrimSpace(docketGovID); text != "" {
		return fmt.Sprintf("Docket %s", text)
	}

	// Last resort: use UUID with prefix
	return fmt.Sprintf("Conversation %s", id)
}

// ValidateConversationData validates conversation data before indexing
func (ci *ConversationIndexer) ValidateConversationData(conversationID string) error {
	if conversationID == "" {
		return fmt.Errorf("conversation ID cannot be empty")
	}

	if _, err := uuid.Parse(conversationID); err != nil {
		return fmt.Errorf("invalid conversation UUID format: %w", err)
	}

	return nil
}

// GetConversationStats returns statistics about conversations in the database
func (ci *ConversationIndexer) GetConversationStats(ctx context.Context) (map[string]interface{}, error) {
	q := database.GetQueries(ci.svc.db)

	// Get total count
	rows, err := q.DocketConversationList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation list: %w", err)
	}

	stats := map[string]interface{}{
		"total_conversations": len(rows),
		"namespace":           ci.svc.defaultNamespace,
	}

	// Count by state if available
	stateCount := make(map[string]int)
	textlessCount := 0

	for _, c := range rows {
		if c.State != "" {
			stateCount[c.State]++
		}

		text := ci.createConversationText(c.Name, c.Description, c.DocketGovID, c.ID.String())
		if text == fmt.Sprintf("Conversation %s", c.ID.String()) {
			textlessCount++
		}
	}

	stats["by_state"] = stateCount
	stats["conversations_without_meaningful_text"] = textlessCount

	return stats, nil
}

// BulkUpdateConversationMetadata updates metadata for multiple conversations
func (ci *ConversationIndexer) BulkUpdateConversationMetadata(ctx context.Context, updates map[string]map[string]interface{}) error {
	client, err := ci.svc.createFuguClient(ctx)
	if err != nil {
		return fmt.Errorf("new fugu client: %w", err)
	}

	for conversationID, metadata := range updates {
		// Validate conversation ID
		if err := ci.ValidateConversationData(conversationID); err != nil {
			logger.Error(ctx, "invalid conversation ID for metadata update",
				zap.String("conversation_id", conversationID),
				zap.Error(err))
			continue
		}

		// Get existing conversation to preserve text content
		id, _ := uuid.Parse(conversationID)
		q := database.GetQueries(ci.svc.db)
		c, err := q.DocketConversationRead(ctx, id)
		if err != nil {
			logger.Error(ctx, "failed to read conversation for metadata update",
				zap.String("conversation_id", conversationID),
				zap.Error(err))
			continue
		}

		// Create text field
		text := ci.createConversationText(c.Name, c.Description, c.DocketGovID, c.ID.String())
		if text == "" {
			logger.Warn(ctx, "skipping conversation metadata update - no valid text content",
				zap.String("conversation_id", conversationID))
			continue
		}

		// Merge existing metadata with updates
		updatedMetadata := map[string]interface{}{
			"docket_gov_id":   c.DocketGovID,
			"description":     c.Description,
			"state":           c.State,
			"matter_type":     c.MatterType,
			"industry_type":   c.IndustryType,
			"conversation_id": c.ID.String(),
			"entity_type":     "conversation",
		}

		// Add custom metadata updates
		for key, value := range metadata {
			updatedMetadata[key] = value
		}

		rec := fugusdk.ObjectRecord{
			ID:             conversationID,
			Text:           text,
			Metadata:       updatedMetadata,
			Namespace:      ci.svc.defaultNamespace,
			ConversationID: conversationID,
			DataType:       "conversation",
		}

		if _, err := client.AddOrUpdateObject(ctx, rec); err != nil {
			logger.Error(ctx, "failed to update conversation metadata",
				zap.String("conversation_id", conversationID),
				zap.Error(err))
		} else {
			logger.Info(ctx, "successfully updated conversation metadata",
				zap.String("conversation_id", conversationID))
		}
	}

	return nil
}
