package indexing

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"kessler/internal/dbstore"
	"kessler/internal/database"
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

// IndexAllConversations retrieves all conversations and batch indexes them.
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
		})
	}
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}
	// Batch ingest
	if err := client.IngestObjects(ctx, recs); err != nil {
		return 0, fmt.Errorf("batch index conversations: %w", err)
	}
	return len(recs), nil
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
	}
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}
	// Single ingest
	if err := client.IngestObjects(ctx, []fugusdk.ObjectRecord{rec}); err != nil {
		return 0, fmt.Errorf("index conversation: %w", err)
	}
	return 1, nil
}

// IndexAllOrganizations retrieves all organizations and batch indexes them.
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
		})
	}
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}
	// Batch ingest
	if err := client.IngestObjects(ctx, recs); err != nil {
		return 0, fmt.Errorf("batch index organizations: %w", err)
	}
	return len(recs), nil
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
	}
	client, err := fugusdk.NewClient(ctx, s.fuguURL)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}
	// Single ingest
	if err := client.IngestObjects(ctx, []fugusdk.ObjectRecord{rec}); err != nil {
		return 0, fmt.Errorf("index organization: %w", err)
	}
	return 1, nil
}
