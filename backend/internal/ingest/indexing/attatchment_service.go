// indexing/attachment_service.go
package indexing

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/database"
	"kessler/pkg/logger"
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

	var recs []fugusdk.ObjectRecord
	skippedCount := 0
	segmentedCount := 0

	for _, row := range rows {
		// Skip if no text content (extra safety check)
		if !row.Text.Valid || strings.TrimSpace(row.Text.String) == "" {
			log.Printf("Skipping attachment %s - no valid text content", row.ID.String())
			skippedCount++
			continue
		}

		// Prepare text content with length validation and segmentation
		text := strings.TrimSpace(row.Text.String)

		// Create the attachment record(s) as data record(s)
		var createdAt *time.Time
		if row.CreatedAt.Valid {
			createdAt = &row.CreatedAt.Time
		}

		baseMetadata := ai.buildAttachmentMetadata(row.ID, &row.Name, createdAt, row.Mdata)

		// Parse date from metadata if available (matching Python script logic)
		if dateStr, ok := baseMetadata["date"].(string); ok {
			if parsedTime, err := ai.parseDate(dateStr); err == nil {
				baseMetadata["date_iso"] = parsedTime.Format(time.RFC3339)
			} else {
				logger.Warn(ctx, "could not parse date from metadata",
					zap.String("attachment_id", row.ID.String()),
					zap.String("date", dateStr),
					zap.Error(err))
			}
		}

		// Split text into segments if necessary
		segments := ai.splitTextIntoSegments(text, 9900) // Leave buffer for safety

		if len(segments) > 1 {
			segmentedCount++
			logger.Info(ctx, "splitting attachment into segments due to length",
				zap.String("attachment_id", row.ID.String()),
				zap.Int("original_length", len(text)),
				zap.Int("num_segments", len(segments)))
		}

		// Create a record for each segment
		for segmentIndex, segment := range segments {
			// Create metadata for this segment
			metadata := make(map[string]interface{})
			for k, v := range baseMetadata {
				metadata[k] = v
			}

			// Add segment-specific metadata
			if len(segments) > 1 {
				metadata["is_segmented"] = true
				metadata["segment_index"] = segmentIndex
				metadata["total_segments"] = len(segments)
				metadata["original_text_length"] = len(text)
				metadata["segment_text_length"] = len(segment)
			} else {
				metadata["is_segmented"] = false
				metadata["segment_index"] = 0
				metadata["total_segments"] = 1
			}

			// Create unique ID for each segment
			recordID := row.ID.String()
			if len(segments) > 1 {
				recordID = fmt.Sprintf("%s-segment-%d", row.ID.String(), segmentIndex)
			}

			recs = append(recs, fugusdk.ObjectRecord{
				ID:       recordID,
				Text:     segment,
				Metadata: metadata,

				// Use proper namespace facet structure for data records
				Namespace: ai.svc.defaultNamespace,
				DataType:  "attachment", // This triggers namespace/NYPUC/data facet
			})
		}
	}

	if skippedCount > 0 {
		log.Printf("Skipped %d attachments with empty content", skippedCount)
	}

	if segmentedCount > 0 {
		log.Printf("Segmented %d attachments due to length limits", segmentedCount)
	}

	if len(recs) == 0 {
		log.Printf("No valid attachments to index")
		return 0, nil
	}

	client, err := ai.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	return ai.svc.processBatchInChunks(ctx, client, recs, "attachments")
}

