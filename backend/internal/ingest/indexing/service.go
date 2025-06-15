package indexing

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"kessler/internal/database"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
)

// IndexService fetches DB records and indexes them into FuguDB.
type IndexService struct {
	fuguURL string
	q       *dbstore.Queries
}

// NewIndexService constructs an IndexService pointing at fuguURL.
func NewIndexService(fuguURL string) *IndexService {
	// Use the shared connection pool (implements dbstore.DBTX)
	db := database.ConnPool
	return &IndexService{
		fuguURL: fuguURL,
		q:       dbstore.New(db),
	}
}

// IndexAllConversations retrieves all conversations and batch indexes them in chunks.
func (s *IndexService) IndexAllConversations(ctx context.Context) (int, error) {
	rows, err := s.q.ConversationCompleteQuickwitListGet(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all conversations: %w", err)
	}

	var recs []fugusdk.ObjectRecord
	for _, c := range rows {
		recs = append(recs, fugusdk.ObjectRecord{
			ID:   c.ID.String(),
			Text: c.Name,
			Metadata: map[string]interface{}{
				"docket_gov_id":   c.DocketGovID,
				"description":     c.Description,
				"total_documents": c.TotalDocuments,
			},
			Namespace: "conversations", // Add namespace for better organization
		})
	}

	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	// Process in chunks to avoid batch size limits
	return s.processBatchInChunks(ctx, client, recs, "conversations")
}

// IndexConversationByID retrieves one conversation by UUID and indexes it.
func (s *IndexService) IndexConversationByID(ctx context.Context, idStr string) (int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid conversation id: %w", err)
	}

	c, err := s.q.DocketConversationRead(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("read conversation: %w", err)
	}

	rec := fugusdk.ObjectRecord{
		ID:   c.ID.String(),
		Text: c.Name,
		Metadata: map[string]interface{}{
			"docket_gov_id": c.DocketGovID,
			"description":   c.Description,
		},
		Namespace: "conversations",
	}

	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	// Use the convenience method for single object upsert
	response, err := client.AddOrUpdateObject(ctx, rec)
	if err != nil {
		return 0, fmt.Errorf("index conversation: %w", err)
	}

	// Log the response
	log.Printf("Successfully indexed conversation %s: %s", idStr, response.Message)
	return 1, nil
}

// IndexAllOrganizations retrieves all organizations and batch indexes them in chunks.
func (s *IndexService) IndexAllOrganizations(ctx context.Context) (int, error) {
	rows, err := s.q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all organizations: %w", err)
	}

	var recs []fugusdk.ObjectRecord
	for _, o := range rows {
		recs = append(recs, fugusdk.ObjectRecord{
			ID:   o.ID.String(),
			Text: o.Name,
			Metadata: map[string]interface{}{
				"description":              o.Description,
				"is_person":                o.IsPerson.Bool,
				"total_documents_authored": o.TotalDocumentsAuthored,
			},
			Namespace: "organizations", // Add namespace for better organization
		})
	}

	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	// Process in chunks to avoid batch size limits
	return s.processBatchInChunks(ctx, client, recs, "organizations")
}

// IndexOrganizationByID retrieves one organization by UUID and indexes it.
func (s *IndexService) IndexOrganizationByID(ctx context.Context, idStr string) (int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid organization id: %w", err)
	}

	o, err := s.q.OrganizationRead(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("read organization: %w", err)
	}

	rec := fugusdk.ObjectRecord{
		ID:   o.ID.String(),
		Text: o.Name,
		Metadata: map[string]interface{}{
			"description": o.Description,
			"is_person":   o.IsPerson.Bool,
		},
		Namespace: "organizations",
	}

	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	// Use the convenience method for single object upsert
	response, err := client.AddOrUpdateObject(ctx, rec)
	if err != nil {
		return 0, fmt.Errorf("index organization: %w", err)
	}

	// Log the response
	log.Printf("Successfully indexed organization %s: %s", idStr, response.Message)
	return 1, nil
}

// DeleteConversationFromIndex removes a conversation from the search index
func (s *IndexService) DeleteConversationFromIndex(ctx context.Context, idStr string) error {
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
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

// DeleteOrganizationFromIndex removes an organization from the search index
func (s *IndexService) DeleteOrganizationFromIndex(ctx context.Context, idStr string) error {
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return fmt.Errorf("new fugu client: %w", err)
	}

	response, err := client.DeleteObject(ctx, idStr)
	if err != nil {
		return fmt.Errorf("delete organization from index: %w", err)
	}

	log.Printf("Successfully deleted organization %s from index: %s", idStr, response.Message)
	return nil
}

// IndexAllData is a convenience method to index all conversations and organizations
func (s *IndexService) IndexAllData(ctx context.Context) (int, int, error) {
	convCount, err := s.IndexAllConversations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("index conversations: %w", err)
	}

	orgCount, err := s.IndexAllOrganizations(ctx)
	if err != nil {
		return convCount, 0, fmt.Errorf("index organizations: %w", err)
	}

	log.Printf("Successfully indexed %d conversations and %d organizations", convCount, orgCount)
	return convCount, orgCount, nil
}

// SearchConversations searches for conversations using the Fugu search API
func (s *IndexService) SearchConversations(ctx context.Context, query string, page, perPage int) (*fugusdk.SanitizedResponse, error) {
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return nil, fmt.Errorf("new fugu client: %w", err)
	}

	// Search with namespace filter for conversations
	filters := []string{"conversations"}
	return client.AdvancedSearch(ctx, query, filters, page, perPage)
}

// SearchOrganizations searches for organizations using the Fugu search API
func (s *IndexService) SearchOrganizations(ctx context.Context, query string, page, perPage int) (*fugusdk.SanitizedResponse, error) {
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return nil, fmt.Errorf("new fugu client: %w", err)
	}

	// Search with namespace filter for organizations
	filters := []string{"organizations"}
	return client.AdvancedSearch(ctx, query, filters, page, perPage)
}

// SearchAll searches across all indexed data
func (s *IndexService) SearchAll(ctx context.Context, query string, page, perPage int) (*fugusdk.SanitizedResponse, error) {
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return nil, fmt.Errorf("new fugu client: %w", err)
	}

	return client.AdvancedSearch(ctx, query, nil, page, perPage)
}

// processBatchInChunks handles large batches by splitting them into smaller chunks
func (s *IndexService) processBatchInChunks(ctx context.Context, client *fugusdk.Client, recs []fugusdk.ObjectRecord, entityType string) (int, error) {
	const chunkSize = 500 // Use 500 to stay well under the 1000 limit
	totalProcessed := 0

	log.Printf("Processing %d %s in chunks of %d", len(recs), entityType, chunkSize)

	for i := 0; i < len(recs); i += chunkSize {
		end := i + chunkSize
		if end > len(recs) {
			end = len(recs)
		}

		chunk := recs[i:end]
		log.Printf("Processing chunk %d-%d of %d %s", i+1, end, len(recs), entityType)

		response, err := client.BatchUpsertObjects(ctx, chunk)
		if err != nil {
			return totalProcessed, fmt.Errorf("batch index %s chunk %d-%d: %w", entityType, i+1, end, err)
		}

		chunkProcessed := len(chunk)
		if response.UpsertedCount != nil {
			chunkProcessed = *response.UpsertedCount
		}

		totalProcessed += chunkProcessed
		log.Printf("Successfully processed chunk %d-%d: %s (processed: %d)", i+1, end, response.Message, chunkProcessed)
	}

	log.Printf("Successfully indexed %d %s in total", totalProcessed, entityType)
	return totalProcessed, nil
}
