// indexing/attachment_service.go
package indexing

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/database"
	"kessler/pkg/logger"
	"kessler/pkg/util"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// AttachmentIndexer handles attachment-specific indexing operations
type AttachmentIndexer struct {
	svc *IndexService
}

// NewAttachmentIndexer creates a new attachment indexer
func NewAttachmentIndexer(svc *IndexService) *AttachmentIndexer {
	return &AttachmentIndexer{svc: svc}
}

// IndexAllAttachments retrieves all attachments and batch indexes them in chunks as data records.
func (ai *AttachmentIndexer) IndexAllAttachments(ctx context.Context) (int, error) {
	q := database.GetQueries(ai.svc.db)

	// Use the generated SQLC method that matches our indexing needs
	rows, err := q.GetAllSearchAttachments(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all searchable attachments: %w", err)
	}

	if len(rows) == 0 {
		log.Printf("No attachments to index")
		return 0, nil
	}

	// Process attachments in parallel with worker pool
	const maxWorkers = 10 // Limit concurrent processing to avoid overwhelming the system
	workers := maxWorkers
	if len(rows) < maxWorkers {
		workers = len(rows)
	}

	// Channels for work distribution
	attachmentChan := make(chan dbstore.GetAllSearchAttachmentsRow, len(rows))
	resultChan := make(chan attachmentProcessingResult, len(rows))

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go ai.attachmentWorker(ctx, q, attachmentChan, resultChan, &wg)
	}

	// Send work to workers
	for _, row := range rows {
		attachmentChan <- row
	}
	close(attachmentChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results
	var allRecords []fugusdk.ObjectRecord
	skippedCount := 0
	segmentedCount := 0

	for result := range resultChan {
		if result.err != nil {
			log.Printf("Skipping attachment %s: %v", result.attachmentID, result.err)
			skippedCount++
			continue
		}
		if result.segmented {
			segmentedCount++
		}
		allRecords = append(allRecords, result.records...)
	}

	if skippedCount > 0 {
		log.Printf("Skipped %d attachments with invalid content", skippedCount)
	}
	if segmentedCount > 0 {
		log.Printf("Segmented %d attachments due to length limits", segmentedCount)
	}
	if len(allRecords) == 0 {
		log.Printf("No valid attachments to index")
		return 0, nil
	}

	log.Printf("Successfully processed %d attachments into %d records using %d workers",
		len(rows)-skippedCount, len(allRecords), workers)

	client, err := ai.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	return ai.svc.processBatchInChunks(ctx, client, allRecords, "attachments")
}

// attachmentProcessingResult holds the result of processing a single attachment
type attachmentProcessingResult struct {
	attachmentID string
	records      []fugusdk.ObjectRecord
	segmented    bool
	err          error
}

// attachmentWorker processes attachments from the channel
func (ai *AttachmentIndexer) attachmentWorker(
	ctx context.Context,
	q *dbstore.Queries,
	attachmentChan <-chan dbstore.GetAllSearchAttachmentsRow,
	resultChan chan<- attachmentProcessingResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for row := range attachmentChan {
		// Prepare createdAt pointer
		var createdAt *time.Time
		if row.CreatedAt.Valid {
			createdAt = &row.CreatedAt.Time
		}

		records, segmented, err := ai.prepareAttachmentRecords(ctx, q, attachmentRecordParams{
			id:        row.ID,
			fileID:    row.FileID,
			name:      row.Name,
			createdAt: createdAt,
			mdata:     row.Mdata,
			rawText:   row.Text.String,
		})

		resultChan <- attachmentProcessingResult{
			attachmentID: row.ID.String(),
			records:      records,
			segmented:    segmented,
			err:          err,
		}
	}
}

