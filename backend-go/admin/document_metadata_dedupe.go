package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/common/objects/authors"
	"kessler/common/objects/conversations"
	"kessler/common/objects/files"
	"kessler/util"
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

//	type FileCheckForDuplicatesRow struct {
//		ID               uuid.UUID
//		Name             string
//		Extension        string
//		Lang             string
//		Verified         pgtype.Bool
//		Hash             string
//		CreatedAt        pgtype.Timestamptz
//		UpdatedAt        pgtype.Timestamptz
//		DatePublished    pgtype.Timestamptz
//		Mdata            []byte
//		ExtraObj         []byte
//		ConversationUuid pgtype.UUID
//		DocketGovID      pgtype.Text
//	}
func HandleCheckDocumentMetadata(w http.ResponseWriter, r *http.Request) {
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
	q := *util.DBQueriesFromRequest(r)
	file, err := CheckDocumentMetadata(args, q, ctx)
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

func CheckDocumentMetadata(args DocumentMetadataCheck, q dbstore.Queries, ctx context.Context) (files.CompleteFileSchema, error) {
	params := dbstore.FileCheckForDuplicatesParams{
		Name:        args.Name,
		Extension:   args.Extension,
		DocketGovID: args.NamedDocketID,
	}
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
		Extension:    matched.row.Extension,
		Lang:         matched.row.Lang,
		Name:         matched.row.Name,
		Hash:         matched.row.Hash,
		Mdata:        matched.metadata,
		Extra:        extra_obj,
		Conversation: conv_info,
		// Note: Authors field is left empty as FileCheckForDuplicatesRow doesn't contain author information
		Authors: []authors.AuthorInformation{},
	}
	return file, nil
}

// Example Metadata Schema
// {"id": "14f92e3e-6748-4f64-8c1a-1200547bf8e9", "url": "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={2071B08B-0000-C730-9754-1C39327CCC7D}", "date": "11/08/2023", "lang": "en", "title": "Staff Proposal on the Transition of Utility Reported Community-Scale Energy Usage Data", "author": "New York State Department of Public Service", "source": "New York State Department of Public Service", "authors": ["New York State Department of Public Service"], "language": "en", "docket_id": "20-M-0082", "extension": "pdf", "file_class": "Plans and Proposals", "item_number": "249", "author_email": "", "author_organisation": "New York State Department of Public Service"}
