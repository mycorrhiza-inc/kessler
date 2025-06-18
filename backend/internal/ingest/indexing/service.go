// indexing/service.go
package indexing

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"

	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
)

// IndexService is the main service that coordinates indexing operations across all entity types
type IndexService struct {
	fuguURL          string
	db               dbstore.DBTX
	defaultNamespace string // e.g., "NYPUC"

	// Entity-specific indexers
	conversationIndexer *ConversationIndexer
	organizationIndexer *OrganizationIndexer
	attachmentIndexer   *AttachmentIndexer
}

// NewIndexService constructs an IndexService pointing at fuguURL.
func NewIndexService(fuguURL string, db dbstore.DBTX) *IndexService {
	svc := &IndexService{
		fuguURL:          fuguURL,
		defaultNamespace: "NYPUC", // Configure this based on your organization
		db:               db,
	}

	// Initialize entity-specific indexers
	svc.conversationIndexer = NewConversationIndexer(svc)
	svc.organizationIndexer = NewOrganizationIndexer(svc)
	svc.attachmentIndexer = NewAttachmentIndexer(svc)

	return svc
}

// DataRecord represents a single record for batch data ingestion
type DataRecord struct {
	ID       string                 `json:"id"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Optional namespace facet fields
	Namespace      string `json:"namespace,omitempty"`
	Organization   string `json:"organization,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	DataType       string `json:"data_type,omitempty"`
}

// DataIngestRequest represents the incoming batch data ingest request
type DataIngestRequest struct {
	Records []DataRecord `json:"records"`
}

// DataIngestResponse represents the response from batch data ingestion
type DataIngestResponse struct {
	Successful  int      `json:"successful"`
	Failed      int      `json:"failed"`
	FailedIDs   []string `json:"failed_ids,omitempty"`
	Message     string   `json:"message"`
	ProcessedAt string   `json:"processed_at"`
}

// Conversation-related methods (delegate to ConversationIndexer)
func (s *IndexService) IndexAllConversations(ctx context.Context) (int, error) {
	return s.conversationIndexer.IndexAllConversations(ctx)
}

func (s *IndexService) IndexConversationByID(ctx context.Context, idStr string) (int, error) {
	return s.conversationIndexer.IndexConversationByID(ctx, idStr)
}

func (s *IndexService) DeleteConversationFromIndex(ctx context.Context, idStr string) error {
	return s.conversationIndexer.DeleteConversationFromIndex(ctx, idStr)
}

// Organization-related methods (delegate to OrganizationIndexer)
func (s *IndexService) IndexAllOrganizations(ctx context.Context) (int, error) {
	return s.organizationIndexer.IndexAllOrganizations(ctx)
}

func (s *IndexService) IndexOrganizationByID(ctx context.Context, idStr string) (int, error) {
	return s.organizationIndexer.IndexOrganizationByID(ctx, idStr)
}

func (s *IndexService) DeleteOrganizationFromIndex(ctx context.Context, idStr string) error {
	return s.organizationIndexer.DeleteOrganizationFromIndex(ctx, idStr)
}

// Attachment-related methods (delegate to AttachmentIndexer)
func (s *IndexService) IndexAllAttachments(ctx context.Context) (int, error) {
	return s.attachmentIndexer.IndexAllAttachments(ctx)
}

func (s *IndexService) IndexAttachmentByID(ctx context.Context, idStr string) (int, error) {
	return s.attachmentIndexer.IndexAttachmentByID(ctx, idStr)
}

func (s *IndexService) DeleteAttachmentFromIndex(ctx context.Context, idStr string) error {
	return s.attachmentIndexer.DeleteAttachmentFromIndex(ctx, idStr)
}

// Bulk operations that coordinate across entity types
func (s *IndexService) IndexAllData(ctx context.Context) (int, int, int, error) {
	convCount, err := s.IndexAllConversations(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("index conversations: %w", err)
	}

	orgCount, err := s.IndexAllOrganizations(ctx)
	if err != nil {
		return convCount, 0, 0, fmt.Errorf("index organizations: %w", err)
	}

	attachCount, err := s.IndexAllAttachments(ctx)
	if err != nil {
		return convCount, orgCount, 0, fmt.Errorf("index attachments: %w", err)
	}

	log.Printf("Successfully indexed %d conversations, %d organizations, and %d attachments", convCount, orgCount, attachCount)
	return convCount, orgCount, attachCount, nil
}

