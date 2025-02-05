package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/authors"
	"kessler/objects/conversations"
	"kessler/objects/files"
	"kessler/objects/timestamp"
	"kessler/quickwit"
	"kessler/util"
	"net/http"
)

func HandleQuickwitFileIngestFromPostgres(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := util.DBQueriesFromRequest(r)
	include_unverified := r.URL.Query().Get("include_unverified") == "true"
	filter_out_unverified := !include_unverified

	fmt.Printf("Starting Quickwit ingest from Postgres (filter_out_unverified=%v)\n", filter_out_unverified)

	err := QuickwitIngestFromPostgres(q, ctx, filter_out_unverified)
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting from postgres: %v", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully completed Quickwit ingest from Postgres")
	w.Write([]byte("Sucessfully ingested from postgres"))
}

func ParseQuickwitFileIntoCompleteSchema(file_raw dbstore.SemiCompleteFileQuickwitListGetRow) (files.CompleteFileSchema, error) {
	var mdata_obj files.FileMetadataSchema
	err := json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		return files.CompleteFileSchema{}, err
	}
	author_list := make([]authors.AuthorInformation, len(file_raw.OrganizationIds))
	for i := range file_raw.OrganizationIds {
		author_list[i] = authors.AuthorInformation{
			AuthorName: file_raw.OrganizationNames[i],
			AuthorID:   file_raw.OrganizationIds[i],
		}
	}
	text_list := make([]files.FileChildTextSource, len(file_raw.FileTexts))
	for i := range file_raw.FileTexts {
		text_list[i] = files.FileChildTextSource{
			IsOriginalText: false,
			Text:           file_raw.FileTexts[i],
			Language:       file_raw.FileTextLanguages[i],
		}
	}

	file := files.CompleteFileSchema{
		ID:            file_raw.ID,
		Verified:      file_raw.Verified.Bool,
		Extension:     file_raw.Extension,
		Lang:          "en",
		Name:          file_raw.Name,
		Hash:          file_raw.Hash,
		IsPrivate:     false,
		DatePublished: timestamp.KesslerTime(file_raw.DatePublished.Time),
		Mdata:         mdata_obj,
		Stage:         files.DocProcStage{},
		Extra:         files.FileGeneratedExtras{},
		Authors:       author_list,
		Conversation: conversations.ConversationInformation{
			DocketGovID: file_raw.DocketGovID.String,
			ID:          file_raw.ConversationUuid.Bytes,
		},
		DocTexts: text_list,
	}
	return file, nil
}

func QuickwitIngestFromPostgres(q *dbstore.Queries, ctx context.Context, filter_out_unverified bool) error {
	indexName := quickwit.NYPUCIndexName
	var files []dbstore.File
	var err error

	files_raw, err := q.SemiCompleteFileQuickwitListGet(ctx)
	if err != nil {
		return err
	}

	if filter_out_unverified {
		fmt.Printf("Got raw n files from postgres: %d\n", len(files))
		var new_raw_files []dbstore.SemiCompleteFileQuickwitListGetRow

		for _, file := range files_raw {
			if file.Verified.Bool {
				new_raw_files = append(new_raw_files, file)
			}
		}
		files_raw = new_raw_files
	}
	complete_file_schema_results := util.MapErrorDiscard(files_raw, ParseQuickwitFileIntoCompleteSchema)

	quickwit_data_list_chunk, err := quickwit.ResolveFileSchemaForDocketIngest(complete_file_schema_results)
	// Randomize the uuids so that you dont have weird unexpected behavior near the beginning or end.
	err = quickwit.IngestIntoIndex(indexName, quickwit_data_list_chunk, true)

	return nil
}
