package admin

import (
	"context"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/quickwit"
	"net/http"

	"github.com/google/uuid"
)

func HandleQuickwitIngestFromPostgresFactory(dbtx_val dbstore.DBTX, filter_out_unverified bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Printf("Starting Quickwit ingest from Postgres (filter_out_unverified=%v)\n", filter_out_unverified)
		err := QuickwitIngestFromPostgresMain(dbtx_val, ctx, filter_out_unverified)
		if err != nil {
			errorstring := fmt.Sprintf("Error ingesting from postgres: %v", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully completed Quickwit ingest from Postgres")
		w.Write([]byte("Sucessfully ingested from postgres"))
	}
}

func QuickwitIngestFromPostgresMain(dbtx_val dbstore.DBTX, ctx context.Context, filter_out_unverified bool) error {
	q := *dbstore.New(dbtx_val)
	indexName := quickwit.NYPUCIndexName
	var files []dbstore.File
	var err error

	if filter_out_unverified {
		return fmt.Errorf("Not implemented yet")
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
	chunkSize := 100
	fmt.Printf("Got %d file ids, processing in chunks of size %d\n", len(ids), chunkSize)

	fmt.Printf("Attempting to clear index %s\n", indexName)
	err = quickwit.ClearIndex(indexName)
	if err != nil {
		return err
	}
	fmt.Printf("Cleared index %s\n", indexName)
	err = quickwit.CreateDocketsQuickwitIndex(indexName)
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
		complete_files_chunk, err := CompleteFileSchemasFromUUIDs(ctx, q, id_chunk)
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
		err = quickwit.IngestIntoIndex(indexName, quickwit_data_list_chunk)
		if err != nil {
			fmt.Printf("Error ingesting into index: %v\n", err)
			return err
		}
		fmt.Printf("Processed chunk %d to %d\n", i, end-1)
	}

	return nil
}
