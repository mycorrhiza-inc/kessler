package jobs

import (
	"context"
	"fmt"
	"kessler/admin"
	"kessler/gen/dbstore"
	"kessler/quickwit"
	"kessler/util"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
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

func QuickwitIngestFromPostgres(q *dbstore.Queries, ctx context.Context, filter_out_unverified bool) error {
	indexName := quickwit.NYPUCIndexName
	var files []dbstore.File
	var err error

	files, err = q.SemiCompleteFileQuickwitListGet(ctx)
	if err != nil {
		return err
	}

	if filter_out_unverified {
		fmt.Printf("Got raw n files from postgres: %d\n", len(files))
		var new_files []dbstore.File

		for _, file := range files {
			if file.Verified.Bool {
				verifiedFiles = append(verifiedFiles, file)
			}
		}
		files = verifiedFiles

	} else {
		files, err = q.FilesList(ctx)
		if err != nil {
			return err
		}
	}
	ids := make([]uuid.UUID, len(files))
	for i, file := range files {
		ids[i] = file.ID
	}
	// Randomize the uuids so that you dont have weird unexpected behavior near the beginning or end.
	for index := range ids {
		rand_index := rand.Intn(index + 1)
		ids[index], ids[rand_index] = ids[rand_index], ids[index]
	}
	// chunkSize := 100
	chunkSize := 100
	fmt.Printf("Got %d file ids, processing in chunks of size %d\n", len(ids), chunkSize)

	fmt.Printf("Attempting to clear index %s\n", indexName)
	err = quickwit.ClearIndex(indexName, false)
	if err != nil {
		return err
	}
	fmt.Printf("Cleared index %s\n", indexName)
	err = quickwit.CreateQuickwitNYFileIndex("")
	if err != nil {
		return err
	}
	fmt.Printf("Created index %s\n", indexName)

	for i := 0; i < len(ids); i += chunkSize {

		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		id_chunk := ids[i:end]
		fmt.Printf("Processing chunk %d to %d\n", i, end-1)
		complete_files_chunk, err := admin.CompleteFileSchemasFromUUIDs(ctx, id_chunk)
		if err != nil {
			fmt.Printf("Error getting complete file schemas: %v\n", err)
			return err
		}
		fmt.Printf("Got %d complete files\n", len(complete_files_chunk))

		quickwit_data_list_chunk, err := quickwit.ResolveFileSchemaForDocketIngest(complete_files_chunk)
		if err != nil {
			fmt.Printf("Error resolving file schema: %v\n", err)
			return err
		}
		fmt.Printf("Sucessfully parsed that into %d quickwit entries\n", len(quickwit_data_list_chunk))
		err = quickwit.IngestIntoIndex(indexName, quickwit_data_list_chunk, false)
		if err != nil {
			fmt.Printf("Error ingesting into index: %v\n", err)
			return err
		}
		fmt.Printf("Processed chunk %d to %d\n", i, end-1)
	}

	return nil
}