// IndexAttachmentByID retrieves one attachment by ID and indexes it as a data record.
func (ai *AttachmentIndexer) IndexAttachmentByID(ctx context.Context, idStr string) (int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid attachment id: %w", err)
	}

	q := database.GetQueries(ai.svc.db)
	row, err := q.GetSearchAttachmentById(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("read attachment: %w", err)
	}

	var createdAt *time.Time
	if row.CreatedAt.Valid {
		createdAt = &row.CreatedAt.Time
	}

	// Prepare records using shared logic
	records, _, err := ai.prepareAttachmentRecords(ctx, q, attachmentRecordParams{
		id:        row.ID,
		fileID:    row.FileID,
		name:      row.Name,
		createdAt: createdAt,
		mdata:     row.Mdata,
		rawText:   row.Text.String,
	})
	if err != nil {
		return 0, fmt.Errorf("prepare attachment %s: %w", idStr, err)
	}

	client, err := ai.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	totalIndexed := 0
	for _, rec := range records {
		resp, err := client.AddOrUpdateObject(ctx, rec)
		if err != nil {
			return totalIndexed, fmt.Errorf("index attachment segment %s: %w", rec.ID, err)
		}
		totalIndexed++
		log.Printf("Successfully indexed %s: %s", rec.ID, resp.Message)
	}

	return totalIndexed, nil
}

// DeleteAttachmentFromIndex removes an attachment from the search index.
func (ai *AttachmentIndexer) DeleteAttachmentFromIndex(ctx context.Context, idStr string) error {
	client, err := ai.svc.createFuguClient(ctx)
	if err != nil {
		return fmt.Errorf("new fugu client: %w", err)
	}

	response, err := client.DeleteObject(ctx, idStr)
	if err != nil {
		return fmt.Errorf("delete attachment from index: %w", err)
	}

	log.Printf("Successfully deleted attachment %s from index: %s", idStr, response.Message)
	return nil
}

// GetAttachmentStats returns statistics about attachments in the database using SQLC queries
func (ai *AttachmentIndexer) GetAttachmentStats(ctx context.Context) (map[string]interface{}, error) {
	q := database.GetQueries(ai.svc.db)

	// Get search stats using SQLC
	searchStats, err := q.GetAttachmentSearchStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment search stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_attachments":          searchStats.TotalCount,
		"attachments_with_text":      searchStats.WithTextCount,
		"attachments_without_text":   searchStats.WithoutTextCount,
		"text_extraction_percentage": float64(searchStats.WithTextCount) / float64(searchStats.TotalCount) * 100,
		"namespace":                  ai.svc.defaultNamespace,
	}

	// Get file extension distribution using SQLC
	extStats, err := q.GetAttachmentExtensionStats(ctx, 20)
	if err != nil {
		logger.Warn(ctx, "failed to get file extension distribution", zap.Error(err))
	} else {
		extensionStats := make(map[string]int64)
		for _, ext := range extStats {
			extensionStats[ext.Extension] = ext.Count
		}
		stats["file_extension_distribution"] = extensionStats
	}

	// Get creation date range using SQLC
	dateRange, err := q.GetAttachmentDateRange(ctx)
	if err != nil {
		logger.Warn(ctx, "failed to get date range", zap.Error(err))
	} else {
		if dateRange.EarliestDate != nil && dateRange.LatestDate != nil {
			// Handle the interface{} types - they should be time.Time
			if earliest, ok := dateRange.EarliestDate.(time.Time); ok {
				if latest, ok := dateRange.LatestDate.(time.Time); ok {
					stats["earliest_attachment"] = earliest.Format(time.RFC3339)
					stats["latest_attachment"] = latest.Format(time.RFC3339)
					stats["date_range_days"] = int(latest.Sub(earliest).Hours() / 24)
				}
			}
		}
	}

	return stats, nil
}

// ReindexAttachmentsWithMissingMetadata finds and reindexes attachments that may be missing important metadata using SQLC
func (ai *AttachmentIndexer) ReindexAttachmentsWithMissingMetadata(ctx context.Context) (int, error) {
	q := database.GetQueries(ai.svc.db)

	// Get all attachment IDs that need reindexing using SQLC
	attachmentIDs, err := q.GetAttachmentsNeedingReindex(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to find attachments needing reindexing: %w", err)
	}

	if len(attachmentIDs) == 0 {
		return 0, nil
	}

	// Reindex each attachment
	successCount := 0
	for _, id := range attachmentIDs {
		if _, err := ai.IndexAttachmentByID(ctx, id.String()); err != nil {
			logger.Error(ctx, "failed to reindex attachment",
				zap.String("attachment_id", id.String()),
				zap.Error(err))
		} else {
			successCount++
		}
	}

	log.Printf("Successfully reindexed %d out of %d attachments", successCount, len(attachmentIDs))
	return successCount, nil
}

