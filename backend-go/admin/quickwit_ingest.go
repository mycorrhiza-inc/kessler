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
		err := QuickwitIngestFromPostgresMain(dbtx_val, ctx, filter_out_unverified)
		if err != nil {
			errorstring := fmt.Sprintf("Error ingesting from postgres: %v", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Sucessfully ingested from postgres"))
	}
}

func QuickwitIngestFromPostgresMain(dbtx_val dbstore.DBTX, ctx context.Context, filter_out_unverified bool) error {
	q := *dbstore.New(dbtx_val)
	indexName := quickwit.NYPUCIndexName
	err := quickwit.ClearIndex(indexName)
	if err != nil {
		return err
	}
	var files []dbstore.File

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
	complete_files, err := CompleteFileSchemasFromUUIDs(ctx, q, ids)
	if err != nil {
		return err
	}
	fmt.Printf("Got %d complete files\n", len(complete_files))

	quickwit_data_list, err := quickwit.ResolveFileSchemaForDocketIngest(complete_files)
	err = quickwit.IngestIntoIndex(indexName, quickwit_data_list)
	if err != nil {
		return err
	}

	return nil
}