// IndexCompleteData indexes all conversations, organizations, and attachments
func (s *IndexService) IndexCompleteData(ctx context.Context) (int, int, int, error) {
	convCount, err := s.IndexAllConversations(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("index conversations: %w", err)
	}

	orgCount, err := s.IndexAllOrganizations(ctx)
	if err != nil {
		return convCount, 0, 0, fmt.Errorf("index organizations: %w", err)
	}

	attachCount, err := s.IndexAllAttachments(ctx)
	if err != nil {
		return convCount, orgCount, 0, fmt.Errorf("index attachments: %w", err)
	}

	log.Printf("Successfully indexed %d conversations, %d organizations, and %d attachments",
		convCount, orgCount, attachCount)
	return convCount, orgCount, attachCount, nil
}

// ProcessBatchDataIngest handles the business logic for batch data ingestion
func (s *IndexService) ProcessBatchDataIngest(ctx context.Context, records []DataRecord) (*DataIngestResponse, error) {
	// Create fugu client
	client, err := s.createFuguClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Convert records to fugu format
	fuguObjects, conversionErrors := s.convertToFuguObjects(ctx, records)

	// Initialize response
	response := &DataIngestResponse{
		ProcessedAt: time.Now().Format(time.RFC3339),
		FailedIDs:   make([]string, 0),
	}

	// Add conversion errors to failed count
	for recordID, err := range conversionErrors {
		logger.Error(ctx, "record conversion failed",
			zap.String("record_id", recordID),
			zap.Error(err))
		response.Failed++
		response.FailedIDs = append(response.FailedIDs, recordID)
	}

	// Ingest valid objects if any
	if len(fuguObjects) > 0 {
		// Use the same chunked processing as other batch operations
		processed, err := s.processBatchInChunks(ctx, client, fuguObjects, "data")
		if err != nil {
			logger.Error(ctx, "fugu data ingestion failed", zap.Error(err))

			// Mark all remaining objects as failed
			for _, obj := range fuguObjects {
				response.Failed++
				response.FailedIDs = append(response.FailedIDs, obj.ID)
			}

			response.Message = "Data ingestion failed"
			return response, nil
		}

		// Ingestion successful
		response.Successful = processed
		logger.Info(ctx, "successfully ingested data records", zap.Int("count", response.Successful))
	}

	// Set final message
	if response.Failed == 0 {
		response.Message = "All records ingested successfully"
	} else if response.Successful == 0 {
		response.Message = "All records failed to ingest"
	} else {
		response.Message = fmt.Sprintf("Partial success: %d succeeded, %d failed", response.Successful, response.Failed)
	}

	return response, nil
}

// IndexDataRecordByID indexes a single data record by ID
func (s *IndexService) IndexDataRecordByID(ctx context.Context, record DataRecord) (int, error) {
	// Validate required fields
	if record.ID == "" {
		return 0, fmt.Errorf("data record ID cannot be empty")
	}

	if strings.TrimSpace(record.Text) == "" {
		return 0, fmt.Errorf("data record text cannot be empty")
	}

	// Create fugu client
	client, err := s.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create fugu client: %w", err)
	}

	// Create fugu object with proper namespace facets
	fuguObj := fugusdk.ObjectRecord{
		ID:       record.ID,
		Text:     strings.TrimSpace(record.Text),
		Metadata: record.Metadata,

		// Use namespace facet structure
		Namespace: s.getRecordNamespace(record),
		DataType:  s.getRecordDataType(record),
	}

	// Set optional namespace facet fields if provided
	if record.Organization != "" {
		fuguObj.Organization = record.Organization
	}
	if record.ConversationID != "" {
		fuguObj.ConversationID = record.ConversationID
	}

	// Add ingestion metadata
	if fuguObj.Metadata == nil {
		fuguObj.Metadata = make(map[string]interface{})
	}

	fuguObj.Metadata["ingested_at"] = time.Now().Format(time.RFC3339)
	fuguObj.Metadata["ingestion_source"] = "kessler-single-ingest"

	// Parse date from metadata if available
	if dateStr, ok := fuguObj.Metadata["date"].(string); ok {
		if parsedTime, err := s.parseDate(dateStr); err == nil {
			fuguObj.Metadata["date_iso"] = parsedTime.Format(time.RFC3339)
		} else {
			logger.Warn(ctx, "could not parse date from metadata",
				zap.String("record_id", record.ID),
				zap.String("date", dateStr),
				zap.Error(err))
		}
	}

	// Index the single record
	response, err := client.AddOrUpdateObject(ctx, fuguObj)
	if err != nil {
		return 0, fmt.Errorf("index data record: %w", err)
	}

	log.Printf("Successfully indexed data record %s: %s", record.ID, response.Message)
	return 1, nil
}

