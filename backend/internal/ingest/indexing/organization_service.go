// indexing/organizations.go
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

// OrganizationIndexer handles organization-specific indexing operations
type OrganizationIndexer struct {
	svc *IndexService
}

// NewOrganizationIndexer creates a new organization indexer
func NewOrganizationIndexer(svc *IndexService) *OrganizationIndexer {
	return &OrganizationIndexer{svc: svc}
}

// IndexAllOrganizations retrieves all organizations and batch indexes them in chunks.
func (oi *OrganizationIndexer) IndexAllOrganizations(ctx context.Context) (int, error) {
	q := database.GetQueries(oi.svc.db)
	rows, err := q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all organizations: %w", err)
	}

	var recs []fugusdk.ObjectRecord
	skippedCount := 0

	for _, o := range rows {
		// Create meaningful text field with fallbacks
		text := oi.createOrganizationText(o.Name, o.Description, o.ID.String())

		if text == "" {
			log.Printf("Skipping organization %s - no valid text content", o.ID.String())
			skippedCount++
			continue
		}

		metadata, facets := oi.buildOrganizationMetadataAndFacets(organizationParams{
			id:                     o.ID,
			name:                   o.Name,
			description:            o.Description,
			isPerson:               o.IsPerson.Bool,
			totalDocumentsAuthored: o.TotalDocumentsAuthored,
		})

		recs = append(recs, fugusdk.ObjectRecord{
			ID:        o.ID.String(),
			Text:      text,
			Metadata:  metadata,
			Facets:    facets,
			Namespace: oi.svc.defaultNamespace,
			DataType:  "data/organization",
		})
	}

	if skippedCount > 0 {
		log.Printf("Skipped %d organizations with empty content", skippedCount)
	}

	if len(recs) == 0 {
		log.Printf("No valid organizations to index")
		return 0, nil
	}

	client, err := oi.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	return oi.svc.processBatchInChunks(ctx, client, recs, "organizations")
}

// IndexOrganizationByID retrieves one organization by UUID and indexes it.
func (oi *OrganizationIndexer) IndexOrganizationByID(ctx context.Context, idStr string) (int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid organization id: %w", err)
	}

	q := database.GetQueries(oi.svc.db)
	o, err := q.OrganizationRead(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("read organization: %w", err)
	}

	// Create meaningful text field with fallbacks
	text := oi.createOrganizationText(o.Name, o.Description, o.ID.String())
	if text == "" {
		return 0, fmt.Errorf("organization %s has no valid text content and cannot be indexed", idStr)
	}

	metadata, facets := oi.buildOrganizationMetadataAndFacets(organizationParams{
		id:          o.ID,
		name:        o.Name,
		description: o.Description,
		isPerson:    o.IsPerson.Bool,
		// Note: OrganizationRead might not have TotalDocumentsAuthored field
		// totalDocumentsAuthored: 0,
	})

	rec := fugusdk.ObjectRecord{
		ID:        o.ID.String(),
		Text:      text,
		Metadata:  metadata,
		Facets:    facets,
		Namespace: oi.svc.defaultNamespace,
		DataType:  "data/organization",
	}

	client, err := oi.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	response, err := client.AddOrUpdateObject(ctx, rec)
	if err != nil {
		return 0, fmt.Errorf("index organization: %w", err)
	}

	log.Printf("Successfully indexed organization %s: %s", idStr, response.Message)
	return 1, nil
}