// SearchAttachmentsByExtension searches for attachments by file extension using SQLC
func (ai *AttachmentIndexer) SearchAttachmentsByExtension(ctx context.Context, extension string) ([]map[string]interface{}, error) {
	q := database.GetQueries(ai.svc.db)

	params := dbstore.GetAttachmentsByExtensionParams{
		Lower: extension,
		Limit: 1000,
	}

	rows, err := q.GetAttachmentsByExtension(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search attachments by extension: %w", err)
	}

	var results []map[string]interface{}
	for _, row := range rows {
		result := map[string]interface{}{
			"id":        row.ID.String(),
			"name":      row.Name,
			"extension": row.Extension,
			"has_text":  row.HasText,
		}

		if row.CreatedAt.Valid {
			result["created_at"] = row.CreatedAt.Time.Format(time.RFC3339)
		}

		results = append(results, result)
	}

	return results, nil
}

// GetAttachmentsByDateRange returns attachments within a date range using SQLC
func (ai *AttachmentIndexer) GetAttachmentsByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	q := database.GetQueries(ai.svc.db)

	if limit <= 0 {
		limit = 100
	}

	params := dbstore.GetAttachmentsByDateRangeParams{
		CreatedAt:   pgtype.Timestamptz{Time: startDate, Valid: true},
		CreatedAt_2: pgtype.Timestamptz{Time: endDate, Valid: true},
		Limit:       int32(limit),
	}

	rows, err := q.GetAttachmentsByDateRange(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachments by date range: %w", err)
	}

	var results []map[string]interface{}
	for _, row := range rows {
		result := map[string]interface{}{
			"id":          row.ID.String(),
			"name":        row.Name,
			"extension":   row.Extension,
			"has_text":    row.HasText,
			"text_length": row.TextLength,
		}

		if row.CreatedAt.Valid {
			result["created_at"] = row.CreatedAt.Time.Format(time.RFC3339)
		}

		results = append(results, result)
	}

	return results, nil
}

// ValidateAttachmentData validates attachment data before indexing
func (ai *AttachmentIndexer) ValidateAttachmentData(attachmentID string) error {
	if attachmentID == "" {
		return fmt.Errorf("attachment ID cannot be empty")
	}

	if _, err := uuid.Parse(attachmentID); err != nil {
		return fmt.Errorf("invalid attachment UUID format: %w", err)
	}

	return nil
}

// Helper methods

// prepareAttachmentRecords prepares object records for an attachment with segmentation and metadata.
type attachmentRecordParams struct {
	id        uuid.UUID
	fileID    uuid.UUID
	name      string
	createdAt *time.Time
	mdata     []byte
	rawText   string
}

