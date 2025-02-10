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
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func HandleQuickwitFileIngestFromPostgresFactory(isTest bool) func(http.ResponseWriter,
	*http.Request) {
	indexName := quickwit.NYPUCIndex
	if isTest {
		indexName = quickwit.TestNYPUCIndex
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := context.Background()
		ctx := util.CreateDBContextWithTimeout(time.Minute*90, 5)
		q := util.DBQueriesFromContext(ctx)
		include_unverified := r.URL.Query().Get("include_unverified") == "true"
		filter_out_unverified := !include_unverified

		log.Info(fmt.Sprintf("Starting Quickwit ingest from Postgres (filter_out_unverified=%v)\n", filter_out_unverified))

		err := QuickwitIngestFromPostgres(q, ctx, filter_out_unverified, indexName)
		if err != nil {
			errorstring := fmt.Sprintf("Error ingesting from postgres: %v", err)
			log.Info(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		log.Info("Successfully completed Quickwit ingest from Postgres")
		w.Write([]byte("Sucessfully ingested from postgres"))
	}
}

func parsePrimativeAuthorSchema(author_schemas []byte) ([]authors.AuthorInformation, error) {
	type SimpleAuthorSchema struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		IsPerson bool      `json:"is_person"`
	}
	var simple_schemas []SimpleAuthorSchema
	err := json.Unmarshal(author_schemas, &simple_schemas)
	if err != nil {
		return []authors.AuthorInformation{}, err
	}
	return_schemas := make([]authors.AuthorInformation, len(simple_schemas))
	for index, simple_schema := range simple_schemas {
		return_schemas[index] = authors.AuthorInformation{
			AuthorID:   simple_schema.ID,
			AuthorName: simple_schema.Name,
			IsPerson:   simple_schema.IsPerson,
		}
	}
	return return_schemas, nil
}

func ParseQuickwitFileIntoCompleteSchema(file_raw dbstore.Testmat) (files.CompleteFileSchema, error) {
	var mdata_obj files.FileMetadataSchema
	err := json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		return files.CompleteFileSchema{}, err
	}
	var extra_obj files.FileGeneratedExtras
	err = json.Unmarshal(file_raw.Mdata, &extra_obj)
	if err != nil {
		log.Info("encountered error decoding extras for file", "file_id", file_raw.ID)
		// return files.CompleteFileSchema{}, err
	}
	author_list, err := parsePrimativeAuthorSchema(file_raw.Organizations)
	if err != nil {
		log.Info("encountered error decoding author list for file", "file_id", file_raw.ID)
		// return files.CompleteFileSchema{}, err
	}
	text_list := []files.FileChildTextSource{
		{
			IsOriginalText: false,
			Text:           file_raw.FileText,
			Language:       "en",
		},
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
		Extra:         extra_obj,
		Authors:       author_list,
		Conversation: conversations.ConversationInformation{
			DocketGovID: file_raw.DocketGovID.String,
			ID:          file_raw.ConversationUuid.Bytes,
		},
		DocTexts: text_list,
	}
	return file, nil
}

func QuickwitIngestFromPostgres(q *dbstore.Queries, ctx context.Context, filter_out_unverified bool, indexName string) error {
	var files_raw []dbstore.Testmat

	page_size := 1000

	// Currently this encounters a hard cap at 10,000,000 files, so this should almost certainly be changed then. But at 1 second per request the ingest job should take 3 hours. So refactoring will be required.
	cap := 10000
	if indexName == quickwit.TestNYPUCIndex {
		cap = 3
	}

	for page := range cap {
		// for page := range 3 {
		pagination_params := dbstore.FilePrecomputedQuickwitListGetPaginatedParams{Limit: int32(page_size), Offset: int32(page * page_size)}
		temporary_file_results, err := q.FilePrecomputedQuickwitListGetPaginated(ctx, pagination_params)
		if err != nil {
			log.Info(fmt.Sprintf("Error getting semi complete file list: %v\n", err))

			time.Sleep(5 * time.Second)
			// return err
			break
		}
		if err == nil {
			files_raw = append(files_raw, temporary_file_results...)
			if len(temporary_file_results) < page_size {
				log.Info(fmt.Sprintf("Finished indexing PG after %v pages\n", page))
				break
			}
			log.Info(fmt.Sprintf("Indexed PG page %v\n", page))
		}
	}

	if filter_out_unverified {
		log.Info(fmt.Sprintf("Got raw n files from postgres: %d\n", len(files_raw)))
		var new_raw_files []dbstore.Testmat

		for _, file := range files_raw {
			if file.Verified.Bool {
				new_raw_files = append(new_raw_files, file)
			}
		}
		files_raw = new_raw_files
	}
	complete_file_schema_results := util.MapErrorDiscard(files_raw, ParseQuickwitFileIntoCompleteSchema)

	quickwit_data_list_chunk, err := quickwit.ResolveFileSchemaForDocketIngest(complete_file_schema_results)
	if err != nil {
		log.Info("Error converting complete file schema into quickwit schema for file inest: %s", err)
	}
	// Randomize the uuids so that you dont have weird unexpected behavior near the beginning or end.
	err = quickwit.IngestIntoIndex(indexName, quickwit_data_list_chunk, true)

	return nil
}