// DeleteOrganizationFromIndex removes an organization from the search index.
func (oi *OrganizationIndexer) DeleteOrganizationFromIndex(ctx context.Context, idStr string) error {
	client, err := oi.svc.createFuguClient(ctx)
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

// organizationParams holds the parameters for building organization metadata and facets
type organizationParams struct {
	id                     uuid.UUID
	name                   string
	description            string
	isPerson               bool
	totalDocumentsAuthored int64
}

// buildOrganizationMetadataAndFacets creates both metadata and facets for an organization record
func (oi *OrganizationIndexer) buildOrganizationMetadataAndFacets(params organizationParams) (map[string]interface{}, []string) {
	metadata := map[string]interface{}{
		"entity_type":       "organization",
		"organization_id":   params.id.String(),
		"organization_name": params.name,
		"is_person":         params.isPerson,
	}

	if params.totalDocumentsAuthored > 0 {
		metadata["total_documents_authored"] = params.totalDocumentsAuthored
	}

	if params.description != "" {
		metadata["description"] = params.description
	}

	var facets []string

	// Add namespace facets
	facets = append(facets, oi.svc.defaultNamespace)
	facets = append(facets, fmt.Sprintf("%s/data/organization", oi.svc.defaultNamespace))

	// Core facets with embedded values
	facets = append(facets, "metadata/entity_type/organization")
	facets = append(facets, fmt.Sprintf("metadata/organization_id/%s", params.id.String()))

	// Add facets for organization name if not empty
	if params.name != "" {
		facets = append(facets, fmt.Sprintf("metadata/organization_name/%s", params.name))
	}

	// Add facets for person type
	if params.isPerson {
		facets = append(facets, "metadata/is_person/true")
	} else {
		facets = append(facets, "metadata/is_person/false")
	}

	return metadata, facets
}

// createOrganizationText creates a meaningful text field for organization indexing
// with multiple fallback options to ensure we always have searchable content
func (oi *OrganizationIndexer) createOrganizationText(name, description, id string) string {
	// Try name first (most common case)
	if text := strings.TrimSpace(name); text != "" {
		return text
	}

	// Fall back to description
	if text := strings.TrimSpace(description); text != "" {
		return text
	}

	// Last resort: use UUID with prefix
	return fmt.Sprintf("Organization %s", id)
}

// ValidateOrganizationData validates organization data before indexing
func (oi *OrganizationIndexer) ValidateOrganizationData(organizationID string) error {
	if organizationID == "" {
		return fmt.Errorf("organization ID cannot be empty")
	}

	if _, err := uuid.Parse(organizationID); err != nil {
		return fmt.Errorf("invalid organization UUID format: %w", err)
	}

	return nil
}

// GetOrganizationStats returns statistics about organizations in the database
func (oi *OrganizationIndexer) GetOrganizationStats(ctx context.Context) (map[string]interface{}, error) {
	q := database.GetQueries(oi.svc.db)

	// Get total count
	rows, err := q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization list: %w", err)
	}

	stats := map[string]interface{}{
		"total_organizations": len(rows),
		"namespace":           oi.svc.defaultNamespace,
	}

	// Analyze organization types and document counts
	personCount := 0
	organizationCount := 0
	totalDocuments := int64(0)
	textlessCount := 0
	documentDistribution := make(map[string]int)

	for _, o := range rows {
		if o.IsPerson.Bool {
			personCount++
		} else {
			organizationCount++
		}

		totalDocuments += o.TotalDocumentsAuthored

		// Categorize by document count
		docCount := o.TotalDocumentsAuthored
		switch {
		case docCount == 0:
			documentDistribution["no_documents"]++
		case docCount <= 5:
			documentDistribution["1_to_5_documents"]++
		case docCount <= 20:
			documentDistribution["6_to_20_documents"]++
		case docCount <= 100:
			documentDistribution["21_to_100_documents"]++
		default:
			documentDistribution["100_plus_documents"]++
		}

		// Check for meaningful text
		text := oi.createOrganizationText(o.Name, o.Description, o.ID.String())
		if text == fmt.Sprintf("Organization %s", o.ID.String()) {
			textlessCount++
		}
	}

	stats["person_count"] = personCount
	stats["organization_count"] = organizationCount
	stats["total_documents_authored"] = totalDocuments
	stats["document_distribution"] = documentDistribution
	stats["organizations_without_meaningful_text"] = textlessCount

	if len(rows) > 0 {
		stats["average_documents_per_organization"] = float64(totalDocuments) / float64(len(rows))
	}

	return stats, nil
}

