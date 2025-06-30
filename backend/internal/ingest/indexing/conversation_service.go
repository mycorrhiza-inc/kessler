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

		metadata, facets := ci.buildConversationMetadataAndFacets(conversationParams{
			id:           c.ID,
			docketGovID:  c.DocketGovID,
			state:        c.State,
			matterType:   c.MatterType,
			industryType: c.IndustryType,
			description:  c.Description,
		})

		recs = append(recs, fugusdk.ObjectRecord{
			ID:        c.ID.String(),
			Text:      text,
			Metadata:  metadata,
			Facets:    facets,
			Namespace: ci.svc.defaultNamespace,
			DataType:  "data/conversation",
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

	metadata, facets := ci.buildConversationMetadataAndFacets(conversationParams{
		id:           c.ID,
		docketGovID:  c.DocketGovID,
		state:        c.State,
		matterType:   c.MatterType,
		industryType: c.IndustryType,
		description:  c.Description,
	})

	rec := fugusdk.ObjectRecord{
		ID:        c.ID.String(),
		Text:      text,
		Metadata:  metadata,
		Facets:    facets,
		Namespace: ci.svc.defaultNamespace,
		DataType:  "data/conversation",
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

// conversationParams holds the parameters for building conversation metadata and facets
type conversationParams struct {
	id           uuid.UUID
	docketGovID  string
	state        string
	matterType   string
	industryType string
	description  string
}

// buildConversationMetadataAndFacets creates both metadata and facets for a conversation record
func (ci *ConversationIndexer) buildConversationMetadataAndFacets(params conversationParams) (map[string]interface{}, []string) {
	metadata := map[string]interface{}{
		"entity_type":     "conversation",
		"conversation_id": params.id.String(),
	}

	var facets []string

	// Add namespace facets
	facets = append(facets, ci.svc.defaultNamespace)
	facets = append(facets, fmt.Sprintf("%s/data/conversation", ci.svc.defaultNamespace))

	// Core facets with embedded values
	facets = append(facets, "metadata/entity_type/conversation")
	facets = append(facets, fmt.Sprintf("metadata/conversation_id/%s", params.id.String()))

	// Add metadata and facets for each field if not empty
	if params.docketGovID != "" {
		metadata["docket_gov_id"] = params.docketGovID
		facets = append(facets, fmt.Sprintf("metadata/docket_gov_id/%s", params.docketGovID))
	}

	if params.state != "" {
		metadata["state"] = params.state
		facets = append(facets, fmt.Sprintf("metadata/state/%s", params.state))
	}

	if params.matterType != "" {
		metadata["matter_type"] = params.matterType
		facets = append(facets, fmt.Sprintf("metadata/matter_type/%s", params.matterType))
	}

	if params.industryType != "" {
		metadata["industry_type"] = params.industryType
		facets = append(facets, fmt.Sprintf("metadata/industry_type/%s", params.industryType))
	}

	if params.description != "" {
		metadata["description"] = params.description
	}

	return metadata, facets

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

		// Build metadata and facets
		baseMetadata, facets := ci.buildConversationMetadataAndFacets(conversationParams{
			id:           c.ID,
			docketGovID:  c.DocketGovID,
			state:        c.State,
			matterType:   c.MatterType,
			industryType: c.IndustryType,
			description:  c.Description,
		})

		// Add custom metadata updates
		for key, value := range metadata {
			baseMetadata[key] = value
		}

		rec := fugusdk.ObjectRecord{
			ID:        conversationID,
			Text:      text,
			Metadata:  baseMetadata,
			Facets:    facets,
			Namespace: ci.svc.defaultNamespace,
			DataType:  "data/conversation",
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
