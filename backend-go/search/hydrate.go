package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/authors"
	"kessler/objects/conversations"
	"kessler/objects/files"
	"log"

	"github.com/google/uuid"
)

func doShittyHydration(uuidList []uuid.UUID, ctx context.Context, q dbstore.Queries) ([]files.CompleteFileSchema, error) {
	intermediate_results, err := q.SemiCompleteFileListGet(ctx, uuidList)
	results := make([]files.CompleteFileSchema, len(intermediate_results))
	if err != nil {
		return nil, err
	}
	for i, schema := range intermediate_results {
		file, err := unrollSingleSchema(schema)
		if err != nil {
			log.Printf("Error hydrating file: %v", err)
		}
		results[i] = file
	}
	return results, nil
}

func unrollSingleSchema(schema dbstore.SemiCompleteFileListGetRow) (files.CompleteFileSchema, error) {
	var mdata_obj map[string]interface{}
	nilSchema := files.CompleteFileSchema{}
	err := json.Unmarshal(schema.Mdata, &mdata_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file metadata %v: %v\n", uuid, err)
		return nilSchema, errors.New(errorstring)
	}
	var extra_obj files.FileGeneratedExtras
	err = json.Unmarshal(schema.ExtraObj, &extra_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file extras %v: %v\n", uuid, err)
		return nilSchema, errors.New(errorstring)
	}
	// Missing info here, it doesnt have the name.
	conv_info := conversations.ConversationInformation{ID: schema.DocketUuid}
	org_id_list := schema.OrganizationIds
	org_name_list := schema.OrganizationNames
	org_personhood_list := schema.IsPersonList

	author_info := make([]authors.AuthorInformation, len(org_id_list))
	for i := range org_id_list {
		author_info[i] = authors.AuthorInformation{
			AuthorName:      org_name_list[i],
			IsPerson:        org_personhood_list[i],
			IsPrimaryAuthor: false,
			AuthorID:        org_id_list[i],
		}
	}

	file := files.CompleteFileSchema{
		ID:           schema.ID,
		Verified:     schema.Verified.Bool,
		Extension:    schema.Extension,
		Lang:         schema.Lang,
		Name:         schema.Name,
		Hash:         schema.Hash,
		Mdata:        mdata_obj,
		Extra:        extra_obj,
		Conversation: conv_info,
		Authors:      author_info,
	}
	return file, nil
}

func HydrateSearchResults(results []SearchData) []SearchData {
	// get a pg connection
	// query the database for all metadata of each result
	// return array of hydrated results
	// This could be done in a more efficient way
	return results
}

func SearchThenHydrate(r SearchRequest, ctx context.Context) ([]files.FileMetadataSchema, error) {
	data, err := SearchQuickwit(r)
	if err != nil {
		log.Printf("search-then-hydrate: error getting quickwit data", err)
	}
	idList := make([]string, len(data))
	for i, d := range data {
		idList[i] = d.SourceID
	}

	// hydratedData := make([]files.FileMetadataSchema, len(idList))

	return []files.FileMetadataSchema{}, nil
}
