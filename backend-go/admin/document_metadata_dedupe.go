package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/authors"
	"kessler/objects/conversations"
	"kessler/objects/files"
	"log"
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
	var matchingFiles []dbstore.FileCheckForDuplicatesRow
	for _, result := range results {
		var metadata map[string]interface{}
		err = json.Unmarshal(result.Mdata, &metadata)
		if err != nil {
			log.Printf("Warning: failed to unmarshal metadata for file %v: %v", result.ID, err)
			continue
		}

		// Check if date and author match
		if date, ok := metadata["date"].(string); ok && date == args.DateString {
			if author, ok := metadata["author"].(string); ok && author == args.AuthorString {
				matchingFiles = append(matchingFiles, result)
			}
		}
	}

	// If no files match the metadata criteria
	if len(matchingFiles) == 0 {
		return files.CompleteFileSchema{}, nil
	}

	// Use the first matching file
	file_row := matchingFiles[0]
	if len(matchingFiles) > 1 {
		log.Printf("Warning: multiple files match all criteria (name, extension, docket, date, author): %v", matchingFiles)
	}

	// Unmarshal metadata for the selected file
	// (bad for perf since your doing it twice, but I also dont care)
	var mdata_obj map[string]interface{}
	nilSchema := files.CompleteFileSchema{}
	err = json.Unmarshal(file_row.Mdata, &mdata_obj)
	if err != nil {
		return nilSchema, fmt.Errorf("error unmarshalling file metadata: %w", err)
	}

	var extra_obj files.FileGeneratedExtras
	err = json.Unmarshal(file_row.ExtraObj, &extra_obj)
	if err != nil {
		return nilSchema, fmt.Errorf("error unmarshalling file extras: %w", err)
	}

	// Create a minimal ConversationInformation since we only have the UUID
	conv_info := conversations.ConversationInformation{
		ID:          file_row.ConversationUuid.Bytes,
		DocketGovID: file_row.DocketGovID.String,
	}

	file := files.CompleteFileSchema{
		ID:           file_row.ID,
		Verified:     file_row.Verified.Bool,
		Extension:    file_row.Extension,
		Lang:         file_row.Lang,
		Name:         file_row.Name,
		Hash:         file_row.Hash,
		Mdata:        mdata_obj,
		Extra:        extra_obj,
		Conversation: conv_info,
		// Note: Authors field is left empty as FileCheckForDuplicatesRow doesn't contain author information
		Authors: []authors.AuthorInformation{},
	}
	return file, nil
}

// Example Metadata Schema
// {"id": "14f92e3e-6748-4f64-8c1a-1200547bf8e9", "url": "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={2071B08B-0000-C730-9754-1C39327CCC7D}", "date": "11/08/2023", "lang": "en", "title": "Staff Proposal on the Transition of Utility Reported Community-Scale Energy Usage Data", "author": "New York State Department of Public Service", "source": "New York State Department of Public Service", "authors": ["New York State Department of Public Service"], "language": "en", "docket_id": "20-M-0082", "extension": "pdf", "file_class": "Plans and Proposals", "item_number": "249", "author_email": "", "author_organisation": "New York State Department of Public Service"}
