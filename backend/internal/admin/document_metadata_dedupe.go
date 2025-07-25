package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/authors"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/files"
	"kessler/pkg/database"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type DocumentMetadataCheck struct {
	NamedDocketID string `json:"named_docket_id"`
	DateString    string `json:"date_string"`
	Name          string `json:"name"`
	Extension     string `json:"extension"`
	AuthorString  string `json:"author_string"`
}

func (h *AdminHandler) HandleCheckDocumentMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var args DocumentMetadataCheck
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	file, err := h.CheckDocumentMetadata(args, ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking document metadata: %v", err), http.StatusInternalServerError)
		return
	}

	if file.ID == uuid.Nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(file); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *AdminHandler) CheckDocumentMetadata(args DocumentMetadataCheck, ctx context.Context) (files.CompleteFileSchema, error) {
	params := dbstore.FileCheckForDuplicatesParams{
		Name:        args.Name,
		Extension:   args.Extension,
		DocketGovID: args.NamedDocketID,
	}
	q := database.GetQueries(h.db)
	results, err := q.FileCheckForDuplicates(ctx, params)
	if err != nil {
		return files.CompleteFileSchema{}, err
	}
	if len(results) == 0 {
		return files.CompleteFileSchema{}, nil
	}

	// Check all results for matching metadata
	type matchedFile struct {
		row      dbstore.FileCheckForDuplicatesRow
		metadata map[string]interface{}
	}
	var matchingFiles []matchedFile
	for _, result := range results {
		var metadata map[string]interface{}
		err = json.Unmarshal(result.Mdata, &metadata)
		if err != nil {
			log.Info(fmt.Sprintf("Warning: failed to unmarshal metadata for file %v: %v", result.ID, err))
			continue
		}

		// Check if date and author match
		if date, ok := metadata["date"].(string); ok && date == args.DateString {
			if author, ok := metadata["author"].(string); ok && author == args.AuthorString {
				matchingFiles = append(matchingFiles, matchedFile{result, metadata})
			}
		}
	}

	// If no files match the metadata criteria
	if len(matchingFiles) == 0 {
		return files.CompleteFileSchema{}, nil
	}

	// Use the first matching file
	matched := matchingFiles[0]
	if len(matchingFiles) > 1 {
		log.Info(fmt.Sprintf("Warning: multiple files match all criteria (name, extension, docket, date, author): %v", matchingFiles))
	}

	nilSchema := files.CompleteFileSchema{}
	var extra_obj files.FileGeneratedExtras
	err = json.Unmarshal(matched.row.ExtraObj, &extra_obj)
	if err != nil {
		return nilSchema, fmt.Errorf("error unmarshalling file extras: %w", err)
	}

	// Create a minimal ConversationInformation since we only have the UUID
	conv_info := conversations.ConversationInformation{
		ID:          matched.row.ConversationUuid.Bytes,
		DocketGovID: matched.row.DocketGovID.String,
	}

	file := files.CompleteFileSchema{
		ID:           matched.row.ID,
		Verified:     matched.row.Verified.Bool,
		Lang:         matched.row.Lang,
		Name:         matched.row.Name,
		Mdata:        matched.metadata,
		Extra:        extra_obj,
		Conversation: conv_info,
		// Note: Authors field is left empty as FileCheckForDuplicatesRow doesn't contain author information
		Authors: []authors.AuthorInformation{},
	}
	return file, nil
}
