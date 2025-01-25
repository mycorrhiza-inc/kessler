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
	if len(results) > 1 {
		log.Printf("Warning, more than one file with the same name and extension and DocketGovID: %v", results)
	}
	file_row := results[0]

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