func (s *IndexService) DeleteDataRecordFromIndex(ctx context.Context, idStr string) error {
	client, err := s.createFuguClient(ctx)
	if err != nil {
		return fmt.Errorf("new fugu client: %w", err)
	}

	response, err := client.DeleteObject(ctx, idStr)
	if err != nil {
		return fmt.Errorf("delete data record from index: %w", err)
	}

	log.Printf("Successfully deleted data record %s from index: %s", idStr, response.Message)
	return nil
}

// Statistics and reporting methods that aggregate across all entity types
func (s *IndexService) GetSystemStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get stats from each indexer
	convStats, err := s.conversationIndexer.GetConversationStats(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get conversation stats", zap.Error(err))
		convStats = map[string]interface{}{"error": "failed to retrieve stats"}
	}
	stats["conversations"] = convStats

	orgStats, err := s.organizationIndexer.GetOrganizationStats(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get organization stats", zap.Error(err))
		orgStats = map[string]interface{}{"error": "failed to retrieve stats"}
	}
	stats["organizations"] = orgStats

	attachStats, err := s.attachmentIndexer.GetAttachmentStats(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get attachment stats", zap.Error(err))
		attachStats = map[string]interface{}{"error": "failed to retrieve stats"}
	}
	stats["attachments"] = attachStats

	// Add system-wide metadata
	stats["namespace"] = s.defaultNamespace
	stats["generated_at"] = time.Now().Format(time.RFC3339)
	stats["fugu_url"] = s.fuguURL

	return stats, nil
}

// GetIndexerByType returns the appropriate indexer for a given entity type
func (s *IndexService) GetIndexerByType(entityType string) interface{} {
	switch strings.ToLower(entityType) {
	case "conversation", "conversations":
		return s.conversationIndexer
	case "organization", "organizations":
		return s.organizationIndexer
	case "attachment", "attachments":
		return s.attachmentIndexer
	default:
		return nil
	}
}

// Helper methods for data processing and utilities

// convertToFuguObjects converts incoming data records to fugu ObjectRecord format
func (s *IndexService) convertToFuguObjects(ctx context.Context, records []DataRecord) ([]fugusdk.ObjectRecord, map[string]error) {
	fuguObjects := make([]fugusdk.ObjectRecord, 0, len(records))
	conversionErrors := make(map[string]error)

	for _, record := range records {
		// Validate required fields
		if record.ID == "" {
			conversionErrors["unknown"] = fmt.Errorf("record ID cannot be empty")
			continue
		}

		if strings.TrimSpace(record.Text) == "" {
			conversionErrors[record.ID] = fmt.Errorf("record text cannot be empty")
			continue
		}

		// Create fugu object with proper namespace facets
		fuguObj := fugusdk.ObjectRecord{
			ID:       record.ID,
			Text:     strings.TrimSpace(record.Text),
			Metadata: record.Metadata,

			// Use namespace facet structure
			Namespace: s.getRecordNamespace(record),
			DataType:  s.getRecordDataType(record),
		}

		// Set optional namespace facet fields if provided
		if record.Organization != "" {
			fuguObj.Organization = record.Organization
		}
		if record.ConversationID != "" {
			fuguObj.ConversationID = record.ConversationID
		}

		// Add ingestion metadata
		if fuguObj.Metadata == nil {
			fuguObj.Metadata = make(map[string]interface{})
		}

		fuguObj.Metadata["ingested_at"] = time.Now().Format(time.RFC3339)
		fuguObj.Metadata["ingestion_source"] = "kessler-batch-ingest"

		// Parse date from metadata if available
		if dateStr, ok := fuguObj.Metadata["date"].(string); ok {
			if parsedTime, err := s.parseDate(dateStr); err == nil {
				fuguObj.Metadata["date_iso"] = parsedTime.Format(time.RFC3339)
			} else {
				logger.Warn(ctx, "could not parse date from metadata",
					zap.String("record_id", record.ID),
					zap.String("date", dateStr),
					zap.Error(err))
			}
		}

		fuguObjects = append(fuguObjects, fuguObj)
	}

	return fuguObjects, conversionErrors
}

// getRecordNamespace returns the namespace for a record, using default if not specified
func (s *IndexService) getRecordNamespace(record DataRecord) string {
	if record.Namespace != "" {
		return record.Namespace
	}
	return s.defaultNamespace
}

// getRecordDataType returns the data type for a record, with intelligent defaults
func (s *IndexService) getRecordDataType(record DataRecord) string {
	if record.DataType != "" {
		return record.DataType
	}

	// Intelligent defaults based on record characteristics
	if record.ConversationID != "" {
		return "conversation-data"
	}
	if record.Organization != "" {
		return "organization-data"
	}

	return "general-data"
}