// IndexAttachmentByID retrieves one attachment by ID and indexes it as a data record.
func (ai *AttachmentIndexer) IndexAttachmentByID(ctx context.Context, idStr string) (int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid attachment id: %w", err)
	}

	q := database.GetQueries(ai.svc.db)

	// Use the generated SQLC method for single attachment retrieval
	row, err := q.GetSearchAttachmentById(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("read attachment: %w", err)
	}

	// Check if we have valid text content
	if !row.Text.Valid || strings.TrimSpace(row.Text.String) == "" {
		return 0, fmt.Errorf("attachment %s has no valid text content and cannot be indexed", idStr)
	}

	// Create metadata
	var createdAt *time.Time
	if row.CreatedAt.Valid {
		createdAt = &row.CreatedAt.Time
	}

	baseMetadata := ai.buildAttachmentMetadata(row.ID, &row.Name, createdAt, row.Mdata)

	// Parse date from metadata if available
	if dateStr, ok := baseMetadata["date"].(string); ok {
		if parsedTime, err := ai.parseDate(dateStr); err == nil {
			baseMetadata["date_iso"] = parsedTime.Format(time.RFC3339)
		} else {
			logger.Warn(ctx, "could not parse date from metadata",
				zap.String("attachment_id", idStr),
				zap.String("date", dateStr),
				zap.Error(err))
		}
	}

	client, err := ai.svc.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	// Split text into segments if necessary
	text := strings.TrimSpace(row.Text.String)
	segments := ai.splitTextIntoSegments(text, 9900) // Leave buffer for safety

	if len(segments) > 1 {
		logger.Info(ctx, "splitting attachment into segments due to length",
			zap.String("attachment_id", idStr),
			zap.Int("original_length", len(text)),
			zap.Int("num_segments", len(segments)))
	}

	// Index each segment
	totalIndexed := 0
	for segmentIndex, segment := range segments {
		// Create metadata for this segment
		metadata := make(map[string]interface{})
		for k, v := range baseMetadata {
			metadata[k] = v
		}

		// Add segment-specific metadata
		if len(segments) > 1 {
			metadata["is_segmented"] = true
			metadata["segment_index"] = segmentIndex
			metadata["total_segments"] = len(segments)
			metadata["original_text_length"] = len(text)
			metadata["segment_text_length"] = len(segment)
		} else {
			metadata["is_segmented"] = false
			metadata["segment_index"] = 0
			metadata["total_segments"] = 1
		}

		// Create unique ID for each segment
		recordID := idStr
		if len(segments) > 1 {
			recordID = fmt.Sprintf("%s-segment-%d", idStr, segmentIndex)
		}

		rec := fugusdk.ObjectRecord{
			ID:       recordID,
			Text:     segment,
			Metadata: metadata,

			// Use proper namespace facet structure for data records
			Namespace: ai.svc.defaultNamespace,
			DataType:  "attachment", // This triggers namespace/NYPUC/data facet
		}

		response, err := client.AddOrUpdateObject(ctx, rec)
		if err != nil {
			return totalIndexed, fmt.Errorf("index attachment segment %d: %w", segmentIndex, err)
		}

		totalIndexed++
		log.Printf("Successfully indexed attachment %s segment %d/%d: %s",
			idStr, segmentIndex+1, len(segments), response.Message)
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

// buildAttachmentMetadata creates metadata for an attachment record
func (ai *AttachmentIndexer) buildAttachmentMetadata(id uuid.UUID, name *string, createdAt *time.Time, mdata []byte) map[string]interface{} {
	metadata := make(map[string]interface{})

	// Add raw mdata if available
	if len(mdata) > 0 {
		metadata["raw_mdata"] = string(mdata)
		// TODO: Parse JSON metadata if needed
		// var mdataMap map[string]interface{}
		// if err := json.Unmarshal(mdata, &mdataMap); err == nil {
		//     for k, v := range mdataMap {
		//         metadata[k] = v
		//     }
		// }
	}

	// Add attachment-specific metadata
	if name != nil {
		metadata["file_name"] = *name
		// Extract file extension if available
		if dotIndex := strings.LastIndex(*name, "."); dotIndex != -1 && dotIndex < len(*name)-1 {
			metadata["file_extension"] = strings.ToLower((*name)[dotIndex+1:])
		}
	}
	if createdAt != nil {
		metadata["created_at"] = createdAt.Format(time.RFC3339)
	}
	metadata["attachment_id"] = id.String()
	metadata["migrated_at"] = time.Now().Format(time.RFC3339)
	metadata["entity_type"] = "attachment"

	return metadata
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
	runes := []rune(text) // Use runes to handle Unicode properly

	for i := 0; i < len(runes); i += maxLength {
		end := i + maxLength
		if end > len(runes) {
			end = len(runes)
		}

		segment := string(runes[i:end])

		// Try to break at word boundaries if we're not at the end
		if end < len(runes) {
			// Look backwards for a space to break at
			lastSpace := strings.LastIndex(segment, " ")
			if lastSpace > maxLength/2 { // Only break at word if it's not too far back
				segment = segment[:lastSpace]
				i = i + lastSpace + 1 - maxLength // Adjust index to continue from after the space
			}
		}

		segments = append(segments, strings.TrimSpace(segment))
	}

	return segments
}
