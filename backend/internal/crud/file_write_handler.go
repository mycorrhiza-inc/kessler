package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

type FileUpsertHandlerConfig struct {
	Private      bool
	Insert       bool
	Deduplicate  bool
	IsAuthorized func(*http.Request) bool
	GetDBQueries func(*http.Request) *dbstore.Queries
}

// Chatgt came up with these, and I actually kind of like them
func respondError(w http.ResponseWriter, message string, statusCode int) {
	log.Info(message)
	http.Error(w, message, statusCode)
}

func makeFileUpsertHandler(config FileUpsertHandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := config.GetDBQueries(r)
		ctx := r.Context()

		// Validate HTTP method and path
		if !config.Insert && r.URL.Path == "/v2/public/files/insert" {
			respondError(w, "UNREACHABLE CODE: Using insert endpoint with update configuration", http.StatusInternalServerError)
			return
		}

		// Parse document UUID for updates
		docUUID, err := parseDocumentUUID(r, config.Insert)
		if err != nil {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse request body
		newDocInfo, err := parseRequestBody(r)
		if err != nil {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if newDocInfo.ID != uuid.Nil {
			newDocInfo.ID = docUUID
		}
		args := IngestDocParams{
			DocInfo:     newDocInfo,
			Insert:      config.Insert,
			Deduplicate: config.Deduplicate,
		}

		// Process file ingestion
		result, err := ingestFile(ctx, q, args)
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send response
		response, err := json.Marshal(result)
		if err != nil {
			respondError(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

type IngestDocParams struct {
	DocInfo     files.CompleteFileSchema
	Insert      bool
	Deduplicate bool
}

func ingestFile(ctx context.Context, q *dbstore.Queries, params IngestDocParams) (files.CompleteFileSchema, error) {
	var err error

	// Deduplication logic
	if config.IsInsert && config.Deduplicate {
		existingUUID, err := DeduplicateFile(ctx, q, docInfo.Hash)
		if err != nil {
			return docInfo, fmt.Errorf("deduplication error: %w", err)
		}
		if existingUUID != uuid.Nil {
			config.IsInsert = false
			docUUID = existingUUID
			docInfo.ID = existingUUID
		}
	}

	// Convert to creation data
	creationData := docInfo.ConvertToCreationData()
	creationData.Verified = pgtype.Bool{Bool: false, Valid: true}

	// Insert or update main file record
	fileSchema, err := upsertFileRecord(ctx, q, creationData, docUUID, config)
	if err != nil {
		return docInfo, fmt.Errorf("file upsert error: %w", err)
	}
	docInfo.ID = fileSchema.ID

	// Process associations
	associationErrors, hasErrored := processAssociations(ctx, q, docInfo)
	docInfo.Stage.IsErrored = docInfo.Stage.IsErrored || hasErrored
	docInfo.Stage.DatabaseErrorMsg += strings.Join(associationErrors, "\n")

	// Update verification status
	if err := updateVerificationStatus(ctx, q, docInfo); err != nil {
		return docInfo, fmt.Errorf("verification update error: %w", err)
	}

	return docInfo, nil
}

// dedupe.go - Deduplication logic

func DeduplicateFile(ctx context.Context, q *dbstore.Queries, hash string) (uuid.UUID, error) {
	ids, err := q.HashGetUUIDsFile(ctx, hash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("database error: %w", err)
	}
	if len(ids) > 0 {
		return ids[0], nil
	}
	return uuid.Nil, nil
}

func parseDocumentUUID(r *http.Request, isInsert bool) (uuid.UUID, error) {
	if isInsert {
		return uuid.Nil, nil
	}

	params := mux.Vars(r)
	fileIDString := params["uuid"]
	return uuid.Parse(fileIDString)
}

func parseRequestBody(r *http.Request) (files.CompleteFileSchema, error) {
	var docInfo files.CompleteFileSchema
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return docInfo, fmt.Errorf("error reading request body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &docInfo); err != nil {
		return docInfo, fmt.Errorf("error parsing request body: %w", err)
	}
	return docInfo, nil
}

func upsertFileRecord(ctx context.Context, q *dbstore.Queries, data dbstore.FileCreationData, docUUID uuid.UUID, config IngestConfig) (files.FileSchema, error) {
	if config.IsInsert {
		return InsertPubPrivateFileObj(q, ctx, data, config.Private)
	}
	return UpdatePubPrivateFileObj(q, ctx, data, config.Private, docUUID)
}

func processAssociations(ctx context.Context, q *dbstore.Queries, docInfo files.CompleteFileSchema) ([]string, bool) {
	var errors []string
	addError := func(err error, context string) {
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", context, err))
		}
	}

	addError(upsertFileAttachments(ctx, q, docInfo.ID, docInfo.Attachments, docInfo.IsInsert), "attachments")
	addError(upsertFileMetadata(ctx, q, docInfo.ID, docInfo.Mdata, docInfo.IsInsert), "metadata")
	addError(upsertFileExtras(ctx, q, docInfo.ID, docInfo.Extra, docInfo.IsInsert), "extras")
	addError(fileAuthorsUpsert(ctx, q, docInfo.ID, docInfo.Authors, docInfo.IsInsert), "authors")
	addError(fileConversationUpsert(ctx, q, docInfo.ID, docInfo.Conversation, docInfo.IsInsert), "conversation")

	return errors, len(errors) > 0
}

func updateVerificationStatus(ctx context.Context, q *dbstore.Queries, docInfo files.CompleteFileSchema) error {
	if !docInfo.Stage.IsErrored && docInfo.Stage.IsCompleted {
		params := dbstore.FileVerifiedUpdateParams{
			Verified: pgtype.Bool{Bool: true, Valid: true},
			ID:       docInfo.ID,
		}
		if _, err := q.FileVerifiedUpdate(ctx, params); err != nil {
			return fmt.Errorf("failed to update verification status: %w", err)
		}
	}
	return nil
}

// response.go - HTTP response helpers