// parseDate attempts to parse various date formats
func (s *IndexService) parseDate(dateStr string) (time.Time, error) {
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

// processBatchInChunks handles large batches by splitting them into smaller chunks
func (s *IndexService) processBatchInChunks(ctx context.Context, client *fugusdk.Client, recs []fugusdk.ObjectRecord, entityType string) (int, error) {
	const chunkSize = 10 // Use 500 to stay well under the 1000 limit
	totalProcessed := 0

	log.Printf("Processing %d %s in chunks of %d", len(recs), entityType, chunkSize)

	for i := 0; i < len(recs); i += chunkSize {
		end := i + chunkSize
		if end > len(recs) {
			end = len(recs)
		}

		chunk := recs[i:end]
		log.Printf("Processing chunk %d-%d of %d %s", i+1, end, len(recs), entityType)

		// Use IngestObjectsWithNamespaceFacets for proper namespace facet handling
		response, err := client.IngestObjectsWithNamespaceFacets(ctx, chunk)
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

// createFuguClient creates a new FuguDB client for health checks and other operations
func (s *IndexService) createFuguClient(ctx context.Context) (*fugusdk.Client, error) {
	return fugusdk.NewClient(ctx, s.fuguURL)
}

// Health check and maintenance methods

// HealthCheck verifies the service and its dependencies are working
func (s *IndexService) HealthCheck(ctx context.Context) error {
	// Test database connection
	if err := s.db.QueryRow(ctx, "SELECT 1").Scan(new(int)); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Test FuguDB connection
	client, err := s.createFuguClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create fugu client: %w", err)
	}

	if err := client.Health(ctx); err != nil {
		return fmt.Errorf("fugu health check failed: %w", err)
	}

	return nil
}

// ValidateConfiguration checks that the service is properly configured
func (s *IndexService) ValidateConfiguration() error {
	if s.fuguURL == "" {
		return fmt.Errorf("fugu URL cannot be empty")
	}

	if s.defaultNamespace == "" {
		return fmt.Errorf("default namespace cannot be empty")
	}

	if s.db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}

	// Validate indexers are initialized
	if s.conversationIndexer == nil {
		return fmt.Errorf("conversation indexer not initialized")
	}
	if s.organizationIndexer == nil {
		return fmt.Errorf("organization indexer not initialized")
	}
	if s.attachmentIndexer == nil {
		return fmt.Errorf("attachment indexer not initialized")
	}

	return nil
}

// MaintenanceReindex performs a full reindex of all entities with error tracking
func (s *IndexService) MaintenanceReindex(ctx context.Context) (map[string]interface{}, error) {
	startTime := time.Now()
	results := make(map[string]interface{})

	// Index conversations
	convCount, convErr := s.IndexAllConversations(ctx)
	results["conversations"] = map[string]interface{}{
		"indexed": convCount,
		"error":   convErr,
	}

	// Index organizations
	orgCount, orgErr := s.IndexAllOrganizations(ctx)
	results["organizations"] = map[string]interface{}{
		"indexed": orgCount,
		"error":   orgErr,
	}

	// Index attachments
	attachCount, attachErr := s.IndexAllAttachments(ctx)
	results["attachments"] = map[string]interface{}{
		"indexed": attachCount,
		"error":   attachErr,
	}

	// Calculate totals and summary
	totalIndexed := convCount + orgCount + attachCount
	hasErrors := convErr != nil || orgErr != nil || attachErr != nil

	results["summary"] = map[string]interface{}{
		"total_indexed":    totalIndexed,
		"duration_seconds": time.Since(startTime).Seconds(),
		"has_errors":       hasErrors,
		"started_at":       startTime.Format(time.RFC3339),
		"completed_at":     time.Now().Format(time.RFC3339),
	}

	if hasErrors {
		return results, fmt.Errorf("maintenance reindex completed with errors")
	}

	log.Printf("Maintenance reindex completed successfully: %d total records indexed in %.2f seconds",
		totalIndexed, time.Since(startTime).Seconds())

	return results, nil
}

// CleanupOrphanedRecords removes records from the index that no longer exist in the database
// This is a placeholder for future implementation
func (s *IndexService) CleanupOrphanedRecords(ctx context.Context) error {
	// TODO: Implement cleanup logic
	// This would involve:
	// 1. Getting all indexed record IDs from FuguDB
	// 2. Checking which ones no longer exist in the database
	// 3. Removing orphaned records from the index

	logger.Info(ctx, "orphaned record cleanup not yet implemented")
	return nil
}
