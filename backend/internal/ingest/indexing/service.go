// indexing/serv// indexing/service.go
package indexing

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"kessler/internal/database"
	"kessler/internal/dbstore"
	"kessler/internal/fugusdk"
	"kessler/pkg/logger"
)

// IndexService fetches DB records and indexes them into FuguDB.
type IndexService struct {
	fuguURL          string
	q                *dbstore.Queries
	defaultNamespace string // e.g., "NYPUC"
}

// NewIndexService constructs an IndexService pointing at fuguURL.
func NewIndexService(fuguURL string) *IndexService {
	// Use the shared connection pool (implements dbstore.DBTX)
	db := database.ConnPool
	return &IndexService{
		fuguURL:          fuguURL,
		q:                dbstore.New(db),
		defaultNamespace: "NYPUC", // Configure this based on your organization
	}
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

// IndexAllConversations retrieves all conversations and batch indexes them in chunks.
func (s *IndexService) IndexAllConversations(ctx context.Context) (int, error) {
	rows, err := s.q.DocketConversationList(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all conversations: %w", err)
	}

	var recs []fugusdk.ObjectRecord
	skippedCount := 0

	for _, c := range rows {
		// Create meaningful text field with fallbacks
		text := s.createConversationText(c.Name, c.Description, c.DocketGovID, c.ID.String())

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
			},
			// Use proper namespace facet structure (categorical only)
			Namespace:      s.defaultNamespace,
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

	client, err := s.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	return s.processBatchInChunks(ctx, client, recs, "conversations")
}

// IndexAllOrganizations retrieves all organizations and batch indexes them in chunks.
func (s *IndexService) IndexAllOrganizations(ctx context.Context) (int, error) {
	rows, err := s.q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetch all organizations: %w", err)
	}

	var recs []fugusdk.ObjectRecord
	skippedCount := 0

	for _, o := range rows {
		// Create meaningful text field with fallbacks
		text := s.createOrganizationText(o.Name, o.Description, o.ID.String())

		if text == "" {
			log.Printf("Skipping organization %s - no valid text content", o.ID.String())
			skippedCount++
			continue
		}

		recs = append(recs, fugusdk.ObjectRecord{
			ID:   o.ID.String(),
			Text: text,
			Metadata: map[string]interface{}{
				"description":              o.Description,
				"is_person":                o.IsPerson.Bool,
				"total_documents_authored": o.TotalDocumentsAuthored,
				"organization_name":        o.Name, // Store specific name in metadata
			},
			// Use proper namespace facet structure (categorical only)
			Namespace:    s.defaultNamespace,
			Organization: o.Name,         // This triggers namespace/NYPUC/organization facet
			DataType:     "organization", // This triggers namespace/NYPUC/data facet
		})
	}

	if skippedCount > 0 {
		log.Printf("Skipped %d organizations with empty content", skippedCount)
	}

	if len(recs) == 0 {
		log.Printf("No valid organizations to index")
		return 0, nil
	}

	client, err := s.createFuguClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("new fugu client: %w", err)
	}

	return s.processBatchInChunks(ctx, client, recs, "organizations")
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

	// Create meaningful text field with fallbacks
	text := s.createConversationText(c.Name, c.Description, c.DocketGovID, c.ID.String())
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
		},
		// Use proper namespace facet structure (categorical only)
		Namespace:      s.defaultNamespace,
		ConversationID: c.ID.String(),
		DataType:       "conversation",
	}

	client, err := s.createFuguClient(ctx)
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

	// Create meaningful text field with fallbacks
	text := s.createOrganizationText(o.Name, o.Description, o.ID.String())
	if text == "" {
		return 0, fmt.Errorf("organization %s has no valid text content and cannot be indexed", idStr)
	}

	rec := fugusdk.ObjectRecord{
		ID:   o.ID.String(),
		Text: text,
		Metadata: map[string]interface{}{
			"description":       o.Description,
			"is_person":         o.IsPerson.Bool,
			"organization_name": o.Name, // Store specific name in metadata
		},
		// Use proper namespace facet structure (categorical only)
		Namespace:    s.defaultNamespace,
		Organization: o.Name,
		DataType:     "organization",
	}

	client, err := s.createFuguClient(ctx)
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

// Helper methods for namespace facet handling

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

// createConversationText creates a meaningful text field for conversation indexing
// with multiple fallback options to ensure we always have searchable content
func (s *IndexService) createConversationText(name, description, docketGovID, id string) string {
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

// createOrganizationText creates a meaningful text field for organization indexing
// with multiple fallback options to ensure we always have searchable content
func (s *IndexService) createOrganizationText(name, description, id string) string {
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

// Delete methods remain the same
func (s *IndexService) DeleteConversationFromIndex(ctx context.Context, idStr string) error {
	client, err := s.createFuguClient(ctx)
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

func (s *IndexService) DeleteOrganizationFromIndex(ctx context.Context, idStr string) error {
	client, err := s.createFuguClient(ctx)
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