func (ai *AttachmentIndexer) prepareAttachmentRecords(ctx context.Context, q *dbstore.Queries, params attachmentRecordParams) ([]fugusdk.ObjectRecord, bool, error) {
	log := logger.FromContext(ctx)
	id := params.id
	name := params.name
	createdAt := params.createdAt
	mdata := params.mdata
	rawText := params.rawText
	fileID := params.fileID
	if fileID == uuid.Nil {
		return nil, false, fmt.Errorf("fileID was nil")
	}

	// Trim and validate text content
	text := strings.TrimSpace(rawText)
	if text == "" {
		return nil, false, fmt.Errorf("attachment %s has no valid text content", id.String())
	}

	// Extract Author IDs
	author_rows, err := q.AuthorshipDocumentListOrganizations(ctx, fileID)
	if err != nil {
		log.Error("Failed author lookup for file ingest", zap.String("file_id", fileID.String()))
		return nil, false, fmt.Errorf("looking up authors for document failed: %s", fileID)
	}
	extract_org_ids := func(row dbstore.AuthorshipDocumentListOrganizationsRow) uuid.UUID {
		return row.OrganizationID
	}
	author_ids := util.Map(author_rows, extract_org_ids)
	// Lookup org id
	convo_rows, err := q.ConversationIDFetchFromFileID(ctx, fileID)
	if err != nil {
		log.Error("Failed conversation lookup for file ingest", zap.String("file_id", fileID.String()))
		return nil, false, fmt.Errorf("looking up conversation for document failed: %s", fileID)
	}
	if len(convo_rows) == 0 {
		log.Error("No conversations found for file", zap.String("file_id", fileID.String()))
		return nil, false, fmt.Errorf("no conversations found for file: %s", fileID)
	}
	if len(convo_rows) > 1 {
		log.Warn("File has more then one conversation", zap.String("file_id", fileID.String()), zap.Int("convo_number", len(convo_rows)))
	}
	convo_id := convo_rows[0].ConversationUuid

	metaParams := attachmentMetadataParams{
		id:        id,
		fileID:    fileID,
		convoID:   convo_id,
		authorIDs: author_ids,
		name:      name,
		createdAt: createdAt,
		mdata:     mdata,
	}
	baseMetadata, facets := ai.buildAttachmentMetadataAndFacets(metaParams)

	// Parse date from metadata if available
	if dateStr, ok := baseMetadata["date"].(string); ok {
		if parsedTime, err := ai.parseDate(dateStr); err == nil {
			baseMetadata["date_iso"] = parsedTime.Format(time.RFC3339)
		} else {
			logger.Warn(ctx, "could not parse date from metadata",
				zap.String("attachment_id", id.String()),
				zap.String("date", dateStr),
				zap.Error(err))
		}
	}

	// Split text into segments
	maxLen := 9500
	segments := ai.splitTextIntoSegments(text, maxLen)
	if len(segments) > 1 {
		logger.Info(ctx, "splitting attachment into segments due to length",
			zap.String("attachment_id", id.String()),
			zap.Int("original_length", len(text)),
			zap.Int("num_segments", len(segments)))
	}

	// Build object records for each segment using goroutines
	records := make([]fugusdk.ObjectRecord, len(segments))
	var wg sync.WaitGroup

	for i, segment := range segments {
		wg.Add(1)
		go func(segmentIndex int, segmentText string) {
			defer wg.Done()

			// Copy base metadata for this goroutine
			metadata := make(map[string]interface{}, len(baseMetadata)+5)
			for k, v := range baseMetadata {
				metadata[k] = v
			}

			// Add segment-specific metadata
			if len(segments) > 1 {
				metadata["is_segmented"] = true
				metadata["segment_index"] = segmentIndex
				metadata["total_segments"] = len(segments)
				metadata["original_text_length"] = len(text)
				metadata["segment_text_length"] = len(segmentText)
			} else {
				metadata["is_segmented"] = false
				metadata["segment_index"] = 0
				metadata["total_segments"] = 1
			}

			// Unique ID per segment
			recID := id.String()
			if len(segments) > 1 {
				recID = fmt.Sprintf("%s-segment-%d", id.String(), segmentIndex)
			}

			// Store record at the correct index (no mutex needed since each goroutine writes to a unique index)
			records[segmentIndex] = fugusdk.ObjectRecord{
				ID:        recID,
				Text:      segmentText,
				Metadata:  metadata,
				Facets:    facets, // Facets are read-only, safe to share
				Namespace: ai.svc.defaultNamespace,
				DataType:  "data/attachment",
			}
		}(i, segment)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return records, len(segments) > 1, nil
}

// Helper methods

type attachmentMetadataParams struct {
	id        uuid.UUID
	fileID    uuid.UUID
	authorIDs []uuid.UUID
	convoID   uuid.UUID
	name      string
	extension string
	createdAt *time.Time
	mdata     []byte
}

// buildAttachmentMetadataAndFacets creates both metadata and facets for an attachment record
func (ai *AttachmentIndexer) buildAttachmentMetadataAndFacets(params attachmentMetadataParams) (map[string]interface{}, []string) {
	metadata := make(map[string]interface{})
	var facets []string

	// Add namespace facets
	facets = append(facets, ai.svc.defaultNamespace)
	facets = append(facets, fmt.Sprintf("%s/data/attachment", ai.svc.defaultNamespace))

	// Add attachment-specific metadata and facets
	if params.name != "" {
		metadata["file_name"] = params.name
	}
	if params.extension != "" {
		metadata["file_extension"] = params.extension
	}
	if params.createdAt != nil {
		metadata["created_at"] = params.createdAt.Format(time.RFC3339)
	}

	// Core metadata fields
	metadata["attachment_id"] = params.id.String()
	metadata["file_id"] = params.fileID.String()
	metadata["conversation_id"] = params.convoID.String()
	metadata["migrated_at"] = time.Now().Format(time.RFC3339)
	metadata["entity_type"] = "attachment"

	// Core facets with embedded values
	facets = append(facets, fmt.Sprintf("metadata/attachment_id/%s", params.id.String()))
	facets = append(facets, fmt.Sprintf("metadata/file_id/%s", params.fileID.String()))
	facets = append(facets, fmt.Sprintf("metadata/conversation_id/%s", params.convoID.String()))
	facets = append(facets, fmt.Sprintf("metadata/entity_type/%s", "attachment"))

	// Author IDs
	if len(params.authorIDs) > 0 {
		transform_into_string := func(id uuid.UUID) string {
			return id.String()
		}
		authorIDStrings := util.Map(params.authorIDs, transform_into_string)
		metadata["author_ids"] = authorIDStrings

		// Add facets for each author ID
		for _, authorID := range authorIDStrings {
			facets = append(facets, fmt.Sprintf("metadata/author_ids/%s", authorID))
		}
	}

	// Parse and handle raw metadata
	if len(params.mdata) > 0 {
		rawMdataStr := string(params.mdata)
		metadata["raw_mdata"] = rawMdataStr

		// Try to parse JSON metadata and create facets from it
		var parsedMdata map[string]interface{}
		if err := json.Unmarshal(params.mdata, &parsedMdata); err == nil {
			// Add parsed metadata fields and create facets
			for key, value := range parsedMdata {
				// Store the parsed value in metadata
				metadata[key] = value

				// Create facets for each parsed field
				if valueStr, ok := value.(string); ok && valueStr != "" {
					facets = append(facets, fmt.Sprintf("metadata/%s/%s", key, valueStr))
				}
			}
		} else {
			logger.Warn(context.Background(), "failed to parse raw metadata as JSON",
				zap.String("attachment_id", params.id.String()),
				zap.Error(err))
		}
	}

	return metadata, facets
}

// parseDate attempts to parse various date formats (matching Python script logic)
func (ai *AttachmentIndexer) parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"1/2/2006",   // M/D/Y
		"01/02/2006", // MM/DD/YYYY
		"2006-01-02", // YYYY-MM-DD
		"2006/01/02", // YYYY/MM/DD
		time.RFC3339, // ISO format
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// splitTextIntoSegments splits long text into smaller segments that fit within the character limit
func (ai *AttachmentIndexer) splitTextIntoSegments(text string, maxLength int) []string {
	if len(text) <= maxLength {
		return []string{text}
	}

	var segments []string

	// Work with the original string to ensure we're counting bytes correctly
	for i := 0; i < len(text); {
		end := i + maxLength
		if end > len(text) {
			end = len(text)
		}

		segment := text[i:end]

		// Try to break at word boundaries if we're not at the end
		if end < len(text) {
			// Look backwards for a space to break at (within last 20% of segment)
			searchStart := len(segment) * 4 / 5 // Start searching from 80% of the way through
			lastSpace := strings.LastIndex(segment[searchStart:], " ")
			if lastSpace >= 0 {
				// Adjust to actual position in segment
				lastSpace += searchStart
				segment = segment[:lastSpace]
				i += lastSpace + 1 // Move past the space for next segment
			} else {
				// No good break point found, use the full segment
				i = end
			}
		} else {
			// This is the last segment
			i = end
		}

		trimmed := strings.TrimSpace(segment)
		if len(trimmed) > 0 {
			// Double-check the length before adding
			if len(trimmed) > maxLength {
				logger.Error(nil, "segment still too long after splitting",
					zap.Int("segment_length", len(trimmed)),
					zap.Int("max_length", maxLength))
				// Force truncate as emergency fallback
				trimmed = trimmed[:maxLength-3] + "..."
			}
			segments = append(segments, trimmed)
		}
	}

	return segments
}