// BulkUpdateOrganizationMetadata updates metadata for multiple organizations
func (oi *OrganizationIndexer) BulkUpdateOrganizationMetadata(ctx context.Context, updates map[string]map[string]interface{}) error {
	client, err := oi.svc.createFuguClient(ctx)
	if err != nil {
		return fmt.Errorf("new fugu client: %w", err)
	}

	for organizationID, metadata := range updates {
		// Validate organization ID
		if err := oi.ValidateOrganizationData(organizationID); err != nil {
			logger.Error(ctx, "invalid organization ID for metadata update",
				zap.String("organization_id", organizationID),
				zap.Error(err))
			continue
		}

		// Get existing organization to preserve text content
		id, _ := uuid.Parse(organizationID)
		q := database.GetQueries(oi.svc.db)
		o, err := q.OrganizationRead(ctx, id)
		if err != nil {
			logger.Error(ctx, "failed to read organization for metadata update",
				zap.String("organization_id", organizationID),
				zap.Error(err))
			continue
		}

		// Create text field
		text := oi.createOrganizationText(o.Name, o.Description, o.ID.String())
		if text == "" {
			logger.Warn(ctx, "skipping organization metadata update - no valid text content",
				zap.String("organization_id", organizationID))
			continue
		}

		// Build metadata and facets
		baseMetadata, facets := oi.buildOrganizationMetadataAndFacets(organizationParams{
			id:          o.ID,
			name:        o.Name,
			description: o.Description,
			isPerson:    o.IsPerson.Bool,
			// Note: OrganizationRead might not have TotalDocumentsAuthored field
			// totalDocumentsAuthored: 0,
		})

		// Add custom metadata updates
		for key, value := range metadata {
			baseMetadata[key] = value
		}

		rec := fugusdk.ObjectRecord{
			ID:        organizationID,
			Text:      text,
			Metadata:  baseMetadata,
			Facets:    facets,
			Namespace: oi.svc.defaultNamespace,
			DataType:  "data/organization",
		}

		if _, err := client.AddOrUpdateObject(ctx, rec); err != nil {
			logger.Error(ctx, "failed to update organization metadata",
				zap.String("organization_id", organizationID),
				zap.Error(err))
		} else {
			logger.Info(ctx, "successfully updated organization metadata",
				zap.String("organization_id", organizationID))
		}
	}

	return nil
}

// SearchOrganizationsByName searches for organizations by name pattern
func (oi *OrganizationIndexer) SearchOrganizationsByName(ctx context.Context, namePattern string) ([]map[string]interface{}, error) {
	q := database.GetQueries(oi.svc.db)

	// Get all organizations (in a real implementation, you might want to add a search query)
	rows, err := q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization list: %w", err)
	}

	var results []map[string]interface{}
	searchPattern := strings.ToLower(namePattern)

	for _, o := range rows {
		if strings.Contains(strings.ToLower(o.Name), searchPattern) {
			result := map[string]interface{}{
				"id":                       o.ID.String(),
				"name":                     o.Name,
				"description":              o.Description,
				"is_person":                o.IsPerson.Bool,
				"total_documents_authored": o.TotalDocumentsAuthored,
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// GetTopDocumentAuthors returns organizations sorted by document count
func (oi *OrganizationIndexer) GetTopDocumentAuthors(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	q := database.GetQueries(oi.svc.db)

	rows, err := q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization list: %w", err)
	}

	// Create a slice of organizations with their document counts
	type orgWithDocs struct {
		ID       string
		Name     string
		DocCount int64
		IsPerson bool
	}

	var orgsWithDocs []orgWithDocs
	for _, o := range rows {
		orgsWithDocs = append(orgsWithDocs, orgWithDocs{
			ID:       o.ID.String(),
			Name:     o.Name,
			DocCount: o.TotalDocumentsAuthored,
			IsPerson: o.IsPerson.Bool,
		})
	}

	// Sort by document count (simple bubble sort for small datasets)
	for i := 0; i < len(orgsWithDocs)-1; i++ {
		for j := 0; j < len(orgsWithDocs)-i-1; j++ {
			if orgsWithDocs[j].DocCount < orgsWithDocs[j+1].DocCount {
				orgsWithDocs[j], orgsWithDocs[j+1] = orgsWithDocs[j+1], orgsWithDocs[j]
			}
		}
	}

	// Convert to result format and apply limit
	var results []map[string]interface{}
	maxResults := limit
	if maxResults > len(orgsWithDocs) || maxResults <= 0 {
		maxResults = len(orgsWithDocs)
	}

	for i := 0; i < maxResults; i++ {
		org := orgsWithDocs[i]
		results = append(results, map[string]interface{}{
			"id":                       org.ID,
			"name":                     org.Name,
			"total_documents_authored": org.DocCount,
			"is_person":                org.IsPerson,
			"rank":                     i + 1,
		})
	}

	return results, nil
}
